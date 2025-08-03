package testutils

import (
	"context"
	"testing"

	"github.com/garnizeh/englog/internal/config"
	"github.com/garnizeh/englog/internal/database"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	postgresImage      = "postgres:17-alpine"
	databaseName       = "englog_test"
	username           = "testuser"
	password           = "testpass"
	postgresInitScript = "../../scripts/init-postgres.sql"
	schemaName         = "englog_test"
	port               = "5432"
)

// Helper function to setup test database
func DB(t testing.TB) *database.DB {
	t.Helper()

	ctx := context.Background()

	// Setup test database container
	pgContainer, err := postgres.Run(
		ctx,
		postgresImage,
		postgres.WithDatabase(databaseName),
		postgres.WithUsername(username),
		postgres.WithPassword(password),
		postgres.WithInitScripts(postgresInitScript),
		postgres.BasicWaitStrategies(),
	)
	require.NoError(t, err)

	// Clean up container when test ends
	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	})

	// Get connection configuration
	host, err := pgContainer.Host(ctx)
	require.NoError(t, err)
	port, err := pgContainer.MappedPort(ctx, port)
	require.NoError(t, err)

	hostPort := host + ":" + port.Port()

	config := config.DBConfig{
		User:          username,
		Password:      password,
		HostReadWrite: hostPort,
		HostReadOnly:  hostPort,
		Name:          databaseName,
		Schema:        schemaName,
	}

	db, err := database.NewDB(ctx, config)
	require.NoError(t, err)
	require.NotNil(t, db)

	// Cleanup database connections when test ends to prevent connection leaks
	t.Cleanup(func() {
		if db != nil {
			db.CloseDB()
		}
	})

	return db
}
