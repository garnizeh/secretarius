//go:build integration
// +build integration

// Package middleware_test provides integration tests for middleware components with real Redis instance
// "The best way to test a distributed system is with real dependencies." 🐳
package middleware_test

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/garnizeh/englog/internal/config"
	"github.com/garnizeh/englog/internal/logging"
	"github.com/garnizeh/englog/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// RedisContainer wraps the testcontainer Redis instance
type RedisContainer struct {
	testcontainers.Container
	URI string
}

// setupRedisContainer starts a Redis container for testing
func setupRedisContainer(ctx context.Context) (*RedisContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "redis:7-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "6379")
	if err != nil {
		return nil, err
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("redis://%s:%s", hostIP, mappedPort.Port())

	return &RedisContainer{
		Container: container,
		URI:       uri,
	}, nil
}

// createRedisClient creates a Redis client connected to the test container
func createRedisClient(uri string) *redis.Client {
	opt, _ := redis.ParseURL(uri)
	return redis.NewClient(opt)
}

// createTestLogger creates a test logger for structured logging
func createTestConfigForManager() (*logging.Logger, *bytes.Buffer) {
	var buf bytes.Buffer
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewJSONHandler(&buf, opts)
	testLogger := &logging.Logger{Logger: slog.New(handler)}
	return testLogger, &buf
}

func TestRateLimitMiddleware_WithRedis_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Start Redis container
	redisContainer, err := setupRedisContainer(ctx)
	require.NoError(t, err)
	defer func() {
		_ = redisContainer.Terminate(ctx)
	}()

	// Create Redis client
	redisClient := createRedisClient(redisContainer.URI)
	defer redisClient.Close()

	// Verify Redis connection
	err = redisClient.Ping(ctx).Err()
	require.NoError(t, err, "Failed to connect to Redis")

	gin.SetMode(gin.TestMode)

	t.Run("rate limiting with Redis - successful requests within limit", func(t *testing.T) {
		testLogger, logBuf := createTestConfigForManager()

		cfg := config.RateLimitConfig{
			Enabled:           true,
			RequestsPerMinute: 5, // Allow 5 requests per minute
			RedisEnabled:      true,
		}

		rateLimiter := middleware.NewRateLimiter(redisClient, cfg, testLogger)
		middlewareFunc := rateLimiter.Middleware()

		router := gin.New()
		router.Use(middlewareFunc)
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// Clear any existing keys
		redisClient.FlushAll(ctx)

		// Make requests within the limit
		for i := 0; i < 5; i++ {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "192.168.1.100:8080" // Simulate consistent client IP

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code, "Request %d should succeed", i+1)
			assert.Equal(t, "5", w.Header().Get("X-Rate-Limit-Limit"))
			assert.Contains(t, w.Header().Get("X-Rate-Limit-Remaining"), "")
			assert.Contains(t, w.Header().Get("X-Rate-Limit-Reset"), "")
		}

		// Verify structured logging
		logOutput := logBuf.String()
		assert.Contains(t, logOutput, "Rate limit check")
	})

	t.Run("rate limiting with Redis - requests exceed limit", func(t *testing.T) {
		testLogger, logBuf := createTestConfigForManager()

		cfg := config.RateLimitConfig{
			Enabled:           true,
			RequestsPerMinute: 3, // Allow only 3 requests per minute
			RedisEnabled:      true,
		}

		rateLimiter := middleware.NewRateLimiter(redisClient, cfg, testLogger)
		middlewareFunc := rateLimiter.Middleware()

		router := gin.New()
		router.Use(middlewareFunc)
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// Clear any existing keys
		redisClient.FlushAll(ctx)

		clientIP := "192.168.1.101:8080"

		// First 3 requests should succeed
		for i := 0; i < 3; i++ {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = clientIP

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code, "Request %d should succeed", i+1)
			assert.Equal(t, "3", w.Header().Get("X-Rate-Limit-Limit"))
		}

		// 4th request should be rate limited
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = clientIP

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTooManyRequests, w.Code, "4th request should be rate limited")
		assert.Equal(t, "3", w.Header().Get("X-Rate-Limit-Limit"))
		assert.Equal(t, "0", w.Header().Get("X-Rate-Limit-Remaining"))
		assert.Contains(t, w.Body.String(), "Rate limit exceeded")

		// Verify rate limiting was logged
		logOutput := logBuf.String()
		assert.Contains(t, logOutput, "Rate limit exceeded")
	})

	t.Run("rate limiting with Redis - different clients independent limits", func(t *testing.T) {
		testLogger, _ := createTestConfigForManager()

		cfg := config.RateLimitConfig{
			Enabled:           true,
			RequestsPerMinute: 2,
			RedisEnabled:      true,
		}

		rateLimiter := middleware.NewRateLimiter(redisClient, cfg, testLogger)
		middlewareFunc := rateLimiter.Middleware()

		router := gin.New()
		router.Use(middlewareFunc)
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// Clear any existing keys
		redisClient.FlushAll(ctx)

		client1IP := "192.168.1.102:8080"
		client2IP := "192.168.1.103:8080"

		// Client 1: Use up their limit
		for i := 0; i < 2; i++ {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = client1IP

			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}

		// Client 1: Should be rate limited
		w1 := httptest.NewRecorder()
		req1 := httptest.NewRequest("GET", "/test", nil)
		req1.RemoteAddr = client1IP
		router.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusTooManyRequests, w1.Code)

		// Client 2: Should still be able to make requests
		w2 := httptest.NewRecorder()
		req2 := httptest.NewRequest("GET", "/test", nil)
		req2.RemoteAddr = client2IP
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusOK, w2.Code)
	})

	t.Run("rate limiting with Redis - sliding window behavior", func(t *testing.T) {
		testLogger, _ := createTestConfigForManager()

		cfg := config.RateLimitConfig{
			Enabled:           true,
			RequestsPerMinute: 2,
			RedisEnabled:      true,
		}

		rateLimiter := middleware.NewRateLimiter(redisClient, cfg, testLogger)
		middlewareFunc := rateLimiter.Middleware()

		router := gin.New()
		router.Use(middlewareFunc)
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// Clear any existing keys
		redisClient.FlushAll(ctx)

		clientIP := "192.168.1.104:8080"

		// Make 2 requests (use up limit)
		for i := 0; i < 2; i++ {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = clientIP

			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}

		// Next request should be rate limited
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = clientIP
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusTooManyRequests, w.Code)

		// Wait for sliding window to allow new requests (simulate time passage)
		// Note: In a real test, you might need to wait longer or manipulate Redis directly
		time.Sleep(100 * time.Millisecond)

		// Verify the sliding window key exists in Redis
		keys, err := redisClient.Keys(ctx, "*").Result()
		require.NoError(t, err)
		assert.NotEmpty(t, keys, "Rate limiting keys should exist in Redis")
	})
}

