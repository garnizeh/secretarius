//go:build integration
// +build integration

package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/garnizeh/englog/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectHandler_Integration(t *testing.T) {
	t.Run("end_to_end_project_lifecycle", testEndToEndProjectLifecycle)
	t.Run("multiple_users_project_isolation", testMultipleUsersProjectIsolation)
	t.Run("concurrent_project_operations", testConcurrentProjectOperations)
	t.Run("default_project_integration", testDefaultProjectIntegration)
	t.Run("validation_and_error_handling_integration", testValidationAndErrorHandlingIntegration)
}

func testEndToEndProjectLifecycle(t *testing.T) {
	// Setup router and services
	router, userService, _, _, _ := setupTestRouterWithServices(t)

	// Create and login user
	user := createTestUser(t, userService)
	token := loginUser(t, router, user.Email, "password123")

	// 1. Create a project
	createData := map[string]any{
		"name":        "Integration Test Project",
		"description": "A project for end-to-end testing",
		"color":       "#FF5733",
		"status":      "active",
		"is_default":  false,
	}
	createBody, _ := json.Marshal(createData)
	createRequest, _ := http.NewRequest("POST", "/v1/projects", bytes.NewBuffer(createBody))
	createRequest.Header.Set("Content-Type", "application/json")
	createRequest.Header.Set("Authorization", "Bearer "+token)
	createResponse := httptest.NewRecorder()
	router.ServeHTTP(createResponse, createRequest)

	assert.Equal(t, http.StatusCreated, createResponse.Code)

	var createResult responseData[models.Project]
	err := json.Unmarshal(createResponse.Body.Bytes(), &createResult)
	require.NoError(t, err)
	projectID := createResult.Data.ID.String()

	// 2. Update the project
	updateData := map[string]any{
		"name":        "Updated Integration Project",
		"description": "Updated description for testing",
		"color":       "#00FF00",
		"status":      "on_hold",
		"is_default":  false,
	}
	updateBody, _ := json.Marshal(updateData)
	updateRequest, _ := http.NewRequest("PUT", "/v1/projects/"+projectID, bytes.NewBuffer(updateBody))
	updateRequest.Header.Set("Content-Type", "application/json")
	updateRequest.Header.Set("Authorization", "Bearer "+token)
	updateResponse := httptest.NewRecorder()
	router.ServeHTTP(updateResponse, updateRequest)

	assert.Equal(t, http.StatusOK, updateResponse.Code)

	// 3. Get the updated project
	getRequest, _ := http.NewRequest("GET", "/v1/projects/"+projectID, nil)
	getRequest.Header.Set("Authorization", "Bearer "+token)
	getResponse := httptest.NewRecorder()
	router.ServeHTTP(getResponse, getRequest)

	assert.Equal(t, http.StatusOK, getResponse.Code)

	var getResult responseData[models.Project]
	err = json.Unmarshal(getResponse.Body.Bytes(), &getResult)
	require.NoError(t, err)
	assert.Equal(t, "Updated Integration Project", getResult.Data.Name)
	assert.Equal(t, models.ProjectOnHold, getResult.Data.Status)

	// 4. List projects and confirm our project exists
	listRequest, _ := http.NewRequest("GET", "/v1/projects", nil)
	listRequest.Header.Set("Authorization", "Bearer "+token)
	listResponse := httptest.NewRecorder()
	router.ServeHTTP(listResponse, listRequest)

	assert.Equal(t, http.StatusOK, listResponse.Code)

	var listResult responseData[[]models.Project]
	err = json.Unmarshal(listResponse.Body.Bytes(), &listResult)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(listResult.Data), 1)

	// Find our project in the list
	found := false
	for _, project := range listResult.Data {
		if project.ID.String() == projectID {
			found = true
			assert.Equal(t, "Updated Integration Project", project.Name)
			break
		}
	}
	assert.True(t, found, "Project should be found in list")

	// 5. Delete the project
	deleteRequest, _ := http.NewRequest("DELETE", "/v1/projects/"+projectID, nil)
	deleteRequest.Header.Set("Authorization", "Bearer "+token)
	deleteResponse := httptest.NewRecorder()
	router.ServeHTTP(deleteResponse, deleteRequest)

	assert.Equal(t, http.StatusOK, deleteResponse.Code, "Delete should return 200")

	// 6. Verify project is deleted
	verifyRequest, _ := http.NewRequest("GET", "/v1/projects/"+projectID, nil)
	verifyRequest.Header.Set("Authorization", "Bearer "+token)
	verifyResponse := httptest.NewRecorder()
	router.ServeHTTP(verifyResponse, verifyRequest)

	assert.Equal(t, http.StatusNotFound, verifyResponse.Code)
}

