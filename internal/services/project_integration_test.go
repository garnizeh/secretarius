//go:build integration
// +build integration

package services_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/garnizeh/englog/internal/logging"
	"github.com/garnizeh/englog/internal/models"
	"github.com/garnizeh/englog/internal/services"
	"github.com/garnizeh/englog/internal/store/testutils"
)

// TestProjectServiceIntegration tests the full project management flow with database
func TestProjectServiceIntegration(t *testing.T) {
	db := testutils.DB(t)

	// Create test logger
	testLogger := logging.NewTestLogger()

	projectService := services.NewProjectService(db, testLogger)
	userService := services.NewUserService(db, testLogger)

	// Create a test user first (required for projects)
	ctx := context.Background()
	testUser, err := userService.CreateUser(ctx, &models.UserRegistration{
		Email:     fmt.Sprintf("project-user-%d@example.com", time.Now().UnixNano()),
		Password:  "password123",
		FirstName: "Project",
		LastName:  "User",
		Timezone:  "UTC",
	})
	require.NoError(t, err)

	defer func() {
		_ = userService.DeleteUser(ctx, testUser.ID.String())
	}()

	t.Run("FullProjectLifecycle", func(t *testing.T) {
		// Create a new project
		createReq := &models.ProjectRequest{
			Name:        fmt.Sprintf("integration-project-%d", time.Now().UnixNano()),
			Description: stringPtr("Integration test project"),
			Color:       "#FF5733",
			Status:      models.ProjectActive,
			IsDefault:   false,
		}

		// Test project creation
		createdProject, err := projectService.CreateProject(ctx, testUser.ID.String(), createReq)
		require.NoError(t, err)
		require.NotNil(t, createdProject)
		assert.Equal(t, createReq.Name, createdProject.Name)
		assert.Equal(t, *createReq.Description, *createdProject.Description)
		assert.Equal(t, createReq.Color, createdProject.Color)
		assert.Equal(t, createReq.Status, createdProject.Status)
		assert.Equal(t, testUser.ID, createdProject.CreatedBy)
		assert.NotEqual(t, uuid.Nil, createdProject.ID)

		// Test project retrieval by ID
		retrievedProject, err := projectService.GetProject(ctx, testUser.ID.String(), createdProject.ID.String())
		require.NoError(t, err)
		assert.Equal(t, createdProject.ID, retrievedProject.ID)
		assert.Equal(t, createdProject.Name, retrievedProject.Name)
		assert.Equal(t, createdProject.Color, retrievedProject.Color)

		// Test project update
		updateReq := &models.ProjectRequest{
			Name:        createdProject.Name + "-updated",
			Description: stringPtr("Updated description"),
			Color:       "#00FF00",
			Status:      models.ProjectCompleted,
			IsDefault:   true, // Note: IsDefault is not updated by UpdateProject
		}

		updatedProject, err := projectService.UpdateProject(ctx, testUser.ID.String(), createdProject.ID.String(), updateReq)
		require.NoError(t, err)
		assert.Equal(t, updateReq.Name, updatedProject.Name)
		assert.Equal(t, *updateReq.Description, *updatedProject.Description)
		assert.Equal(t, updateReq.Color, updatedProject.Color)
		assert.Equal(t, updateReq.Status, updatedProject.Status)
		// Note: IsDefault is not updated by UpdateProject method, so we don't test it here

		// Test project deletion
		err = projectService.DeleteProject(ctx, testUser.ID.String(), createdProject.ID.String())
		require.NoError(t, err)

		// Verify project is deleted
		_, err = projectService.GetProject(ctx, testUser.ID.String(), createdProject.ID.String())
		assert.Error(t, err)
	})

	t.Run("ProjectListingAndFiltering", func(t *testing.T) {
		// Create multiple projects for testing
		testProjects := []*models.ProjectRequest{
			{
				Name:        fmt.Sprintf("active-project-%d", time.Now().UnixNano()),
				Description: stringPtr("Active project"),
				Color:       "#FF0000",
				Status:      models.ProjectActive,
				IsDefault:   false,
			},
			{
				Name:        fmt.Sprintf("completed-project-%d", time.Now().UnixNano()),
				Description: stringPtr("Completed project"),
				Color:       "#00FF00",
				Status:      models.ProjectCompleted,
				IsDefault:   false,
			},
			{
				Name:        fmt.Sprintf("onhold-project-%d", time.Now().UnixNano()),
				Description: stringPtr("On hold project"),
				Color:       "#0000FF",
				Status:      models.ProjectOnHold,
				IsDefault:   true,
			},
		}

		var createdProjectIDs []string
		for _, projectReq := range testProjects {
			project, err := projectService.CreateProject(ctx, testUser.ID.String(), projectReq)
			require.NoError(t, err)
			createdProjectIDs = append(createdProjectIDs, project.ID.String())
		}

		// Clean up created projects
		defer func() {
			for _, projectID := range createdProjectIDs {
				_ = projectService.DeleteProject(ctx, testUser.ID.String(), projectID)
			}
		}()

		// Test GetUserProjects
		allProjects, err := projectService.GetUserProjects(ctx, testUser.ID.String())
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(allProjects), 3)

		// Test GetActiveProjects
		activeProjects, err := projectService.GetActiveProjects(ctx, testUser.ID.String())
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(activeProjects), 1)

		// Verify active projects only contain active status
		for _, project := range activeProjects {
			assert.Equal(t, models.ProjectActive, project.Status)
		}

		// Test GetDefaultProject
		defaultProject, err := projectService.GetDefaultProject(ctx, testUser.ID.String())
		require.NoError(t, err)
		assert.True(t, defaultProject.IsDefault)
		assert.Equal(t, models.ProjectOnHold, defaultProject.Status)
	})

	t.Run("ProjectOwnershipValidation", func(t *testing.T) {
		// Create another user
		otherUser, err := userService.CreateUser(ctx, &models.UserRegistration{
			Email:     fmt.Sprintf("other-user-%d@example.com", time.Now().UnixNano()),
			Password:  "password123",
			FirstName: "Other",
			LastName:  "User",
			Timezone:  "UTC",
		})
		require.NoError(t, err)

		defer func() {
			_ = userService.DeleteUser(ctx, otherUser.ID.String())
		}()

		// Create project with testUser
		createReq := &models.ProjectRequest{
			Name:   fmt.Sprintf("ownership-project-%d", time.Now().UnixNano()),
			Color:  "#FF5733",
			Status: models.ProjectActive,
		}

		project, err := projectService.CreateProject(ctx, testUser.ID.String(), createReq)
		require.NoError(t, err)

		defer func() {
			_ = projectService.DeleteProject(ctx, testUser.ID.String(), project.ID.String())
		}()

		// Try to access project with otherUser (should fail)
		_, err = projectService.GetProject(ctx, otherUser.ID.String(), project.ID.String())
		assert.Error(t, err)

		// Try to update project with otherUser (should fail)
		updateReq := &models.ProjectRequest{
			Name:   "Updated by other user",
			Color:  "#00FF00",
			Status: models.ProjectActive,
		}
		_, err = projectService.UpdateProject(ctx, otherUser.ID.String(), project.ID.String(), updateReq)
		assert.Error(t, err)

		// Try to delete project with otherUser (should fail)
		err = projectService.DeleteProject(ctx, otherUser.ID.String(), project.ID.String())
		assert.Error(t, err)
	})

	t.Run("ConcurrentProjectOperations", func(t *testing.T) {
		// Test concurrent project creation
		const numGoroutines = 5
		results := make(chan *models.Project, numGoroutines)
		errors := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(index int) {
				projectReq := &models.ProjectRequest{
					Name:   fmt.Sprintf("concurrent-project-%d-%d", index, time.Now().UnixNano()),
					Color:  "#FF0000",
					Status: models.ProjectActive,
				}

				project, err := projectService.CreateProject(ctx, testUser.ID.String(), projectReq)
				if err != nil {
					errors <- err
					return
				}
				results <- project
			}(i)
		}

		// Collect results
		var createdProjects []*models.Project
		for i := 0; i < numGoroutines; i++ {
			select {
			case project := <-results:
				createdProjects = append(createdProjects, project)
			case err := <-errors:
				t.Errorf("Concurrent operation failed: %v", err)
			case <-time.After(10 * time.Second):
				t.Error("Timeout waiting for concurrent operations")
			}
		}

		// Clean up created projects
		defer func() {
			for _, project := range createdProjects {
				_ = projectService.DeleteProject(ctx, testUser.ID.String(), project.ID.String())
			}
		}()

		assert.Equal(t, numGoroutines, len(createdProjects))

		// Verify all projects are unique
		projectNames := make(map[string]bool)
		for _, project := range createdProjects {
			assert.False(t, projectNames[project.Name], "Project name should be unique: %s", project.Name)
			projectNames[project.Name] = true
		}
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		// Test getting non-existent project
		_, err := projectService.GetProject(ctx, testUser.ID.String(), uuid.New().String())
		assert.Error(t, err)

		// Test updating non-existent project
		updateReq := &models.ProjectRequest{
			Name:   "Updated name",
			Color:  "#00FF00",
			Status: models.ProjectActive,
		}
		_, err = projectService.UpdateProject(ctx, testUser.ID.String(), uuid.New().String(), updateReq)
		assert.Error(t, err)

		// Test deleting non-existent project
		err = projectService.DeleteProject(ctx, testUser.ID.String(), uuid.New().String())
		assert.Error(t, err)

		// Test invalid UUID formats
		invalidUUIDs := []string{"", "invalid", "not-a-uuid", "12345"}
		for _, invalidUUID := range invalidUUIDs {
			_, err = projectService.GetProject(ctx, testUser.ID.String(), invalidUUID)
			assert.Error(t, err)
		}

		// Test with invalid userID
		validProjectReq := &models.ProjectRequest{
			Name:   "Valid Project",
			Color:  "#FF5733",
			Status: models.ProjectActive,
		}

		for _, invalidUserID := range invalidUUIDs {
			_, err = projectService.CreateProject(ctx, invalidUserID, validProjectReq)
			assert.Error(t, err)
		}
	})

	t.Run("EdgeCases", func(t *testing.T) {
		// Test with minimal project data
		minimalProject := &models.ProjectRequest{
			Name:   fmt.Sprintf("minimal-%d", time.Now().UnixNano()),
			Color:  "#000000",
			Status: models.ProjectActive,
		}

		project, err := projectService.CreateProject(ctx, testUser.ID.String(), minimalProject)
		require.NoError(t, err)
		assert.Nil(t, project.Description)
		assert.False(t, project.IsDefault)

		defer func() {
			_ = projectService.DeleteProject(ctx, testUser.ID.String(), project.ID.String())
		}()

		// Test with maximum length fields
		maxProject := &models.ProjectRequest{
			Name:        fmt.Sprintf("max-length-project-name-with-many-chars-%d", time.Now().UnixNano()),
			Description: stringPtr("This is a very long description that tests the maximum length handling for project descriptions in the database and service layer. It should be within the allowed limit."),
			Color:       "#FFFFFF",
			Status:      models.ProjectCompleted,
			IsDefault:   true,
		}

		project2, err := projectService.CreateProject(ctx, testUser.ID.String(), maxProject)
		require.NoError(t, err)

		defer func() {
			_ = projectService.DeleteProject(ctx, testUser.ID.String(), project2.ID.String())
		}()

		// Test getting projects when user has no projects (different user)
		emptyUser, err := userService.CreateUser(ctx, &models.UserRegistration{
			Email:     fmt.Sprintf("empty-user-%d@example.com", time.Now().UnixNano()),
			Password:  "password123",
			FirstName: "Empty",
			LastName:  "User",
			Timezone:  "UTC",
		})
		require.NoError(t, err)

		defer func() {
			_ = userService.DeleteUser(ctx, emptyUser.ID.String())
		}()

		emptyProjects, err := projectService.GetUserProjects(ctx, emptyUser.ID.String())
		require.NoError(t, err)
		assert.Empty(t, emptyProjects)

		emptyActiveProjects, err := projectService.GetActiveProjects(ctx, emptyUser.ID.String())
		require.NoError(t, err)
		assert.Empty(t, emptyActiveProjects)

		_, err = projectService.GetDefaultProject(ctx, emptyUser.ID.String())
		assert.Error(t, err) // Should error when no default project exists
	})
}

