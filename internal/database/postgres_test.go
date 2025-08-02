//go:build integration
// +build integration

package database_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/garnizeh/englog/internal/config"
	"github.com/garnizeh/englog/internal/database"
	"github.com/garnizeh/englog/internal/store"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	testTimezone = pgtype.Text{String: "UTC", Valid: true}
)

// TestNewDB tests database initialization and migration
func TestNewDB(t *testing.T) {
	ctx := context.Background()

	// Setup test database container
	pgContainer, err := postgres.Run(
		ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		postgres.BasicWaitStrategies(),
		postgres.WithInitScripts("../../scripts/init-postgres.sql"),
	)
	require.NoError(t, err)
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	}()

	// Get connection configuration
	host, err := pgContainer.Host(ctx)
	require.NoError(t, err)
	port, err := pgContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	hostPort := host + ":" + port.Port()

	t.Run("SuccessfulConnection", func(t *testing.T) {
		config := config.DBConfig{
			User:          "testuser",
			Password:      "testpass",
			HostReadWrite: hostPort,
			HostReadOnly:  hostPort,
			Name:          "testdb",
			Schema:        "public",
		}

		db, err := database.NewDB(ctx, config)
		require.NoError(t, err)
		require.NotNil(t, db)
		defer db.CloseDB()

		// Verify schema was created and migrations ran
		err = db.Check(ctx)
		require.NoError(t, err)
	})

	t.Run("InvalidConnectionParams", func(t *testing.T) {
		config := config.DBConfig{
			User:          "invaliduser",
			Password:      "invalidpass",
			HostReadWrite: "invalid:5432",
			HostReadOnly:  "invalid:5432",
			Name:          "invaliddb",
			Schema:        "public",
		}

		db, err := database.NewDB(ctx, config)
		require.Error(t, err)
		require.Nil(t, db)
	})

	t.Run("WithCustomSchema", func(t *testing.T) {
		config := config.DBConfig{
			User:          "testuser",
			Password:      "testpass",
			HostReadWrite: hostPort,
			HostReadOnly:  hostPort,
			Name:          "testdb",
			Schema:        "test",
		}

		db, err := database.NewDB(ctx, config)
		require.NoError(t, err)
		require.NotNil(t, db)
		defer db.CloseDB()

		// Verify schema was created and migrations ran
		err = db.Check(ctx)
		require.NoError(t, err)
	})
}