func testMultipleUsersProjectIsolation(t *testing.T) {
	// Setup router and services
	router, userService, _, _, _ := setupTestRouterWithServices(t)

	// Create and login first user
	user1 := createTestUser(t, userService)
	token1 := loginUser(t, router, user1.Email, "password123")

	// Create and login second user
	user2 := createTestUser(t, userService)
	token2 := loginUser(t, router, user2.Email, "password123")

	// User 1 creates a project
	user1ProjectData := map[string]any{
		"name":        "User 1 Project",
		"description": "Project belonging to user 1",
		"color":       "#FF5733",
		"status":      "active",
		"is_default":  false,
	}
	user1Body, _ := json.Marshal(user1ProjectData)
	user1Request, _ := http.NewRequest("POST", "/v1/projects", bytes.NewBuffer(user1Body))
	user1Request.Header.Set("Content-Type", "application/json")
	user1Request.Header.Set("Authorization", "Bearer "+token1)
	user1Response := httptest.NewRecorder()
	router.ServeHTTP(user1Response, user1Request)

	assert.Equal(t, http.StatusCreated, user1Response.Code)

	var user1Result responseData[models.Project]
	err := json.Unmarshal(user1Response.Body.Bytes(), &user1Result)
	require.NoError(t, err)
	user1ProjectID := user1Result.Data.ID.String()

	// User 2 creates a project
	user2ProjectData := map[string]any{
		"name":        "User 2 Project",
		"description": "Project belonging to user 2",
		"color":       "#00FF00",
		"status":      "active",
		"is_default":  false,
	}
	user2Body, _ := json.Marshal(user2ProjectData)
	user2Request, _ := http.NewRequest("POST", "/v1/projects", bytes.NewBuffer(user2Body))
	user2Request.Header.Set("Content-Type", "application/json")
	user2Request.Header.Set("Authorization", "Bearer "+token2)
	user2Response := httptest.NewRecorder()
	router.ServeHTTP(user2Response, user2Request)

	assert.Equal(t, http.StatusCreated, user2Response.Code)

	var user2Result responseData[models.Project]
	err = json.Unmarshal(user2Response.Body.Bytes(), &user2Result)
	require.NoError(t, err)
	user2ProjectID := user2Result.Data.ID.String()

	// User 1 lists projects (should only see their own)
	user1ListRequest, _ := http.NewRequest("GET", "/v1/projects", nil)
	user1ListRequest.Header.Set("Authorization", "Bearer "+token1)
	user1ListResponse := httptest.NewRecorder()
	router.ServeHTTP(user1ListResponse, user1ListRequest)

	assert.Equal(t, http.StatusOK, user1ListResponse.Code)

	var user1ListResult responseData[[]models.Project]
	err = json.Unmarshal(user1ListResponse.Body.Bytes(), &user1ListResult)
	require.NoError(t, err)

	// User 2 lists projects (should only see their own)
	user2ListRequest, _ := http.NewRequest("GET", "/v1/projects", nil)
	user2ListRequest.Header.Set("Authorization", "Bearer "+token2)
	user2ListResponse := httptest.NewRecorder()
	router.ServeHTTP(user2ListResponse, user2ListRequest)

	assert.Equal(t, http.StatusOK, user2ListResponse.Code)

	var user2ListResult responseData[[]models.Project]
	err = json.Unmarshal(user2ListResponse.Body.Bytes(), &user2ListResult)
	require.NoError(t, err)

	// User 1 should not be able to access User 2's project
	user1AccessUser2Request, _ := http.NewRequest("GET", "/v1/projects/"+user2ProjectID, nil)
	user1AccessUser2Request.Header.Set("Authorization", "Bearer "+token1)
	user1AccessUser2Response := httptest.NewRecorder()
	router.ServeHTTP(user1AccessUser2Response, user1AccessUser2Request)

	assert.Equal(t, http.StatusNotFound, user1AccessUser2Response.Code)

	// User 2 should not be able to access User 1's project
	user2AccessUser1Request, _ := http.NewRequest("GET", "/v1/projects/"+user1ProjectID, nil)
	user2AccessUser1Request.Header.Set("Authorization", "Bearer "+token2)
	user2AccessUser1Response := httptest.NewRecorder()
	router.ServeHTTP(user2AccessUser1Response, user2AccessUser1Request)

	assert.Equal(t, http.StatusNotFound, user2AccessUser1Response.Code)
}

