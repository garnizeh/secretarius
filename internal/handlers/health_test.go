package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestHealthHandler_HealthCheck tests the health check endpoint
func TestHealthHandler_HealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful_health_check", func(t *testing.T) {
		// Setup
		handler := NewHealthHandler()
		router := gin.New()
		router.GET("/health", handler.HealthCheck)

		// Test
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check response structure
		assert.Equal(t, "healthy", response["status"])
		assert.Contains(t, response, "timestamp")
		assert.Contains(t, response, "uptime")
		assert.Equal(t, "1.0.0", response["version"])

		// Validate timestamp format
		timestamp, ok := response["timestamp"].(string)
		assert.True(t, ok)
		_, err = time.Parse(time.RFC3339, timestamp)
		assert.NoError(t, err)

		// Validate uptime is present and is a string
		uptime, ok := response["uptime"].(string)
		assert.True(t, ok)
		assert.NotEmpty(t, uptime)
	})

	t.Run("uptime_progression", func(t *testing.T) {
		// Setup
		handler := NewHealthHandler()
		router := gin.New()
		router.GET("/health", handler.HealthCheck)

		// First request
		req1, _ := http.NewRequest("GET", "/health", nil)
		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)

		var response1 map[string]any
		err := json.Unmarshal(w1.Body.Bytes(), &response1)
		assert.NoError(t, err)

		// Sleep briefly to ensure uptime changes
		time.Sleep(10 * time.Millisecond)

		// Second request
		req2, _ := http.NewRequest("GET", "/health", nil)
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		var response2 map[string]any
		err = json.Unmarshal(w2.Body.Bytes(), &response2)
		assert.NoError(t, err)

		// Both should be successful
		assert.Equal(t, http.StatusOK, w1.Code)
		assert.Equal(t, http.StatusOK, w2.Code)

		// Uptime should be different (the second one should be larger)
		uptime1 := response1["uptime"].(string)
		uptime2 := response2["uptime"].(string)
		assert.NotEqual(t, uptime1, uptime2)
	})

	t.Run("consistent_version", func(t *testing.T) {
		// Setup
		handler := NewHealthHandler()
		router := gin.New()
		router.GET("/health", handler.HealthCheck)

		// Multiple requests to ensure consistent version
		for range 3 {
			req, _ := http.NewRequest("GET", "/health", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response map[string]any
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "1.0.0", response["version"])
		}
	})
}

// TestHealthHandler_ReadinessCheck tests the readiness check endpoint
func TestHealthHandler_ReadinessCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful_readiness_check", func(t *testing.T) {
		// Setup
		handler := NewHealthHandler()
		router := gin.New()
		router.GET("/ready", handler.ReadinessCheck)

		// Test
		req, _ := http.NewRequest("GET", "/ready", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check response structure
		assert.Equal(t, "ready", response["status"])
		assert.Contains(t, response, "timestamp")

		// Validate timestamp format
		timestamp, ok := response["timestamp"].(string)
		assert.True(t, ok)
		_, err = time.Parse(time.RFC3339, timestamp)
		assert.NoError(t, err)

		// Readiness check should not include uptime or version
		assert.NotContains(t, response, "uptime")
		assert.NotContains(t, response, "version")
	})

	t.Run("multiple_readiness_checks", func(t *testing.T) {
		// Setup
		handler := NewHealthHandler()
		router := gin.New()
		router.GET("/ready", handler.ReadinessCheck)

		// Multiple rapid requests
		for range 5 {
			req, _ := http.NewRequest("GET", "/ready", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response map[string]any
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "ready", response["status"])
		}
	})

	t.Run("timestamp_format_validation", func(t *testing.T) {
		// Setup
		handler := NewHealthHandler()
		router := gin.New()
		router.GET("/ready", handler.ReadinessCheck)

		// Test
		req, _ := http.NewRequest("GET", "/ready", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Validate timestamp is in correct UTC RFC3339 format
		timestamp, ok := response["timestamp"].(string)
		assert.True(t, ok)

		parsedTime, err := time.Parse(time.RFC3339, timestamp)
		assert.NoError(t, err)

		// Timestamp should be very recent (within last few seconds)
		now := time.Now().UTC()
		timeDiff := now.Sub(parsedTime)
		assert.True(t, timeDiff >= 0, "Timestamp should not be in the future")
		assert.True(t, timeDiff < 5*time.Second, "Timestamp should be very recent")
	})
}

// TestHealthHandler_ErrorHandling tests error scenarios
func TestHealthHandler_ErrorHandling(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("wrong_http_method_health", func(t *testing.T) {
		// Setup
		handler := NewHealthHandler()
		router := gin.New()
		router.GET("/health", handler.HealthCheck)

		// Test POST instead of GET
		req, _ := http.NewRequest("POST", "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return 404 since route doesn't exist for POST
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("wrong_http_method_readiness", func(t *testing.T) {
		// Setup
		handler := NewHealthHandler()
		router := gin.New()
		router.GET("/ready", handler.ReadinessCheck)

		// Test PUT instead of GET
		req, _ := http.NewRequest("PUT", "/ready", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return 404 since route doesn't exist for PUT
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("non_existent_endpoint", func(t *testing.T) {
		// Setup
		handler := NewHealthHandler()
		router := gin.New()
		router.GET("/health", handler.HealthCheck)
		router.GET("/ready", handler.ReadinessCheck)

		// Test non-existent endpoint
		req, _ := http.NewRequest("GET", "/nonexistent", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return 404
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// TestHealthHandler_ResponseHeaders tests HTTP headers
func TestHealthHandler_ResponseHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("health_check_headers", func(t *testing.T) {
		// Setup
		handler := NewHealthHandler()
		router := gin.New()
		router.GET("/health", handler.HealthCheck)

		// Test
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	})

	t.Run("readiness_check_headers", func(t *testing.T) {
		// Setup
		handler := NewHealthHandler()
		router := gin.New()
		router.GET("/ready", handler.ReadinessCheck)

		// Test
		req, _ := http.NewRequest("GET", "/ready", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	})
}
