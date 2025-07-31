//go:build integration
// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/garnizeh/englog/internal/auth"
	"github.com/garnizeh/englog/internal/config"
	"github.com/garnizeh/englog/internal/handlers"
	"github.com/garnizeh/englog/internal/logging"
	"github.com/garnizeh/englog/internal/models"
	"github.com/garnizeh/englog/internal/services"
	"github.com/garnizeh/englog/internal/store/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// APITestSuite tests HTTP API endpoints with real database
// "Integration testing is where the rubber meets the road." 🛣️
type APITestSuite struct {
	suite.Suite
	router      *gin.Engine
	authService *auth.AuthService
	userService *services.UserService
	logger      *logging.Logger
}

func (suite *APITestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	// Initialize logger
	logConfig := config.LoggingConfig{
		Level:     slog.LevelDebug,
		Format:    "text",
		AddSource: true,
	}
	suite.logger = logging.NewLogger(logConfig)

	// Use the centralized test database utility
	db := testutils.DB(suite.T())

	// Initialize services
	suite.authService = auth.NewAuthServiceForTest(db, suite.logger, "test-secret-key")
	suite.userService = services.NewUserService(db, suite.logger)

	// Setup router
	suite.setupRouter()
}

func (suite *APITestSuite) setupRouter() {
	suite.router = gin.New()

	// Add middleware
	suite.router.Use(gin.Recovery())

	// Setup routes
	v1 := suite.router.Group("/v1")

	// Health routes
	v1.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Auth routes - using AuthService methods directly
	authGroup := v1.Group("/auth")
	{
		authGroup.POST("/register", suite.authService.RegisterHandler)
		authGroup.POST("/login", suite.authService.LoginHandler)
		authGroup.POST("/refresh", suite.authService.RefreshHandler)
		authGroup.POST("/logout", suite.authService.LogoutHandler)
	}

	// User routes with auth middleware
	userHandler := handlers.NewUserHandler(suite.userService)
	userGroup := v1.Group("/users")
	userGroup.Use(suite.authService.RequireAuth())
	{
		userGroup.GET("/profile", userHandler.GetProfile)
		userGroup.PUT("/profile", userHandler.UpdateProfile)
		userGroup.POST("/change-password", userHandler.ChangePassword)
		userGroup.DELETE("/account", userHandler.DeleteAccount)
	}
}

func (suite *APITestSuite) TearDownSuite() {
	// Cleanup is handled automatically by testutils.DB()
}

func (suite *APITestSuite) SetupTest() {
	// Database cleanup is handled by the testutils.DB() setup
	// Each test gets a fresh database state
}

func (suite *APITestSuite) TestHealthEndpoint() {
	req, _ := http.NewRequest("GET", "/v1/health", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "ok", response["status"])
}

func (suite *APITestSuite) TestAuthFlow() {
	// Test user registration
	registerReq := models.UserRegistration{
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Password:  "password123",
		Timezone:  "UTC",
	}
	body, _ := json.Marshal(registerReq)

	req, _ := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	// Test user login
	loginReq := models.UserLogin{
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ = json.Marshal(loginReq)

	req, _ = http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var loginResponse map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &loginResponse)
	require.NoError(suite.T(), err)

	assert.Contains(suite.T(), loginResponse, "tokens")
	tokens, ok := loginResponse["tokens"].(map[string]any)
	require.True(suite.T(), ok)

	token, ok := tokens["access_token"].(string)
	require.True(suite.T(), ok)
	assert.NotEmpty(suite.T(), token)

	// Test authenticated endpoint
	req, _ = http.NewRequest("GET", "/v1/users/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *APITestSuite) TestUnauthorizedAccess() {
	// Test access without token
	req, _ := http.NewRequest("GET", "/v1/users/profile", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)

	// Test access with invalid token
	req, _ = http.NewRequest("GET", "/v1/users/profile", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *APITestSuite) TestInvalidLogin() {
	loginReq := models.UserLogin{
		Email:    "nonexistent@example.com",
		Password: "wrongpassword",
	}
	body, _ := json.Marshal(loginReq)

	req, _ := http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *APITestSuite) TestDuplicateUserRegistration() {
	// Register first user
	suite.registerTestUser("duplicate@example.com", "password123")

	// Try to register same user again
	registerReq := models.UserRegistration{
		Email:     "duplicate@example.com",
		FirstName: "Duplicate",
		LastName:  "User",
		Password:  "password456",
		Timezone:  "UTC",
	}
	body, _ := json.Marshal(registerReq)

	req, _ := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusConflict, w.Code)
}

func (suite *APITestSuite) TestTokenRefresh() {
	// First register and login
	registerReq := models.UserRegistration{
		Email:     "refresh@example.com",
		FirstName: "Refresh",
		LastName:  "User",
		Password:  "password123",
		Timezone:  "UTC",
	}
	body, _ := json.Marshal(registerReq)

	req, _ := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	require.Equal(suite.T(), http.StatusCreated, w.Code)

	// Login to get tokens
	loginReq := models.UserLogin{
		Email:    "refresh@example.com",
		Password: "password123",
	}
	body, _ = json.Marshal(loginReq)

	req, _ = http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	require.Equal(suite.T(), http.StatusOK, w.Code)

	var loginResponse map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &loginResponse)
	require.NoError(suite.T(), err)

	tokens, ok := loginResponse["tokens"].(map[string]any)
	require.True(suite.T(), ok)

	refreshToken, ok := tokens["refresh_token"].(string)
	require.True(suite.T(), ok)
	require.NotEmpty(suite.T(), refreshToken)

	// Use refresh token to get new access token
	refreshReq := struct {
		RefreshToken string `json:"refresh_token"`
	}{
		RefreshToken: refreshToken,
	}
	body, _ = json.Marshal(refreshReq)

	req, _ = http.NewRequest("POST", "/v1/auth/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var refreshResponse map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &refreshResponse)
	require.NoError(suite.T(), err)

	newTokens, ok := refreshResponse["tokens"].(map[string]any)
	require.True(suite.T(), ok)

	newAccessToken, ok := newTokens["access_token"].(string)
	require.True(suite.T(), ok)
	assert.NotEmpty(suite.T(), newAccessToken)
}

// Helper method to register a test user
func (suite *APITestSuite) registerTestUser(email, password string) {
	registerReq := models.UserRegistration{
		Email:     email,
		FirstName: "Test",
		LastName:  "User",
		Password:  password,
		Timezone:  "UTC",
	}
	body, _ := json.Marshal(registerReq)

	req, _ := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	require.Equal(suite.T(), http.StatusCreated, w.Code)
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