func testConcurrentProjectOperations(t *testing.T) {
	// Setup router and services
	router, userService, _, _, _ := setupTestRouterWithServices(t)

	// Create and login user
	user := createTestUser(t, userService)
	token := loginUser(t, router, user.Email, "password123")

	// Create a project first
	projectData := map[string]any{
		"name":        "Concurrent Test Project",
		"description": "Project for concurrent operations testing",
		"color":       "#FF5733",
		"status":      "active",
		"is_default":  false,
	}
	body, _ := json.Marshal(projectData)
	request, _ := http.NewRequest("POST", "/v1/projects", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusCreated, response.Code)

	var result responseData[models.Project]
	err := json.Unmarshal(response.Body.Bytes(), &result)
	require.NoError(t, err)
	projectID := result.Data.ID.String()

	// Test concurrent reads
	var wg sync.WaitGroup
	const numConcurrentReads = 5

	for i := 0; i < numConcurrentReads; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			getRequest, _ := http.NewRequest("GET", "/v1/projects/"+projectID, nil)
			getRequest.Header.Set("Authorization", "Bearer "+token)
			getResponse := httptest.NewRecorder()
			router.ServeHTTP(getResponse, getRequest)

			assert.Equal(t, http.StatusOK, getResponse.Code, "Concurrent read %d should succeed", index)
		}(i)
	}

	wg.Wait()

	// Test sequential updates to avoid concurrency issues
	updateData1 := map[string]any{
		"name":        "Updated Name 1",
		"description": "First update",
		"color":       "#00FF00",
		"status":      "on_hold",
		"is_default":  false,
	}
	updateBody1, _ := json.Marshal(updateData1)
	updateRequest1, _ := http.NewRequest("PUT", "/v1/projects/"+projectID, bytes.NewBuffer(updateBody1))
	updateRequest1.Header.Set("Content-Type", "application/json")
	updateRequest1.Header.Set("Authorization", "Bearer "+token)
	updateResponse1 := httptest.NewRecorder()
	router.ServeHTTP(updateResponse1, updateRequest1)

	assert.Equal(t, http.StatusOK, updateResponse1.Code, "First update should succeed")

	// Second update
	updateData2 := map[string]any{
		"name":        "Updated Name 2",
		"description": "Second update",
		"color":       "#0000FF",
		"status":      "completed",
		"is_default":  false,
	}
	updateBody2, _ := json.Marshal(updateData2)
	updateRequest2, _ := http.NewRequest("PUT", "/v1/projects/"+projectID, bytes.NewBuffer(updateBody2))
	updateRequest2.Header.Set("Content-Type", "application/json")
	updateRequest2.Header.Set("Authorization", "Bearer "+token)
	updateResponse2 := httptest.NewRecorder()
	router.ServeHTTP(updateResponse2, updateRequest2)

	assert.Equal(t, http.StatusOK, updateResponse2.Code, "Second update should succeed")

	// Verify final state
	finalRequest, _ := http.NewRequest("GET", "/v1/projects/"+projectID, nil)
	finalRequest.Header.Set("Authorization", "Bearer "+token)
	finalResponse := httptest.NewRecorder()
	router.ServeHTTP(finalResponse, finalRequest)

	assert.Equal(t, http.StatusOK, finalResponse.Code)

	var finalResult responseData[models.Project]
	err = json.Unmarshal(finalResponse.Body.Bytes(), &finalResult)
	require.NoError(t, err)

	// One of the updates should have succeeded
	assert.Equal(t, models.ProjectCompleted, finalResult.Data.Status)
	assert.Equal(t, "Updated Name 2", finalResult.Data.Name, "Final name should be from the last update")
}