// TestDBOperations tests read/write operations
func TestDBOperations(t *testing.T) {
	ctx := context.Background()

	// Setup test database
	db := setupTestDB(t)
	defer db.CloseDB()

	t.Run("WriteOperation", func(t *testing.T) {
		var userID uuid.UUID

		err := db.Write(ctx, func(qtx *store.Queries) error {
			user, err := qtx.CreateUser(ctx, store.CreateUserParams{
				Email:        "write-test@example.com",
				PasswordHash: "hashedpassword",
				FirstName:    "Write",
				LastName:     "Test",
				Timezone:     testTimezone,
				Preferences:  []byte(`{"theme": "dark"}`),
			})
			if err != nil {
				return err
			}
			userID = user.ID
			return nil
		})
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, userID)
	})

	t.Run("ReadOperation", func(t *testing.T) {
		// First create a user
		var createdUserID uuid.UUID
		err := db.Write(ctx, func(qtx *store.Queries) error {
			user, err := qtx.CreateUser(ctx, store.CreateUserParams{
				Email:        "read-test@example.com",
				PasswordHash: "hashedpassword",
				FirstName:    "Read",
				LastName:     "Test",
				Timezone:     testTimezone,
				Preferences:  []byte(`{"theme": "light"}`),
			})
			if err != nil {
				return err
			}
			createdUserID = user.ID
			return nil
		})
		require.NoError(t, err)

		// Then read it back
		var fetchedUser store.User
		err = db.Read(ctx, func(qtx *store.Queries) error {
			user, err := qtx.GetUserByID(ctx, createdUserID)
			if err != nil {
				return err
			}
			fetchedUser = user
			return nil
		})
		require.NoError(t, err)
		require.Equal(t, createdUserID, fetchedUser.ID)
		require.Equal(t, "read-test@example.com", fetchedUser.Email)
		require.Equal(t, "Read", fetchedUser.FirstName)
		require.Equal(t, "Test", fetchedUser.LastName)
	})

	t.Run("TransactionRollback", func(t *testing.T) {
		initialCount := countUsers(t, db)

		// Attempt to create user but force an error
		err := db.Write(ctx, func(qtx *store.Queries) error {
			_, err := qtx.CreateUser(ctx, store.CreateUserParams{
				Email:        "rollback-test@example.com",
				PasswordHash: "hashedpassword",
				FirstName:    "Rollback",
				LastName:     "Test",
				Timezone:     testTimezone,
				Preferences:  []byte(`{}`),
			})
			if err != nil {
				return err
			}
			// Force an error to trigger rollback
			return fmt.Errorf("forced error for rollback test")
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "forced error for rollback test")

		// Verify rollback worked - user count should be the same
		finalCount := countUsers(t, db)
		require.Equal(t, initialCount, finalCount)
	})

	t.Run("ReadOnlyTransactionConstraint", func(t *testing.T) {
		// Attempt to write in a read transaction should fail
		err := db.Read(ctx, func(qtx *store.Queries) error {
			_, err := qtx.CreateUser(ctx, store.CreateUserParams{
				Email:        "readonly-fail@example.com",
				PasswordHash: "hashedpassword",
				FirstName:    "ReadOnly",
				LastName:     "Fail",
				Timezone:     testTimezone,
				Preferences:  []byte(`{}`),
			})
			return err
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "read-only")
	})
}

// TestNoRows tests the NoRows utility function
func TestNoRows(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	defer db.CloseDB()

	t.Run("NoRowsError", func(t *testing.T) {
		var err error
		readErr := db.Read(ctx, func(qtx *store.Queries) error {
			// Try to get a user that doesn't exist
			_, err = qtx.GetUserByID(ctx, uuid.New())
			return err
		})

		require.Error(t, readErr)
		require.True(t, database.NoRows(readErr))
	})

	t.Run("OtherError", func(t *testing.T) {
		otherErr := fmt.Errorf("some other error")
		require.False(t, database.NoRows(otherErr))
	})

	t.Run("NilError", func(t *testing.T) {
		require.False(t, database.NoRows(nil))
	})
}

// TestDBCheck tests the database health check functionality
func TestDBCheck(t *testing.T) {
	ctx := context.Background()

	t.Run("HealthyDatabase", func(t *testing.T) {
		db := setupTestDB(t)
		defer db.CloseDB()

		err := db.Check(ctx)
		require.NoError(t, err)
	})

	t.Run("CheckWithTimeout", func(t *testing.T) {
		db := setupTestDB(t)
		defer db.CloseDB()

		// Create a context with timeout
		timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		err := db.Check(timeoutCtx)
		require.NoError(t, err)
	})
}

// TestRDBMS tests the RDBMS getter
func TestRDBMS(t *testing.T) {
	db := setupTestDB(t)
	defer db.CloseDB()

	pool := db.RDBMS()
	require.NotNil(t, pool)

	// Verify the pool is functional
	ctx := context.Background()
	err := pool.Ping(ctx)
	require.NoError(t, err)
}

// TestConcurrentOperations tests concurrent database operations
func TestConcurrentOperations(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	defer db.CloseDB()

	const numGoroutines = 10
	const usersPerGoroutine = 5

	errChan := make(chan error, numGoroutines)
	userChan := make(chan uuid.UUID, numGoroutines*usersPerGoroutine)

	// Start multiple goroutines creating users concurrently
	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			for j := 0; j < usersPerGoroutine; j++ {
				err := db.Write(ctx, func(qtx *store.Queries) error {
					user, err := qtx.CreateUser(ctx, store.CreateUserParams{
						Email:        fmt.Sprintf("concurrent-%d-%d@example.com", goroutineID, j),
						PasswordHash: "hashedpassword",
						FirstName:    fmt.Sprintf("User%d", goroutineID),
						LastName:     fmt.Sprintf("Test%d", j),
						Timezone:     testTimezone,
						Preferences:  []byte(`{}`),
					})
					if err != nil {
						return err
					}
					userChan <- user.ID
					return nil
				})
				if err != nil {
					errChan <- err
					return
				}
			}
			errChan <- nil
		}(i)
	}

	// Collect results
	var createdUsers []uuid.UUID
	for i := 0; i < numGoroutines; i++ {
		err := <-errChan
		require.NoError(t, err)
	}

	// Collect all user IDs
	close(userChan)
	for userID := range userChan {
		createdUsers = append(createdUsers, userID)
	}

	require.Len(t, createdUsers, numGoroutines*usersPerGoroutine)

	// Verify all users were created by reading them back
	for _, userID := range createdUsers {
		err := db.Read(ctx, func(qtx *store.Queries) error {
			_, err := qtx.GetUserByID(ctx, userID)
			return err
		})
		require.NoError(t, err)
	}
}

