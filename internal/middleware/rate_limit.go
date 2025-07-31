package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/garnizeh/englog/internal/config"
	"github.com/garnizeh/englog/internal/logging"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimiter provides Redis-based rate limiting with sliding window
type RateLimiter struct {
	redis  *redis.Client
	config config.RateLimitConfig
	logger *logging.Logger
}

// NewRateLimiter creates a new rate limiter instance
func NewRateLimiter(redisClient *redis.Client, cfg config.RateLimitConfig, logger *logging.Logger) *RateLimiter {
	return &RateLimiter{
		redis:  redisClient,
		config: cfg,
		logger: logger,
	}
}

// Middleware returns a Gin middleware function for rate limiting
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	if !rl.config.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		if !rl.config.RedisEnabled {
			// Fallback to memory-based rate limiting (no-op for now)
			c.Next()
			return
		}

		clientIP := c.ClientIP()
		key := fmt.Sprintf("rate_limit:%s", clientIP)

		allowed, remaining, resetTime, err := rl.checkRateLimit(c.Request.Context(), key)
		if err != nil {
			rl.logger.LogError(c.Request.Context(), err, "Rate limit check failed",
				"client_ip", clientIP,
				"path", c.Request.URL.Path,
			)
			// Allow request if rate limiting fails
			c.Next()
			return
		}

		rl.logger.WithContext(c.Request.Context()).Info("Rate limit check",
			"client_ip", clientIP,
			"path", c.Request.URL.Path,
			"allowed", allowed,
			"remaining", remaining,
		)

		// Add rate limit headers
		c.Header("X-Rate-Limit-Limit", strconv.Itoa(rl.config.RequestsPerMinute))
		c.Header("X-Rate-Limit-Remaining", strconv.Itoa(remaining))
		c.Header("X-Rate-Limit-Reset", strconv.FormatInt(resetTime, 10))

		if !allowed {
			rl.logger.WithContext(c.Request.Context()).Warn("Rate limit exceeded",
				"client_ip", clientIP,
				"path", c.Request.URL.Path,
				"remaining", remaining,
			)

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"retry_after": resetTime - time.Now().Unix(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// checkRateLimit implements sliding window rate limiting with Redis
func (rl *RateLimiter) checkRateLimit(ctx context.Context, key string) (allowed bool, remaining int, resetTime int64, err error) {
	now := time.Now()
	windowStart := now.Add(-time.Minute) // Look back 1 minute
	resetTime = now.Add(time.Minute).Unix()

	// Use sliding window algorithm with Redis
	pipeline := rl.redis.Pipeline()

	// Remove old entries (older than 1 minute)
	pipeline.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart.Unix()))

	// Count current requests in the window
	countCmd := pipeline.ZCard(ctx, key)

	_, err = pipeline.Exec(ctx)
	if err != nil {
		return false, 0, resetTime, err
	}

	count := int(countCmd.Val())
	allowed = count < rl.config.RequestsPerMinute
	remaining = rl.config.RequestsPerMinute - count

	if allowed {
		// Only add the request if it's allowed
		rl.redis.ZAdd(ctx, key, redis.Z{
			Score:  float64(now.Unix()),
			Member: fmt.Sprintf("%d", now.UnixNano()),
		})
		// Set expiration for the key
		rl.redis.Expire(ctx, key, time.Minute*2)
		remaining-- // Decrement after adding
	}

	if remaining < 0 {
		remaining = 0
	}

	return allowed, remaining, resetTime, nil
}
