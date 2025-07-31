//go:build integration
// +build integration

package store_test

import (
	"context"
	"database/sql"
	"encoding/json"
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
	timezoneUTC = pgtype.Text{String: "UTC", Valid: true}
)

// TestStoreTypes validates that SQLC generated types work correctly
func TestStoreTypes(t *testing.T) {
	t.Run("CreateUserParams", testCreateUserParams)
	t.Run("CreateProjectParams", testCreateProjectParams)
	t.Run("CreateLogEntryParams", testCreateLogEntryParams)
	t.Run("AuthParams", testAuthParams)
}

func testCreateUserParams(t *testing.T) {
	// Test that CreateUserParams struct can be created with correct types
	preferences := json.RawMessage(`{"theme": "dark", "language": "en"}`)

	params := store.CreateUserParams{
		Email:        "test@example.com",
		PasswordHash: "hashed_password_123",
		FirstName:    "John",
		LastName:     "Doe",
		Timezone:     timezoneUTC,
		Preferences:  preferences,
	}

	require.Equal(t, "test@example.com", params.Email)
	require.Equal(t, "hashed_password_123", params.PasswordHash)
	require.Equal(t, "John", params.FirstName)
	require.Equal(t, "Doe", params.LastName)
	require.NotNil(t, params.Timezone)
	require.Equal(t, timezoneUTC, params.Timezone)
	require.NotNil(t, params.Preferences)
}

func testCreateProjectParams(t *testing.T) {
	// Test CreateProjectParams with correct pgtype usage
	userID := uuid.New()
	description := pgtype.Text{String: "Test project description", Valid: true}
	color := pgtype.Text{String: "#FF5722", Valid: true}
	status := pgtype.Text{String: "active", Valid: true}
	isDefault := pgtype.Bool{Bool: true, Valid: true}

	now := time.Now()
	startDate := pgtype.Date{}
	err := startDate.Scan(now)
	require.NoError(t, err)

	endDate := pgtype.Date{} // NULL date

	params := store.CreateProjectParams{
		Name:        "Test Project",
		Description: description,
		Color:       color,
		Status:      status,
		StartDate:   startDate,
		EndDate:     endDate,
		CreatedBy:   userID,
		IsDefault:   isDefault,
	}

	require.Equal(t, "Test Project", params.Name)
	require.NotNil(t, params.Description)
	require.Equal(t, description, params.Description)
	require.NotNil(t, params.Color)
	require.Equal(t, color, params.Color)
	require.NotNil(t, params.Status)
	require.Equal(t, status, params.Status)
	require.True(t, params.StartDate.Valid)
	require.False(t, params.EndDate.Valid) // EndDate is NULL
	require.NotNil(t, params.IsDefault)
	require.True(t, params.IsDefault.Bool)
	require.Equal(t, userID, params.CreatedBy)
}

func testCreateLogEntryParams(t *testing.T) {
	// Test CreateLogEntryParams with pgtype timestamps
	userID := uuid.New()
	projectID := uuid.New()
	description := pgtype.Text{String: "Working on test implementation", Valid: true}

	startTime := pgtype.Timestamptz{}
	err := startTime.Scan(time.Now().Add(-2 * time.Hour))
	require.NoError(t, err)

	endTime := pgtype.Timestamptz{}
	err = endTime.Scan(time.Now().Add(-1 * time.Hour))
	require.NoError(t, err)

	projectUUID := pgtype.UUID{Bytes: projectID, Valid: true}

	params := store.CreateLogEntryParams{
		UserID:      userID,
		ProjectID:   projectUUID,
		Title:       "Test Development Work",
		Description: description,
		Type:        "development",
		StartTime:   startTime,
		EndTime:     endTime,
		ValueRating: "high",
		ImpactLevel: "team",
	}

	require.Equal(t, userID, params.UserID)
	require.Equal(t, "Test Development Work", params.Title)
	require.NotNil(t, params.Description)
	require.Equal(t, description, params.Description)
	require.Equal(t, "development", params.Type)
	require.True(t, params.StartTime.Valid)
	require.True(t, params.EndTime.Valid)
	require.Equal(t, "high", params.ValueRating)
	require.Equal(t, "team", params.ImpactLevel)
	require.True(t, params.ProjectID.Valid)
	require.Equal(t, projectID[:], params.ProjectID.Bytes[:])
}