// TestComplexTransactions tests more complex transaction scenarios
func TestComplexTransactions(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	defer db.CloseDB()

	t.Run("MultipleOperationsInTransaction", func(t *testing.T) {
		var userID uuid.UUID
		var projectID uuid.UUID
		var logEntryID uuid.UUID

		err := db.Write(ctx, func(qtx *store.Queries) error {
			// Create user
			user, err := qtx.CreateUser(ctx, store.CreateUserParams{
				Email:        "complex-tx@example.com",
				PasswordHash: "hashedpassword",
				FirstName:    "Complex",
				LastName:     "Transaction",
				Timezone:     testTimezone,
				Preferences:  []byte(`{}`),
			})
			if err != nil {
				return err
			}
			userID = user.ID

			// Create project
			project, err := qtx.CreateProject(ctx, store.CreateProjectParams{
				Name:        "Complex Project",
				Description: pgtype.Text{String: "A complex project", Valid: true},
				Color:       pgtype.Text{String: "#FF5722", Valid: true},
				Status:      pgtype.Text{String: "active", Valid: true},
				StartDate:   pgtype.Date{Time: time.Now(), Valid: true},
				EndDate:     pgtype.Date{},
				CreatedBy:   userID,
				IsDefault:   pgtype.Bool{Bool: false, Valid: true},
			})
			if err != nil {
				return err
			}
			projectID = project.ID

			// Create log entry
			logEntry, err := qtx.CreateLogEntry(ctx, store.CreateLogEntryParams{
				UserID:      userID,
				ProjectID:   pgtype.UUID{Bytes: projectID, Valid: true},
				Title:       "Complex Log Entry",
				Description: pgtype.Text{String: "A complex log entry", Valid: true},
				Type:        "development",
				StartTime:   pgtype.Timestamptz{Time: time.Now().Add(-time.Hour), Valid: true},
				EndTime:     pgtype.Timestamptz{Time: time.Now(), Valid: true},
				ValueRating: "high",
				ImpactLevel: "team",
			})
			if err != nil {
				return err
			}
			logEntryID = logEntry.ID

			return nil
		})
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, userID)
		require.NotEqual(t, uuid.Nil, projectID)
		require.NotEqual(t, uuid.Nil, logEntryID)

		// Verify all entities were created
		err = db.Read(ctx, func(qtx *store.Queries) error {
			_, err := qtx.GetUserByID(ctx, userID)
			if err != nil {
				return err
			}
			_, err = qtx.GetProjectByID(ctx, projectID)
			if err != nil {
				return err
			}
			_, err = qtx.GetLogEntryByID(ctx, logEntryID)
			return err
		})
		require.NoError(t, err)
	})
}