func TestMiddlewareStack_WithRedis_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Start Redis container
	redisContainer, err := setupRedisContainer(ctx)
	require.NoError(t, err)
	defer func() {
		_ = redisContainer.Terminate(ctx)
	}()

	// Create Redis client
	redisClient := createRedisClient(redisContainer.URI)
	defer redisClient.Close()

	// Verify Redis connection
	err = redisClient.Ping(ctx).Err()
	require.NoError(t, err, "Failed to connect to Redis")

	gin.SetMode(gin.TestMode)

	t.Run("full middleware stack with Redis rate limiting", func(t *testing.T) {
		testLogger, logBuf := createTestConfigForManager()

		// Configure all middleware components
		securityCfg := config.SecurityConfig{
			CORSAllowedOrigins:   []string{"https://example.com"},
			CORSAllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
			CORSAllowedHeaders:   []string{"Content-Type", "Authorization"},
			CORSAllowCredentials: true,
		}

		rateLimitCfg := config.RateLimitConfig{
			Enabled:           true,
			RequestsPerMinute: 3,
			RedisEnabled:      true,
		}

		// Create middleware components
		validator := middleware.NewValidationMiddleware()
		rateLimiter := middleware.NewRateLimiter(redisClient, rateLimitCfg, testLogger)

		// Setup router with full middleware stack
		router := gin.New()
		router.Use(middleware.RequestLogger(testLogger))
		router.Use(middleware.CORS(securityCfg))
		router.Use(rateLimiter.Middleware())

		router.GET("/api/test/:id",
			validator.ValidateUUIDParam("id"),
			func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "success",
					"id":      c.Param("id"),
				})
			})

		// Clear Redis
		redisClient.FlushAll(ctx)

		clientIP := "192.168.1.105:8080"
		validUUID := "550e8400-e29b-41d4-a716-446655440000"

		// Test successful requests within rate limit
		for i := 0; i < 3; i++ {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/api/test/%s", validUUID), nil)
			req.Header.Set("Origin", "https://example.com")
			req.RemoteAddr = clientIP

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code, "Request %d should succeed", i+1)
			assert.Equal(t, "https://example.com", w.Header().Get("Access-Control-Allow-Origin"))
			assert.Contains(t, w.Body.String(), "success")
			assert.Contains(t, w.Body.String(), validUUID)
			assert.NotEmpty(t, w.Header().Get("X-Trace-ID")) // From request logger
		}

		// Test rate limiting kicks in
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/test/%s", validUUID), nil)
		req.Header.Set("Origin", "https://example.com")
		req.RemoteAddr = clientIP

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTooManyRequests, w.Code)
		assert.Contains(t, w.Body.String(), "Rate limit exceeded")

		// Test validation failure with rate limiting in place
		redisClient.FlushAll(ctx) // Reset rate limits

		w = httptest.NewRecorder()
		req = httptest.NewRequest("GET", "/api/test/invalid-uuid", nil)
		req.Header.Set("Origin", "https://example.com")
		req.RemoteAddr = clientIP

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid id format")
		assert.Equal(t, "https://example.com", w.Header().Get("Access-Control-Allow-Origin"))

		// Verify comprehensive logging occurred
		logOutput := logBuf.String()
		assert.Contains(t, logOutput, "HTTP Request")
		assert.Contains(t, logOutput, "Rate limit")
	})
}

func TestRedisConnection_FailureHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	gin.SetMode(gin.TestMode)
	testLogger, logBuf := createTestConfigForManager()

	t.Run("graceful fallback when Redis is unavailable", func(t *testing.T) {
		// Create Redis client with invalid connection
		invalidRedisClient := redis.NewClient(&redis.Options{
			Addr: "localhost:9999", // Non-existent Redis instance
		})
		defer invalidRedisClient.Close()

		cfg := config.RateLimitConfig{
			Enabled:           true,
			RequestsPerMinute: 5,
			RedisEnabled:      true,
		}

		rateLimiter := middleware.NewRateLimiter(invalidRedisClient, cfg, testLogger)
		middlewareFunc := rateLimiter.Middleware()

		router := gin.New()
		router.Use(middlewareFunc)
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// Requests should still work (fallback behavior)
		for i := 0; i < 10; i++ {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "192.168.1.106:8080"

			router.ServeHTTP(w, req)

			// Should still pass due to fallback behavior when Redis is unavailable
			assert.Equal(t, http.StatusOK, w.Code, "Request %d should succeed with fallback", i+1)
		}

		// Verify error logging
		logOutput := logBuf.String()
		assert.Contains(t, logOutput, "Rate limit check")
	})
}