func testAuthParams(t *testing.T) {
	// Test authentication-related parameters
	userID := uuid.New()
	tokenJTI := uuid.New().String()
	reason := pgtype.Text{String: "logout", Valid: true}

	expiresAt := pgtype.Timestamptz{}
	err := expiresAt.Scan(time.Now().Add(24 * time.Hour))
	require.NoError(t, err)

	denylistParams := store.CreateRefreshTokenDenylistParams{
		Jti:       tokenJTI,
		UserID:    userID,
		ExpiresAt: expiresAt,
		Column4:   reason,
	}

	require.Equal(t, tokenJTI, denylistParams.Jti)
	require.Equal(t, userID, denylistParams.UserID)
	require.NotNil(t, denylistParams.Column4)
	require.Equal(t, reason, denylistParams.Column4)
	require.True(t, denylistParams.ExpiresAt.Valid)
}

// TestModelValidation ensures that generated models have correct field types
func TestModelValidation(t *testing.T) {
	t.Run("UserModel", testUserModel)
	t.Run("ProjectModel", testProjectModel)
	t.Run("LogEntryModel", testLogEntryModel)
}

func testUserModel(t *testing.T) {
	// Test User model structure
	var user store.User

	// Verify field types exist and are accessible
	user.ID = uuid.New()
	user.Email = "test@example.com"
	user.PasswordHash = "hashed"
	user.FirstName = "John"
	user.LastName = "Doe"

	user.Timezone = timezoneUTC

	preferences := []byte(`{"theme": "dark"}`)
	user.Preferences = preferences

	// Basic validation
	require.NotEqual(t, uuid.Nil, user.ID)
	require.Equal(t, "test@example.com", user.Email)
	require.Equal(t, "hashed", user.PasswordHash)
	require.Equal(t, "John", user.FirstName)
	require.Equal(t, "Doe", user.LastName)
	require.NotNil(t, user.Timezone)
	require.Equal(t, timezoneUTC, user.Timezone)
	require.NotNil(t, user.Preferences)
}

func testProjectModel(t *testing.T) {
	// Test Project model structure
	var project store.Project

	project.ID = uuid.New()
	project.Name = "Test Project"
	project.CreatedBy = uuid.New()

	description := pgtype.Text{String: "Test description", Valid: true}
	project.Description = description

	color := pgtype.Text{String: "#FF5722", Valid: true}
	project.Color = color

	status := pgtype.Text{String: "active", Valid: true}
	project.Status = status

	isDefault := pgtype.Bool{Bool: true, Valid: true}
	project.IsDefault = isDefault

	// Basic validation
	require.NotEqual(t, uuid.Nil, project.ID)
	require.Equal(t, "Test Project", project.Name)
	require.NotEqual(t, uuid.Nil, project.CreatedBy)
	require.NotNil(t, project.Description)
	require.Equal(t, description, project.Description)
	require.NotNil(t, project.Color)
	require.Equal(t, color, project.Color)
	require.NotNil(t, project.Status)
	require.Equal(t, status, project.Status)
	require.NotNil(t, project.IsDefault)
	require.True(t, project.IsDefault.Bool)
}

