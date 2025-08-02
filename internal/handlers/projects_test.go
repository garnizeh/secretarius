package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/garnizeh/englog/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectHandler_CreateProject(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful project creation", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create and login user
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Test project creation
		projectReq := models.ProjectRequest{
			Name:        "Test Project",
			Description: stringPtr("A test project"),
			Color:       "#FF5733",
			Status:      models.ProjectActive,
			IsDefault:   false,
		}

		body, _ := json.Marshal(projectReq)
		req := httptest.NewRequest(http.MethodPost, "/v1/projects", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response responseData[models.Project]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Test Project", response.Data.Name)
		assert.Equal(t, models.ProjectActive, response.Data.Status)
	})

	t.Run("unauthorized access", func(t *testing.T) {
		// Setup test environment
		router, _, _, _, _ := RouterWithServices(t)

		projectReq := models.ProjectRequest{
			Name:   "Test Project",
			Status: models.ProjectActive,
		}

		body, _ := json.Marshal(projectReq)
		req := httptest.NewRequest(http.MethodPost, "/v1/projects", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create and login user
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		req := httptest.NewRequest(http.MethodPost, "/v1/projects", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestProjectHandler_GetProject(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful project retrieval", func(t *testing.T) {
		// Setup test environment
		router, userService, projectService, _, _ := RouterWithServices(t)

		// Create and login user
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Create a project first
		project := createTestProject(t, projectService, user.ID.String())

		req := httptest.NewRequest(http.MethodGet, "/v1/projects/"+project.ID.String(), nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response responseData[models.Project]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, project.ID, response.Data.ID)
		assert.Equal(t, project.Name, response.Data.Name)
	})

	t.Run("project not found", func(t *testing.T) {
		// Setup test environment
		router, userService, _, _, _ := RouterWithServices(t)

		// Create and login user
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		req := httptest.NewRequest(http.MethodGet, "/v1/projects/00000000-0000-0000-0000-000000000000", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestProjectHandler_GetProjects(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful projects list", func(t *testing.T) {
		// Setup test environment
		router, userService, projectService, _, _ := RouterWithServices(t)

		// Create and login user
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Create some projects
		createTestProject(t, projectService, user.ID.String())
		createTestProject(t, projectService, user.ID.String())

		req := httptest.NewRequest(http.MethodGet, "/v1/projects", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response responseData[[]*models.Project]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(response.Data), 2)
	})
}

func TestProjectHandler_UpdateProject(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful project update", func(t *testing.T) {
		// Setup test environment
		router, userService, projectService, _, _ := RouterWithServices(t)

		// Create and login user
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Create a project first
		project := createTestProject(t, projectService, user.ID.String())

		updateReq := models.ProjectRequest{
			Name:        "Updated Project",
			Description: stringPtr("Updated description"),
			Color:       "#00FF00",
			Status:      models.ProjectCompleted,
			IsDefault:   false,
		}

		body, _ := json.Marshal(updateReq)
		req := httptest.NewRequest(http.MethodPut, "/v1/projects/"+project.ID.String(), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response responseData[models.Project]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Updated Project", response.Data.Name)
		assert.Equal(t, models.ProjectCompleted, response.Data.Status)
	})
}

func TestProjectHandler_DeleteProject(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful project deletion", func(t *testing.T) {
		// Setup test environment
		router, userService, projectService, _, _ := RouterWithServices(t)

		// Create and login user
		user := createTestUser(t, userService)
		token := loginUser(t, router, user.Email, "password123")

		// Create a project first
		project := createTestProject(t, projectService, user.ID.String())

		req := httptest.NewRequest(http.MethodDelete, "/v1/projects/"+project.ID.String(), nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