// Helper function to setup a test database
func setupTestDB(t *testing.T) *database.DB {
	t.Helper()
	ctx := context.Background()

	// Setup test database container
	pgContainer, err := postgres.Run(
		ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		postgres.WithInitScripts("../../scripts/init-postgres.sql"),
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
	port, err := pgContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	hostPort := host + ":" + port.Port()

	config := config.DBConfig{
		User:          "testuser",
		Password:      "testpass",
		HostReadWrite: hostPort,
		HostReadOnly:  hostPort,
		Name:          "testdb",
		Schema:        "test",
	}

	db, err := database.NewDB(ctx, config)
	require.NoError(t, err)
	require.NotNil(t, db)

	return db
}

// Helper function to count users in the database
func countUsers(t *testing.T, db *database.DB) int {
	t.Helper()
	ctx := context.Background()

	var count int64
	err := db.Read(ctx, func(qtx *store.Queries) error {
		userCount, err := qtx.GetUserCount(ctx)
		if err != nil {
			return err
		}
		count = userCount
		return nil
	})
	require.NoError(t, err)

	return int(count)
}

// BenchmarkDatabaseOperations benchmarks database operations
func BenchmarkDatabaseOperations(b *testing.B) {
	ctx := context.Background()
	db := setupBenchmarkDB(b)
	defer db.CloseDB()

	b.Run("CreateUser", func(b *testing.B) {
		b.ResetTimer()
		for i := range b.N {
			err := db.Write(ctx, func(qtx *store.Queries) error {
				_, err := qtx.CreateUser(ctx, store.CreateUserParams{
					Email:        fmt.Sprintf("benchmark-user-%d-%d@example.com", time.Now().UnixNano(), i),
					PasswordHash: "hashedpassword",
					FirstName:    "Benchmark",
					LastName:     "User",
					Timezone:     testTimezone,
					Preferences:  []byte(`{}`),
				})
				return err
			})
			if err != nil {
				b.Fatalf("Failed to create user: %v", err)
			}
		}
	})

	b.Run("ReadUser", func(b *testing.B) {
		// Create a user first
		var userID uuid.UUID
		err := db.Write(ctx, func(qtx *store.Queries) error {
			user, err := qtx.CreateUser(ctx, store.CreateUserParams{
				Email:        fmt.Sprintf("benchmark-read-%d@example.com", time.Now().UnixNano()),
				PasswordHash: "hashedpassword",
				FirstName:    "Benchmark",
				LastName:     "Read",
				Timezone:     testTimezone,
				Preferences:  []byte(`{}`),
			})
			if err != nil {
				return err
			}
			userID = user.ID
			return nil
		})
		if err != nil {
			b.Fatalf("Failed to create user for benchmark: %v", err)
		}

		b.ResetTimer()
		for range b.N {
			err := db.Read(ctx, func(qtx *store.Queries) error {
				_, err := qtx.GetUserByID(ctx, userID)
				return err
			})
			if err != nil {
				b.Fatalf("Failed to read user: %v", err)
			}
		}
	})
}

// Helper function for benchmarks
func setupBenchmarkDB(b *testing.B) *database.DB {
	b.Helper()
	ctx := context.Background()

	// For benchmarks, we'll use a simpler setup
	// In a real scenario, you might want to use a persistent test database
	pgContainer, err := postgres.Run(
		ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("benchdb"),
		postgres.WithUsername("benchuser"),
		postgres.WithPassword("benchpass"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		b.Fatalf("Failed to start postgres container: %v", err)
	}

	b.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			b.Logf("Failed to terminate container: %v", err)
		}
	})

	host, err := pgContainer.Host(ctx)
	if err != nil {
		b.Fatalf("Failed to get container host: %v", err)
	}
	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		b.Fatalf("Failed to get container port: %v", err)
	}

	hostPort := host + ":" + port.Port()

	config := config.DBConfig{
		User:          "benchuser",
		Password:      "benchpass",
		HostReadWrite: hostPort,
		HostReadOnly:  hostPort,
		Name:          "benchdb",
		Schema:        "public",
	}

	db, err := database.NewDB(ctx, config)
	if err != nil {
		b.Fatalf("Failed to create database: %v", err)
	}

	return db
}