func testLogEntryModel(t *testing.T) {
	// Test LogEntry model structure
	var entry store.LogEntry

	entry.ID = uuid.New()
	entry.UserID = uuid.New()
	entry.Title = "Test Entry"
	entry.Type = "development"
	entry.ValueRating = "high"
	entry.ImpactLevel = "team"

	description := pgtype.Text{String: "Test description", Valid: true}
	entry.Description = description

	// Test pgtype fields
	err := entry.StartTime.Scan(time.Now())
	require.NoError(t, err)

	err = entry.EndTime.Scan(time.Now().Add(time.Hour))
	require.NoError(t, err)

	projectID := uuid.New()
	entry.ProjectID = pgtype.UUID{Bytes: projectID, Valid: true}

	// Basic validation
	require.NotEqual(t, uuid.Nil, entry.ID)
	require.NotEqual(t, uuid.Nil, entry.UserID)
	require.Equal(t, "Test Entry", entry.Title)
	require.Equal(t, "development", entry.Type)
	require.Equal(t, "high", entry.ValueRating)
	require.Equal(t, "team", entry.ImpactLevel)
	require.NotNil(t, entry.Description)
	require.Equal(t, description, entry.Description)
	require.True(t, entry.StartTime.Valid)
	require.True(t, entry.EndTime.Valid)
	require.True(t, entry.ProjectID.Valid)
	require.Equal(t, projectID[:], entry.ProjectID.Bytes[:])
}

// BenchmarkStoreStructs benchmarks struct creation and manipulation
func BenchmarkStoreStructs(b *testing.B) {
	b.Run("CreateUserParams", benchmarkCreateUserParams)
	b.Run("CreateProjectParams", benchmarkCreateProjectParams)
	b.Run("CreateLogEntryParams", benchmarkCreateLogEntryParams)
}

func benchmarkCreateUserParams(b *testing.B) {
	preferences := json.RawMessage(`{"theme": "dark"}`)

	for i := 0; i < b.N; i++ {
		_ = store.CreateUserParams{
			Email:        "test@example.com",
			PasswordHash: "hashed",
			FirstName:    "John",
			LastName:     "Doe",
			Timezone:     timezoneUTC,
			Preferences:  preferences,
		}
	}
}

func benchmarkCreateProjectParams(b *testing.B) {
	userID := uuid.New()
	description := pgtype.Text{String: "Test project", Valid: true}
	color := pgtype.Text{String: "#FF5722", Valid: true}
	status := pgtype.Text{String: "active", Valid: true}
	isDefault := pgtype.Bool{Bool: true, Valid: true}

	for i := 0; i < b.N; i++ {
		_ = store.CreateProjectParams{
			Name:        "Test Project",
			Description: description,
			Color:       color,
			Status:      status,
			StartDate:   pgtype.Date{},
			EndDate:     pgtype.Date{},
			CreatedBy:   userID,
			IsDefault:   isDefault,
		}
	}
}

