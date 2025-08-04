package testutils

import (
	"context"
	"sync"
	"testing"

	"github.com/garnizeh/englog/internal/config"
	"github.com/garnizeh/englog/internal/database"
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

var (
	sharedContainer *postgres.PostgresContainer
	sharedDB        *database.DB
	sharedOnce      sync.Once
	sharedMutex     sync.RWMutex
)

// setupSharedContainer initializes a shared test container for all tests
func setupSharedContainer() {
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
	if err != nil {
		panic("Failed to start test container: " + err.Error())
	}

	// Get connection configuration
	host, err := pgContainer.Host(ctx)
	if err != nil {
		panic("Failed to get container host: " + err.Error())
	}
	mappedPort, err := pgContainer.MappedPort(ctx, port)
	if err != nil {
		panic("Failed to get container port: " + err.Error())
	}

	hostPort := host + ":" + mappedPort.Port()

	config := config.DBConfig{
		User:          username,
		Password:      password,
		HostReadWrite: hostPort,
		HostReadOnly:  hostPort,
		Name:          databaseName,
		Schema:        schemaName,
	}

	db, err := database.NewDB(ctx, config)
	if err != nil {
		panic("Failed to connect to test database: " + err.Error())
	}

	// Clean up any development data that was seeded by migrations
	// This ensures tests start with a clean database
	cleanAllTestData(db)

	sharedContainer = pgContainer
	sharedDB = db
}

// Helper function to get shared test database
func DB(t testing.TB) *database.DB {
	t.Helper()

	sharedOnce.Do(setupSharedContainer)

	sharedMutex.RLock()
	defer sharedMutex.RUnlock()

	// Clean up test data BEFORE each test to ensure isolation
	// This prevents previous test data from affecting the current test
	cleanTestData(t)

	// Also clean up after the test completes
	t.Cleanup(func() {
		cleanTestData(t)
	})

	// Return a new connection from the same database pool
	// This ensures each test gets its own connection but shares the container
	return sharedDB
}

// cleanAllTestData removes all data from database tables (used during setup)
func cleanAllTestData(db *database.DB) {
	ctx := context.Background()

	// Delete data in dependency order to avoid foreign key constraints
	tables := []string{
		"log_entry_tags",
		"log_entries",
		"projects",
		"generated_insights",
		"tasks",
		"scheduled_deletions",
		"refresh_token_denylist",
		"user_sessions",
		"tags",
		"users",
	}

	for _, table := range tables {
		_, err := db.RDBMS().Exec(ctx, "DELETE FROM "+table)
		if err != nil {
			// Ignore errors during initial cleanup
			continue
		}
	}
}

// cleanTestData removes all test data from the database tables
func cleanTestData(t testing.TB) {
	if sharedDB == nil {
		return
	}

	ctx := context.Background()

	// Delete data in dependency order to avoid foreign key constraints
	tables := []string{
		"log_entry_tags",
		"log_entries",
		"projects",
		"generated_insights",
		"tasks",
		"scheduled_deletions",
		"refresh_token_denylist",
		"user_sessions",
		"tags",
		"users",
	}

	for _, table := range tables {
		_, err := sharedDB.RDBMS().Exec(ctx, "DELETE FROM "+table)
		if err != nil {
			t.Logf("Warning: failed to clean table %s: %v", table, err)
		}
	}

	// Note: We don't reset sequences since they don't affect test correctness
	// and the sequence names may vary depending on the migration strategy used
} // CleanupSharedResources should be called in TestMain to cleanup shared resources
func CleanupSharedResources() {
	sharedMutex.Lock()
	defer sharedMutex.Unlock()

	if sharedDB != nil {
		sharedDB.CloseDB()
		sharedDB = nil
	}

	if sharedContainer != nil {
		ctx := context.Background()
		if err := sharedContainer.Terminate(ctx); err != nil {
			// Log error but don't panic during cleanup
			panic("Failed to terminate shared container: " + err.Error())
		}
	}
}
