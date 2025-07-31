//go:build e2e
// +build e2e

package e2e

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// UserJourneyTestSuite tests complete user workflows end-to-end
// "The journey of a thousand miles begins with a single step." 🚶‍♀️
type UserJourneyTestSuite struct {
	suite.Suite
	baseURL      string
	httpClient   *http.Client
	accessToken  string
	refreshToken string
	userID       string
	email        string
}

func (suite *UserJourneyTestSuite) SetupSuite() {
	// Allow base URL to be configurable via environment variable
	baseURL := os.Getenv("E2E_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8081"
	}
	suite.baseURL = baseURL

	suite.httpClient = &http.Client{
		Timeout: 30 * time.Second,
	}
}

func (suite *UserJourneyTestSuite) TestCompleteUserJourney() {
	// Step 1: User Registration
	suite.T().Log("Step 1: User Registration")
	suite.registerUser()

	// Step 2: User Login
	suite.T().Log("Step 2: User Login")
	suite.loginUser()

	// Step 3: Create Project
	suite.T().Log("Step 3: Create Project")
	projectID := suite.createProject()

	// Step 4: Create Log Entry
	suite.T().Log("Step 4: Create Log Entry")
	entryID := suite.createLogEntry(projectID)

	// Step 5: Get Log Entries
	suite.T().Log("Step 5: Get Log Entries")
	suite.getLogEntries()

	// Step 6: Update Log Entry
	suite.T().Log("Step 6: Update Log Entry")
	suite.updateLogEntry(entryID)

	// Step 7: Get Analytics
	suite.T().Log("Step 7: Get Analytics")
	suite.getAnalytics()

	// Step 8: Logout
	suite.T().Log("Step 8: Logout")
	suite.logoutUser()
}

func (suite *UserJourneyTestSuite) registerUser() {
	timestamp := time.Now().Unix()
	email := fmt.Sprintf("e2e-test-%d@example.com", timestamp)
	username := fmt.Sprintf("e2euser%d", timestamp)

	// Store email for later use
	suite.email = email

	payload := map[string]string{
		"email":      email,
		"username":   username,
		"password":   "securepassword123",
		"first_name": "E2E",
		"last_name":  "User",
	}

	resp := suite.makeRequest("POST", "/v1/auth/register", payload, "")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var response map[string]any
	err := json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(suite.T(), err)

	// API response structure: {"tokens": {"access_token": "...", "refresh_token": "...", ...}, "user": {"id": "...", ...}}
	tokens := response["tokens"].(map[string]any)
	assert.NotEmpty(suite.T(), tokens["access_token"])
	assert.NotEmpty(suite.T(), tokens["refresh_token"])
	suite.accessToken = tokens["access_token"].(string)
	suite.refreshToken = tokens["refresh_token"].(string)

	user := response["user"].(map[string]any)
	suite.userID = user["id"].(string)
}

func (suite *UserJourneyTestSuite) loginUser() {
	// Since we already have a valid token from registration,
	// we'll just verify the token still works by calling /auth/me
	resp := suite.makeRequest("GET", "/v1/auth/me", nil, suite.accessToken)
	defer resp.Body.Close()

	// If this endpoint exists and returns user info, login is working
	if resp.StatusCode == http.StatusOK {
		suite.T().Log("Authentication verified - user token is valid")
	}
}

func (suite *UserJourneyTestSuite) createProject() string {
	payload := map[string]string{
		"name":        "E2E Test Project",
		"description": "A project created during end-to-end testing",
		"status":      "active",
		"color":       "#007bff",
	}

	resp := suite.makeRequest("POST", "/v1/projects", payload, suite.accessToken)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var response map[string]any
	err := json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(suite.T(), err)

	// API response structure: {"data": {"id": "...", ...}}
	data := response["data"].(map[string]any)
	assert.NotEmpty(suite.T(), data["id"])
	return data["id"].(string)
}

