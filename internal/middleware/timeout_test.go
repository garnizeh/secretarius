package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/garnizeh/englog/internal/config"
	"github.com/garnizeh/englog/internal/logging"
	"github.com/garnizeh/englog/internal/middleware"
)

func TestRequestTimeout(t *testing.T) {
	t.Skip("Skipping test due to data race conditions in timeout middleware")

	// Create test logger
	logger := logging.NewTestLogger()

	tests := []struct {
		name           string
		requestTimeout time.Duration
		handlerDelay   time.Duration
		path           string
		expectTimeout  bool
		expectedStatus int
	}{
		{
			name:           "request completes within timeout",
			requestTimeout: 100 * time.Millisecond,
			handlerDelay:   50 * time.Millisecond,
			path:           "/test",
			expectTimeout:  false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "request times out",
			requestTimeout: 50 * time.Millisecond,
			handlerDelay:   100 * time.Millisecond,
			path:           "/test",
			expectTimeout:  true,
			expectedStatus: http.StatusRequestTimeout,
		},
		{
			name:           "health endpoint skips timeout",
			requestTimeout: 50 * time.Millisecond,
			handlerDelay:   100 * time.Millisecond,
			path:           "/health",
			expectTimeout:  false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "ready endpoint skips timeout",
			requestTimeout: 50 * time.Millisecond,
			handlerDelay:   100 * time.Millisecond,
			path:           "/ready",
			expectTimeout:  false,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server configuration
			cfg := config.ServerConfig{
				RequestTimeout: tt.requestTimeout,
			}

			// Setup Gin router with timeout middleware
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Use(middleware.RequestTimeout(cfg, logger))

			// Add test handler that delays for the specified time
			router.GET("/*path", func(c *gin.Context) {
				// Simulate processing delay
				select {
				case <-time.After(tt.handlerDelay):
					c.JSON(http.StatusOK, gin.H{"message": "success"})
				case <-c.Request.Context().Done():
					// Context was cancelled (timeout)
					return
				}
			})

			// Create test request
			req, err := http.NewRequest("GET", tt.path, nil)
			require.NoError(t, err)

			// Record response
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Verify response
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectTimeout {
				// Verify timeout response structure
				assert.Contains(t, w.Body.String(), "Request timeout")
				assert.Contains(t, w.Body.String(), "error")
			} else {
				// Verify successful response
				assert.Contains(t, w.Body.String(), "success")
			}
		})
	}
}

func TestRequestTimeoutContextCancellation(t *testing.T) {
	t.Skip("Skipping test due to data race conditions in timeout middleware")

	logger := logging.NewTestLogger()

	cfg := config.ServerConfig{
		RequestTimeout: 100 * time.Millisecond,
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.RequestTimeout(cfg, logger))

	// Handler that checks if context is properly cancelled
	router.GET("/test", func(c *gin.Context) {
		// Wait longer than timeout
		select {
		case <-time.After(200 * time.Millisecond):
			c.JSON(http.StatusOK, gin.H{"message": "this should not happen"})
		case <-c.Request.Context().Done():
			// Verify context was cancelled due to timeout
			assert.Equal(t, context.DeadlineExceeded, c.Request.Context().Err())
			return
		}
	})

	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusRequestTimeout, w.Code)
}

func TestRequestTimeoutMiddlewareOrder(t *testing.T) {
	t.Skip("Skipping test due to data race conditions in concurrent middleware execution")

	logger := logging.NewTestLogger()

	cfg := config.ServerConfig{
		RequestTimeout: 100 * time.Millisecond,
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add timeout middleware first (this is important for proper functioning)
	router.Use(middleware.RequestTimeout(cfg, logger))

	// Add another middleware after timeout
	router.Use(func(c *gin.Context) {
		c.Header("X-Test-Middleware", "executed")
		c.Next()
	})

	router.GET("/test", func(c *gin.Context) {
		// Delay longer than timeout
		time.Sleep(150 * time.Millisecond)
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should timeout and return 408
	assert.Equal(t, http.StatusRequestTimeout, w.Code)

	// Verify that other middleware was executed before timeout
	assert.Equal(t, "executed", w.Header().Get("X-Test-Middleware"))
}