func benchmarkCreateLogEntryParams(b *testing.B) {
	userID := uuid.New()
	projectID := pgtype.UUID{Bytes: uuid.New(), Valid: true}
	description := pgtype.Text{String: "Test work", Valid: true}

	startTime := pgtype.Timestamptz{}
	_ = startTime.Scan(time.Now())
	endTime := pgtype.Timestamptz{}
	_ = endTime.Scan(time.Now().Add(time.Hour))

	for i := 0; i < b.N; i++ {
		_ = store.CreateLogEntryParams{
			UserID:      userID,
			ProjectID:   projectID,
			Title:       "Test Entry",
			Description: description,
			Type:        "development",
			StartTime:   startTime,
			EndTime:     endTime,
			ValueRating: "high",
			ImpactLevel: "team",
		}
	}
}
func TestStore(t *testing.T) {
	ctx := context.Background()

	// Setup test database container
	pgContainer, err := postgres.Run(
		ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		postgres.BasicWaitStrategies(),
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

	// Setup database connection
	config := config.DBConfig{
		User:          "testuser",
		Password:      "testpass",
		HostReadWrite: host + ":" + port.Port(),
		HostReadOnly:  host + ":" + port.Port(),
		Name:          "testdb",
		Schema:        "testschema",
	}

	db, err := database.NewDB(ctx, config)
	require.NoError(t, err)
	defer db.CloseDB()

	// Run tests
	t.Run("CreateAndGetUser", testCreateAndGetUser(db))
	t.Run("CreateProject", testCreateProject(db))
	t.Run("CreateLogEntry", testCreateLogEntry(db))
	t.Run("UserAuthentication", testUserAuthentication(db))
	t.Run("TransactionRollback", testTransactionRollback(db))
}

func testCreateAndGetUser(db *database.DB) func(*testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()

		var createdUser store.User
		var fetchedUser store.User

		// Test user creation
		err := db.Write(ctx, func(qtx *store.Queries) error {
			preferences := json.RawMessage(`{"theme": "dark", "language": "en"}`)

			user, err := qtx.CreateUser(ctx, store.CreateUserParams{
				Email:        "test@example.com",
				PasswordHash: "hashed_password_123",
				FirstName:    "John",
				LastName:     "Doe",
				Timezone:     timezoneUTC,
				Preferences:  preferences,
			})
			if err != nil {
				return err
			}
			createdUser = user
			return nil
		})
		require.NoError(t, err)
		require.NotEmpty(t, createdUser.ID)
		require.Equal(t, "test@example.com", createdUser.Email)
		require.Equal(t, "John", createdUser.FirstName)
		require.Equal(t, "Doe", createdUser.LastName)

		// Test user retrieval by ID
		err = db.Read(ctx, func(qtx *store.Queries) error {
			user, err := qtx.GetUserByID(ctx, createdUser.ID)
			if err != nil {
				return err
			}
			fetchedUser = user
			return nil
		})
		require.NoError(t, err)
		require.Equal(t, createdUser.ID, fetchedUser.ID)
		require.Equal(t, createdUser.Email, fetchedUser.Email)

		// Test user retrieval by email
		err = db.Read(ctx, func(qtx *store.Queries) error {
			user, err := qtx.GetUserByEmail(ctx, "test@example.com")
			if err != nil {
				return err
			}
			fetchedUser = user
			return nil
		})
		require.NoError(t, err)
		require.Equal(t, createdUser.ID, fetchedUser.ID)
	}
}

func testCreateProject(db *database.DB) func(*testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()

		// First create a user
		var userID uuid.UUID
		err := db.Write(ctx, func(qtx *store.Queries) error {
			user, err := qtx.CreateUser(ctx, store.CreateUserParams{
				Email:        "project-owner@example.com",
				PasswordHash: "hashed_password",
				FirstName:    "Project",
				LastName:     "Owner",
				Timezone:     timezoneUTC,
				Preferences:  json.RawMessage(`{}`),
			})
			if err != nil {
				return err
			}
			userID = user.ID
			return nil
		})
		require.NoError(t, err)

		// Create a project
		var project store.Project
		status := pgtype.Text{String: "active", Valid: true}
		err = db.Write(ctx, func(qtx *store.Queries) error {
			description := pgtype.Text{String: "Test project description", Valid: true}
			color := pgtype.Text{String: "#FF5722", Valid: true}
			startDate := pgtype.Date{Time: time.Now(), Valid: true}
			endDate := pgtype.Date{} // NULL date
			isDefault := pgtype.Bool{Bool: true, Valid: true}

			proj, err := qtx.CreateProject(ctx, store.CreateProjectParams{
				Name:        "Test Project",
				Description: description,
				Color:       color,
				Status:      status,
				StartDate:   startDate,
				EndDate:     endDate,
				CreatedBy:   userID,
				IsDefault:   isDefault,
			})
			if err != nil {
				return err
			}
			project = proj
			return nil
		})
		require.NoError(t, err)
		require.NotEmpty(t, project.ID)
		require.Equal(t, "Test Project", project.Name)
		require.Equal(t, status, project.Status)
		require.True(t, project.IsDefault.Bool)
	}
}

