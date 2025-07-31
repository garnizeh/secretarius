//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/garnizeh/englog/internal/config"
	"github.com/garnizeh/englog/internal/database"
	"github.com/garnizeh/englog/internal/logging"
	"github.com/garnizeh/englog/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// DatabaseTestSuite tests database operations using real PostgreSQL container
// "Testing is the process of evaluating a system to find discrepancies." 🔍
type DatabaseTestSuite struct {
	suite.Suite
	container *postgres.PostgresContainer
	db        *database.DB
	logger    *logging.Logger
}

func (suite *DatabaseTestSuite) SetupSuite() {
	ctx := context.Background()

	// Start PostgreSQL container
	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:17-alpine"),
		postgres.WithDatabase("test_englog"),
		postgres.WithUsername("test_user"),
		postgres.WithPassword("test_password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Minute)),
	)
	require.NoError(suite.T(), err)
	suite.container = container

	// Get connection info from container
	host, err := container.Host(ctx)
	require.NoError(suite.T(), err)

	port, err := container.MappedPort(ctx, "5432")
	require.NoError(suite.T(), err)

	// Initialize logger
	logConfig := config.LoggingConfig{
		Level:     slog.LevelDebug,
		Format:    "text",
		AddSource: true,
	}
	suite.logger = logging.NewLogger(logConfig)

	// Create config for database
	// Initialize database using the project's method
	dbConfig := config.DBConfig{
		HostReadWrite: fmt.Sprintf("%s:%s", host, port.Port()),
		HostReadOnly:  fmt.Sprintf("%s:%s", host, port.Port()),
		User:          "test_user",
		Password:      "test_password",
		Name:          "test_englog",
		Schema:        "englog_test", // Add missing schema field
	}
	db, err := database.NewDB(ctx, dbConfig)
	require.NoError(suite.T(), err)
	suite.db = db
}

func (suite *DatabaseTestSuite) TearDownSuite() {
	if suite.container != nil {
		ctx := context.Background()
		suite.container.Terminate(ctx)
	}
}

func (suite *DatabaseTestSuite) SetupTest() {
	// Clean up data before each test
	ctx := context.Background()
	_, err := suite.db.RDBMS().Exec(ctx, "TRUNCATE TABLE users CASCADE")
	require.NoError(suite.T(), err)
	_, err = suite.db.RDBMS().Exec(ctx, "TRUNCATE TABLE projects CASCADE")
	require.NoError(suite.T(), err)
	_, err = suite.db.RDBMS().Exec(ctx, "TRUNCATE TABLE log_entries CASCADE")
	require.NoError(suite.T(), err)
}

func (suite *DatabaseTestSuite) TestBasicDatabaseOperations() {
	ctx := context.Background()

	// Test database connectivity
	err := suite.db.RDBMS().Ping(ctx)
	require.NoError(suite.T(), err)

	// Test basic SQL execution
	var result int
	err = suite.db.RDBMS().QueryRow(ctx, "SELECT 1").Scan(&result)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, result)
}

func (suite *DatabaseTestSuite) TestDatabaseConnections() {
	ctx := context.Background()

	// Test RDBMS connection
	rdbmsConn := suite.db.RDBMS()
	assert.NotNil(suite.T(), rdbmsConn)
	err := rdbmsConn.Ping(ctx)
	require.NoError(suite.T(), err)
}

func (suite *DatabaseTestSuite) TestTransactionSupport() {
	ctx := context.Background()

	// Test Read transaction
	err := suite.db.Read(ctx, func(qtx *store.Queries) error {
		// This is a read-only transaction
		var count int64
		err := suite.db.RDBMS().QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
		return err
	})
	require.NoError(suite.T(), err)

	// Test Write transaction
	err = suite.db.Write(ctx, func(qtx *store.Queries) error {
		// This is a read-write transaction
		var count int64
		err := suite.db.RDBMS().QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
		return err
	})
	require.NoError(suite.T(), err)
}

func (suite *DatabaseTestSuite) TestDatabaseHealthCheck() {
	ctx := context.Background()

	// Test Check method
	err := suite.db.Check(ctx)
	require.NoError(suite.T(), err)
}

func TestDatabaseIntegration(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}
