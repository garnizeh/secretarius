package database

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/garnizeh/englog/internal/config"
	"github.com/garnizeh/englog/internal/sqlc"
	"github.com/garnizeh/englog/internal/store"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const (
	defaultMaxConns        = 100
	defaultMinConns        = 5
	defaultMaxConnIdleTime = time.Hour * 12
	defaultConnectTimeout  = time.Second
)

// DB is the main struct for the database connection pool.
type DB struct {
	rddb *pgxpool.Pool
	wrdb *pgxpool.Pool
}

// NewDB creates a new database connection pool and runs the migrations.
func NewDB(
	ctx context.Context,
	config config.DBConfig,
) (*DB, error) {
	config.Host = config.HostReadWrite

	// Migrate the database
	if err := migrate(ctx, sqlc.Migrations, config); err != nil {
		return nil, fmt.Errorf("error migrating database: %w", err)
	}

	// Setup the database connection
	dbRW, err := setupDBConn(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error setting up read write database connection: %w", err)
	}

	config.Host = config.HostReadOnly
	dbRO, err := setupDBConn(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error setting up read only database connection: %w", err)
	}

	db := DB{
		rddb: dbRO,
		wrdb: dbRW,
	}

	// Check if the database connection is alive
	if err := db.Check(ctx); err != nil {
		return nil, fmt.Errorf("error checking database connection: %w", err)
	}

	return &db, nil
}

// RDBMS returns the read write database connection pool.
func (db *DB) RDBMS() *pgxpool.Pool {
	return db.wrdb
}

// CloseDB closes the database connection pool.
func (db *DB) CloseDB() {
	db.rddb.Close()
	db.wrdb.Close()
}

// Check checks if the read write and the read only database connections are alive.
func (db *DB) Check(ctx context.Context) error {
	if err := db.wrdb.Ping(ctx); err != nil {
		return fmt.Errorf("error pinging read write database: %w", err)
	}

	if err := db.rddb.Ping(ctx); err != nil {
		return fmt.Errorf("error pinging read only database: %w", err)
	}

	if err := checkPool(ctx, db.wrdb); err != nil {
		return fmt.Errorf("error checking read write database: %w", err)
	}

	if err := checkPool(ctx, db.rddb); err != nil {
		return fmt.Errorf("error checking read only database: %w", err)
	}

	return nil
}

// Read executes a read-only transaction.
func (db *DB) Read(
	ctx context.Context,
	f func(qtx *store.Queries) error,
) error {
	txOpts := pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	}
	return db.transaction(ctx, db.rddb, txOpts, f)
}

// Write executes a read-write transaction.
func (db *DB) Write(
	ctx context.Context,
	f func(qtx *store.Queries) error,
) error {
	txOpts := pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
		IsoLevel:   pgx.Serializable,
	}

	return db.transaction(ctx, db.wrdb, txOpts, f)
}

// NoRows is useful to check if a query returned no rows.
func NoRows(err error) bool {
	return err != nil && errors.Is(err, sql.ErrNoRows)
}

func (db *DB) transaction(
	ctx context.Context,
	rdbms *pgxpool.Pool,
	txOpts pgx.TxOptions,
	f func(qtx *store.Queries) error,
) error {
	tx, err := rdbms.BeginTx(ctx, txOpts)
	if err != nil {
		return fmt.Errorf("error creating transaction: %w", err)
	}

	retries := 3
	for {
		if err := f(store.New(tx)); err != nil {
			rbErr := tx.Rollback(ctx)
			if rbErr != nil {
				err = errors.Join(err, rbErr)
			}
			return err
		}

		err = tx.Commit(ctx)
		if err == nil {
			return nil
		}

		// We found a commit error. We only retry if it is a serialization error.
		if isSerializationError(err) {
			retries--
			if retries >= 0 {
				continue
			}
		}

		// Not a serialization error, or we no longer can retry.
		return err
	}
}

func migrate(
	ctx context.Context,
	migrations embed.FS,
	config config.DBConfig,
) error {
	// We need to create the schema before running the migrations.
	// This is because the migrations are run in the context of the schema.
	// If the schema does not exist, the migrations will fail.
	if err := createSchema(ctx, migrations, config); err != nil {
		return fmt.Errorf("error creating schema %q: %w", config.Schema, err)
	}

	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("error setting goose dialect: %w", err)
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?application_name=englog&search_path=%s,public&connect_timeout=10&sslmode=disable",
		config.User, config.Password, config.Host, config.Name, config.Schema,
	)
	fmt.Println("DSN:", dsn)

	cfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("error parsing database config: %w", err)
	}

	db := stdlib.OpenDB(*cfg)

	if err := goose.Up(db, "schema"); err != nil {
		return errors.Join(
			db.Close(),
			fmt.Errorf("error running migrations: %w", err),
		)
	}

	return db.Close()
}

func createSchema(
	ctx context.Context,
	migrations embed.FS,
	config config.DBConfig,
) error {
	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("error setting dialect: %w", err)
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?application_name=englog&connect_timeout=10&sslmode=disable",
		config.User, config.Password, config.Host, config.Name,
	)
	fmt.Println("DSN:", dsn)

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return fmt.Errorf("error opening database connection: %w", err)
	}
	defer conn.Close(ctx)

	// Create required extensions first (globally in public schema)
	// This makes the functions available across all schemas
	if _, err := conn.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\""); err != nil {
		return fmt.Errorf("error creating uuid-ossp extension: %w", err)
	}

	// pg_stat_statements is a global extension
	if _, err := conn.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS \"pg_stat_statements\""); err != nil {
		return fmt.Errorf("error creating pg_stat_statements extension: %w", err)
	}

	// Create the schema
	sql := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", config.Schema)
	if _, err := conn.Exec(ctx, sql); err != nil {
		return fmt.Errorf("error creating schema: %w", err)
	}

	return nil
}

func isSerializationError(err error) bool {
	// This code was pulled directly from the pgx/tx_test.go to ensure we are detecting a serialization error.
	if pgErr, ok := err.(*pgconn.PgError); !ok || pgErr.Code != "40001" {
		return false
	}

	return true
}

func setupDBConn(ctx context.Context, config config.DBConfig) (*pgxpool.Pool, error) {
	poolConfig, err := toPoolConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creating pool config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error pinging pool: %w", err)
	}

	return pool, nil
}

// toPoolConfig converts the DBInfo struct to a pgxpool.Config struct.
func toPoolConfig(config config.DBConfig) (*pgxpool.Config, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?application_name=englog&search_path=%s,public&sslmode=disable",
		config.User, config.Password, config.Host, config.Name, config.Schema,
	)
	fmt.Println("DSN:", dsn)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	poolConfig.MinConns = defaultMinConns
	poolConfig.MaxConns = defaultMaxConns
	poolConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	poolConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	return poolConfig, nil
}

func checkPool(ctx context.Context, pool *pgxpool.Pool) error {
	// If the user doesn't give us a deadline set 1 second.
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Second)
		defer cancel()
	}

	for attempts := 1; ; attempts++ {
		if err := pool.Ping(ctx); err == nil {
			break
		}

		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Run a simple query to determine connectivity.
	// Running this query forces a round trip through the database.
	_, err := pool.Exec(ctx, "SELECT TRUE")
	return err
}