func testCreateLogEntry(db *database.DB) func(*testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()

		// Setup user and project
		var userID uuid.UUID
		var projectID uuid.UUID

		err := db.Write(ctx, func(qtx *store.Queries) error {
			// Create user
			user, err := qtx.CreateUser(ctx, store.CreateUserParams{
				Email:        "logger@example.com",
				PasswordHash: "hashed_password",
				FirstName:    "Logger",
				LastName:     "User",
				Timezone:     timezoneUTC,
				Preferences:  json.RawMessage(`{}`),
			})
			if err != nil {
				return err
			}
			userID = user.ID

			// Create project
			project, err := qtx.CreateProject(ctx, store.CreateProjectParams{
				Name:        "Log Project",
				Description: pgtype.Text{String: "Project for logging", Valid: true},
				Color:       pgtype.Text{String: "#2196F3", Valid: true},
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
			return nil
		})
		require.NoError(t, err)

		// Create log entry
		startTime := time.Now().Add(-2 * time.Hour)
		endTime := time.Now().Add(-1 * time.Hour)

		var logEntry store.LogEntry
		err = db.Write(ctx, func(qtx *store.Queries) error {
			entry, err := qtx.CreateLogEntry(ctx, store.CreateLogEntryParams{
				UserID:      userID,
				ProjectID:   pgtype.UUID{Bytes: projectID, Valid: true},
				Title:       "Test Development Work",
				Description: pgtype.Text{String: "Working on test implementation", Valid: true},
				Type:        "development",
				StartTime:   pgtype.Timestamptz{Time: startTime, Valid: true},
				EndTime:     pgtype.Timestamptz{Time: endTime, Valid: true},
				ValueRating: "high",
				ImpactLevel: "team",
			})
			if err != nil {
				return err
			}
			logEntry = entry
			return nil
		})
		require.NoError(t, err)
		require.NotEmpty(t, logEntry.ID)
		require.Equal(t, "Test Development Work", logEntry.Title)
		require.Equal(t, "development", logEntry.Type)
		require.Equal(t, "high", logEntry.ValueRating)
	}
}

func testUserAuthentication(db *database.DB) func(*testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()

		// Create user
		var userID uuid.UUID
		err := db.Write(ctx, func(qtx *store.Queries) error {
			user, err := qtx.CreateUser(ctx, store.CreateUserParams{
				Email:        "auth-user@example.com",
				PasswordHash: "secure_hash",
				FirstName:    "Auth",
				LastName:     "User",
				Timezone:     timezoneUTC,
				Preferences:  json.RawMessage(`{}`),
			})
			if err != nil {
				return err
			}
			userID = user.ID
			return nil
		})
		require.NoError(t, err)

		// Test token denylisting
		tokenJTI := uuid.New().String()
		err = db.Write(ctx, func(qtx *store.Queries) error {
			return qtx.CreateRefreshTokenDenylist(ctx, store.CreateRefreshTokenDenylistParams{
				Jti:       tokenJTI,
				UserID:    userID,
				ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(24 * time.Hour), Valid: true},
				Column4:   "logout",
			})
		})
		require.NoError(t, err)

		// Test token denylist check
		var isDenylisted bool
		err = db.Read(ctx, func(qtx *store.Queries) error {
			result, err := qtx.IsRefreshTokenDenylisted(ctx, tokenJTI)
			if err != nil {
				return err
			}
			isDenylisted = result
			return nil
		})
		require.NoError(t, err)
		require.True(t, isDenylisted)
	}
}

func testTransactionRollback(db *database.DB) func(*testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()

		// Test that transaction rollback works correctly
		err := db.Write(ctx, func(qtx *store.Queries) error {
			// Create a user
			_, err := qtx.CreateUser(ctx, store.CreateUserParams{
				Email:        "rollback-test@example.com",
				PasswordHash: "password",
				FirstName:    "Rollback",
				LastName:     "Test",
				Timezone:     timezoneUTC,
				Preferences:  json.RawMessage(`{}`),
			})
			if err != nil {
				return err
			}

			// Force an error to trigger rollback
			return sql.ErrTxDone
		})
		require.Error(t, err)

		// Verify the user was not created due to rollback
		err = db.Read(ctx, func(qtx *store.Queries) error {
			_, err := qtx.GetUserByEmail(ctx, "rollback-test@example.com")
			return err
		})
		require.Error(t, err) // Should fail because user wasn't created
		require.True(t, database.NoRows(err))
	}
}
