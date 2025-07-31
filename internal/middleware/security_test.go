package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/garnizeh/englog/internal/config"
	"github.com/garnizeh/englog/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSecurityHeaders(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	cfg := config.SecurityConfig{
		SecurityHeaders: config.SecurityHeadersConfig{
			ContentTypeOptions:      "nosniff",
			FrameOptions:            "DENY",
			XSSProtection:           "1; mode=block",
			StrictTransportSecurity: "max-age=31536000",
			ContentSecurityPolicy:   "default-src 'self'",
			ReferrerPolicy:          "strict-origin",
		},
	}

	router := gin.New()
	router.Use(middleware.SecurityHeaders(cfg))

	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test"})
	})

	// Test request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))
	assert.Equal(t, "max-age=31536000", w.Header().Get("Strict-Transport-Security"))
	assert.Equal(t, "default-src 'self'", w.Header().Get("Content-Security-Policy"))
	assert.Equal(t, "strict-origin", w.Header().Get("Referrer-Policy"))
}

func TestCORS(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	cfg := config.SecurityConfig{
		CORSAllowedOrigins:   []string{"https://example.com", "https://app.example.com"},
		CORSAllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		CORSAllowedHeaders:   []string{"Content-Type", "Authorization"},
		CORSAllowCredentials: true,
	}

	router := gin.New()
	router.Use(middleware.CORS(cfg))

	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test"})
	})

	t.Run("allowed origin", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "https://example.com")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "https://example.com", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
	})

	t.Run("options request", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/test", nil)
		req.Header.Set("Origin", "https://example.com")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 204, w.Code)
		assert.Equal(t, "https://example.com", w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("disallowed origin", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "https://malicious.com")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
	})
}

// Performance/benchmark tests
func BenchmarkCORSMiddleware(b *testing.B) {
	cfg := config.SecurityConfig{
		CORSAllowedOrigins:   []string{"https://example.com", "https://app.example.com"},
		CORSAllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		CORSAllowedHeaders:   []string{"Content-Type", "Authorization"},
		CORSAllowCredentials: true,
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.CORS(cfg))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)
	}
}