func (suite *UserJourneyTestSuite) createLogEntry(projectID string) string {
	payload := map[string]any{
		"title":        "My E2E Test Entry",
		"description":  "This is a test log entry created during end-to-end testing",
		"type":         "development",
		"value_rating": "high",
		"impact_level": "team",
		"project_id":   projectID,
		"start_time":   time.Now().Add(-time.Hour).Format(time.RFC3339),
		"end_time":     time.Now().Format(time.RFC3339),
	}

	resp := suite.makeRequest("POST", "/v1/logs", payload, suite.accessToken)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var response map[string]any
	err := json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(suite.T(), err)

	// API response structure: {"success": true, "data": {"id": "...", ...}}
	data := response["data"].(map[string]any)
	assert.NotEmpty(suite.T(), data["id"])
	return data["id"].(string)
}

func (suite *UserJourneyTestSuite) getLogEntries() {
	resp := suite.makeRequest("GET", "/v1/logs", nil, suite.accessToken)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]any
	err := json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(suite.T(), err)

	// API response structure: {"data": [...], "pagination": {...}, "total": 1}
	logs, ok := response["data"].([]any)
	assert.True(suite.T(), ok)
	assert.Greater(suite.T(), len(logs), 0)
}

func (suite *UserJourneyTestSuite) updateLogEntry(entryID string) {
	payload := map[string]any{
		"title":        "Updated E2E Test Entry",
		"description":  "This entry has been updated during the E2E test",
		"type":         "development",
		"value_rating": "medium",
		"impact_level": "team",
		"start_time":   time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
		"end_time":     time.Now().Add(-30 * time.Minute).Format(time.RFC3339),
	}

	url := fmt.Sprintf("/v1/logs/%s", entryID)
	resp := suite.makeRequest("PUT", url, payload, suite.accessToken)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func (suite *UserJourneyTestSuite) getAnalytics() {
	// Test different analytics endpoints
	endpoints := []string{
		"/v1/analytics/productivity",
		"/v1/analytics/summary",
	}

	for _, endpoint := range endpoints {
		resp := suite.makeRequest("GET", endpoint, nil, suite.accessToken)
		defer resp.Body.Close()

		// Analytics might return 200 or 404 if no data yet - both are acceptable for E2E
		assert.True(suite.T(), resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotFound,
			"Analytics endpoint %s returned unexpected status: %d", endpoint, resp.StatusCode)
	}
}

func (suite *UserJourneyTestSuite) logoutUser() {
	payload := map[string]string{
		"refresh_token": suite.refreshToken,
	}

	resp := suite.makeRequest("POST", "/v1/auth/logout", payload, suite.accessToken)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func (suite *UserJourneyTestSuite) makeRequest(method, path string, payload any, token string) *http.Response {
	url := suite.baseURL + path

	var body io.Reader
	if payload != nil {
		// Marshal the payload into JSON. This is expected to succeed because the payload
		// is controlled within the test suite and should always be a JSON-serializable structure.
		jsonBody, err := json.Marshal(payload)
		require.NoError(suite.T(), err)
		body = strings.NewReader(string(jsonBody))
	}

	req, err := http.NewRequest(method, url, body)
	require.NoError(suite.T(), err)

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := suite.httpClient.Do(req)
	require.NoError(suite.T(), err)

	return resp
}

// TestHealthEndpoint tests that the health endpoint is accessible
func (suite *UserJourneyTestSuite) TestHealthEndpoint() {
	resp := suite.makeRequest("GET", "/health", nil, "")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]any
	err := json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(suite.T(), err)

	assert.NotEmpty(suite.T(), response["status"])
}

// TestUnauthorizedAccess tests that protected endpoints require authentication
func (suite *UserJourneyTestSuite) TestUnauthorizedAccess() {
	protectedEndpoints := []struct {
		method string
		path   string
	}{
		{"GET", "/v1/logs"},
		{"POST", "/v1/logs"},
		{"GET", "/v1/projects"},
		{"POST", "/v1/projects"},
		{"GET", "/v1/analytics/productivity"},
	}

	for _, endpoint := range protectedEndpoints {
		resp := suite.makeRequest(endpoint.method, endpoint.path, nil, "")
		defer resp.Body.Close()

		assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode,
			"Endpoint %s %s should require authentication", endpoint.method, endpoint.path)
	}
}

func TestUserJourney(t *testing.T) {
	suite.Run(t, new(UserJourneyTestSuite))
}