func testDefaultProjectIntegration(t *testing.T) {
	// Setup router and services
	router, userService, _, _, _ := setupTestRouterWithServices(t)

	// Create and login user
	user := createTestUser(t, userService)
	token := loginUser(t, router, user.Email, "password123")

	// Create a default project
	defaultProjectData := map[string]any{
		"name":        "Default Project",
		"description": "This is the default project",
		"color":       "#FF5733",
		"status":      "active",
		"is_default":  true,
	}
	defaultBody, _ := json.Marshal(defaultProjectData)
	defaultRequest, _ := http.NewRequest("POST", "/v1/projects", bytes.NewBuffer(defaultBody))
	defaultRequest.Header.Set("Content-Type", "application/json")
	defaultRequest.Header.Set("Authorization", "Bearer "+token)
	defaultResponse := httptest.NewRecorder()
	router.ServeHTTP(defaultResponse, defaultRequest)

	assert.Equal(t, http.StatusCreated, defaultResponse.Code)

	var defaultResult responseData[models.Project]
	err := json.Unmarshal(defaultResponse.Body.Bytes(), &defaultResult)
	require.NoError(t, err)
	defaultProjectID := defaultResult.Data.ID.String()

	// Create a second project (not default)
	secondProjectData := map[string]any{
		"name":        "Second Project",
		"description": "This will become the new default",
		"color":       "#00FF00",
		"status":      "active",
		"is_default":  false,
	}
	secondBody, _ := json.Marshal(secondProjectData)
	secondRequest, _ := http.NewRequest("POST", "/v1/projects", bytes.NewBuffer(secondBody))
	secondRequest.Header.Set("Content-Type", "application/json")
	secondRequest.Header.Set("Authorization", "Bearer "+token)
	secondResponse := httptest.NewRecorder()
	router.ServeHTTP(secondResponse, secondRequest)

	assert.Equal(t, http.StatusCreated, secondResponse.Code)

	var secondResult responseData[models.Project]
	err = json.Unmarshal(secondResponse.Body.Bytes(), &secondResult)
	require.NoError(t, err)
	secondProjectID := secondResult.Data.ID.String()

	// Create a regular project
	regularProjectData := map[string]any{
		"name":        "Regular Project",
		"description": "This is not a default project",
		"color":       "#0000FF",
		"status":      "active",
		"is_default":  false,
	}
	regularBody, _ := json.Marshal(regularProjectData)
	regularRequest, _ := http.NewRequest("POST", "/v1/projects", bytes.NewBuffer(regularBody))
	regularRequest.Header.Set("Content-Type", "application/json")
	regularRequest.Header.Set("Authorization", "Bearer "+token)
	regularResponse := httptest.NewRecorder()
	router.ServeHTTP(regularResponse, regularRequest)

	assert.Equal(t, http.StatusCreated, regularResponse.Code)

	// Update the first project to not be default anymore
	updateData := map[string]any{
		"name":        "Updated Default Project",
		"description": "No longer default",
		"color":       "#FF5733",
		"status":      "active",
		"is_default":  false,
	}
	updateBody, _ := json.Marshal(updateData)
	updateRequest, _ := http.NewRequest("PUT", "/v1/projects/"+defaultProjectID, bytes.NewBuffer(updateBody))
	updateRequest.Header.Set("Content-Type", "application/json")
	updateRequest.Header.Set("Authorization", "Bearer "+token)
	updateResponse := httptest.NewRecorder()
	router.ServeHTTP(updateResponse, updateRequest)

	assert.Equal(t, http.StatusOK, updateResponse.Code)

	// Now update the second project to be the default
	makeDefaultData := map[string]any{
		"name":        "New Default Project",
		"description": "This is now the default",
		"color":       "#00FF00",
		"status":      "active",
		"is_default":  true,
	}
	makeDefaultBody, _ := json.Marshal(makeDefaultData)
	makeDefaultRequest, _ := http.NewRequest("PUT", "/v1/projects/"+secondProjectID, bytes.NewBuffer(makeDefaultBody))
	makeDefaultRequest.Header.Set("Content-Type", "application/json")
	makeDefaultRequest.Header.Set("Authorization", "Bearer "+token)
	makeDefaultResponse := httptest.NewRecorder()
	router.ServeHTTP(makeDefaultResponse, makeDefaultRequest)

	assert.Equal(t, http.StatusOK, makeDefaultResponse.Code)

	// List all projects and verify default status
	listRequest, _ := http.NewRequest("GET", "/v1/projects", nil)
	listRequest.Header.Set("Authorization", "Bearer "+token)
	listResponse := httptest.NewRecorder()
	router.ServeHTTP(listResponse, listRequest)

	assert.Equal(t, http.StatusOK, listResponse.Code)

	var listResult responseData[[]models.Project]
	err = json.Unmarshal(listResponse.Body.Bytes(), &listResult)
	require.NoError(t, err)

	// Should have exactly one default project
	defaultCount := 0
	var currentDefault *models.Project
	for _, project := range listResult.Data {
		if project.IsDefault {
			defaultCount++
			currentDefault = &project
		}
	}

	assert.Equal(t, 1, defaultCount, "Should have exactly one default project")
	require.NotNil(t, currentDefault, "Should have exactly one default project")
	assert.Equal(t, secondProjectID, currentDefault.ID.String(), "Second project should be the default")
	assert.Equal(t, "New Default Project", currentDefault.Name, "Default project should have updated name")
}