// BenchmarkProjectService benchmarks project service operations
func BenchmarkProjectService(b *testing.B) {
	db := testutils.DB(b)

	// Create test logger
	testLogger := logging.NewTestLogger()

	projectService := services.NewProjectService(db, testLogger)
	userService := services.NewUserService(db, testLogger)
	ctx := context.Background()

	// Create test user
	testUser, err := userService.CreateUser(ctx, &models.UserRegistration{
		Email:     "bench-project-user@example.com",
		Password:  "password123",
		FirstName: "Bench",
		LastName:  "User",
		Timezone:  "UTC",
	})
	if err != nil {
		b.Fatalf("Failed to create test user: %v", err)
	}

	defer func() {
		_ = userService.DeleteUser(ctx, testUser.ID.String())
	}()

	b.Run("CreateProject", func(b *testing.B) {
		var createdIDs []string

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req := &models.ProjectRequest{
				Name:   fmt.Sprintf("bench-project-%d", i),
				Color:  "#FF0000",
				Status: models.ProjectActive,
			}

			project, err := projectService.CreateProject(ctx, testUser.ID.String(), req)
			if err != nil {
				b.Fatalf("Failed to create project: %v", err)
			}
			createdIDs = append(createdIDs, project.ID.String())
		}
		b.StopTimer()

		// Cleanup
		for _, id := range createdIDs {
			_ = projectService.DeleteProject(ctx, testUser.ID.String(), id)
		}
	})

	b.Run("GetProject", func(b *testing.B) {
		// Setup: create a project to fetch
		req := &models.ProjectRequest{
			Name:   "bench-get-project",
			Color:  "#FF0000",
			Status: models.ProjectActive,
		}

		project, err := projectService.CreateProject(ctx, testUser.ID.String(), req)
		if err != nil {
			b.Fatalf("Failed to create test project: %v", err)
		}

		defer func() {
			_ = projectService.DeleteProject(ctx, testUser.ID.String(), project.ID.String())
		}()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := projectService.GetProject(ctx, testUser.ID.String(), project.ID.String())
			if err != nil {
				b.Fatalf("Failed to get project: %v", err)
			}
		}
	})

	b.Run("GetUserProjects", func(b *testing.B) {
		// Setup: create some projects
		for i := 0; i < 10; i++ {
			req := &models.ProjectRequest{
				Name:   fmt.Sprintf("bench-list-project-%d", i),
				Color:  "#FF0000",
				Status: models.ProjectActive,
			}
			_, _ = projectService.CreateProject(ctx, testUser.ID.String(), req)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := projectService.GetUserProjects(ctx, testUser.ID.String())
			if err != nil {
				b.Fatalf("Failed to get user projects: %v", err)
			}
		}
	})
}