func testValidationAndErrorHandlingIntegration(t *testing.T) {
	// Setup router and services
	router, userService, _, _, _ := setupTestRouterWithServices(t)

	// Create and login user
	user := createTestUser(t, userService)
	token := loginUser(t, router, user.Email, "password123")

	// Test invalid project creation - missing required fields
	invalidData := map[string]any{
		"description": "Missing name and other required fields",
	}
	invalidBody, _ := json.Marshal(invalidData)
	invalidRequest, _ := http.NewRequest("POST", "/v1/projects", bytes.NewBuffer(invalidBody))
	invalidRequest.Header.Set("Content-Type", "application/json")
	invalidRequest.Header.Set("Authorization", "Bearer "+token)
	invalidResponse := httptest.NewRecorder()
	router.ServeHTTP(invalidResponse, invalidRequest)

	assert.Equal(t, http.StatusBadRequest, invalidResponse.Code)

	// Test invalid project update
	nonExistentID := "b1f021ae-76a5-4c8d-8ab0-51c6a634e1e8"
	updateData := map[string]any{
		"name": "Updated Name",
	}
	updateBody, _ := json.Marshal(updateData)
	updateRequest, _ := http.NewRequest("PUT", "/v1/projects/"+nonExistentID, bytes.NewBuffer(updateBody))
	updateRequest.Header.Set("Content-Type", "application/json")
	updateRequest.Header.Set("Authorization", "Bearer "+token)
	updateResponse := httptest.NewRecorder()
	router.ServeHTTP(updateResponse, updateRequest)

	assert.Equal(t, http.StatusBadRequest, updateResponse.Code)

	// Test deleting non-existent project
	deleteRequest, _ := http.NewRequest("DELETE", "/v1/projects/"+nonExistentID, nil)
	deleteRequest.Header.Set("Authorization", "Bearer "+token)
	deleteResponse := httptest.NewRecorder()
	router.ServeHTTP(deleteResponse, deleteRequest)

	assert.Equal(t, http.StatusInternalServerError, deleteResponse.Code)

	// Test invalid UUID in path
	invalidUUIDRequest, _ := http.NewRequest("GET", "/v1/projects/invalid-uuid", nil)
	invalidUUIDRequest.Header.Set("Authorization", "Bearer "+token)
	invalidUUIDResponse := httptest.NewRecorder()
	router.ServeHTTP(invalidUUIDResponse, invalidUUIDRequest)

	assert.Equal(t, http.StatusBadRequest, invalidUUIDResponse.Code)

	// Test unauthorized access
	unauthorizedRequest, _ := http.NewRequest("GET", "/v1/projects", nil)
	unauthorizedResponse := httptest.NewRecorder()
	router.ServeHTTP(unauthorizedResponse, unauthorizedRequest)

	assert.Equal(t, http.StatusUnauthorized, unauthorizedResponse.Code)

	// Test with invalid token
	invalidTokenRequest, _ := http.NewRequest("GET", "/v1/projects", nil)
	invalidTokenRequest.Header.Set("Authorization", "Bearer invalid-token")
	invalidTokenResponse := httptest.NewRecorder()
	router.ServeHTTP(invalidTokenResponse, invalidTokenRequest)

	assert.Equal(t, http.StatusUnauthorized, invalidTokenResponse.Code)
}
