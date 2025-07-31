package database

import (
	"context"
	"fmt"
	"time"

	"github.com/garnizeh/englog/internal/config"
	"github.com/garnizeh/englog/internal/logging"
	"github.com/redis/go-redis/v9"
)

// NewRedisClient creates a new Redis client instance with the provided configuration
func NewRedisClient(cfg config.RedisConfig, logger *logging.Logger) (*redis.Client, error) {
	ctx := context.Background()

	// Create Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		DialTimeout:  10 * time.Second,
	})

	// Test the connection
	err := rdb.Ping(ctx).Err()
	if err != nil {
		logger.LogError(ctx, err, "Failed to connect to Redis",
			"host", cfg.Host,
			"port", cfg.Port,
			"db", cfg.DB,
		)
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	logger.Info("Redis connection established successfully",
		"host", cfg.Host,
		"port", cfg.Port,
		"db", cfg.DB,
		"pool_size", cfg.PoolSize,
	)

	return rdb, nil
}

// CloseRedisClient gracefully closes the Redis client connection
func CloseRedisClient(client *redis.Client, logger *logging.Logger) {
	if client != nil {
		if err := client.Close(); err != nil {
			logger.LogError(context.Background(), err, "Error closing Redis connection")
		} else {
			logger.Info("Redis connection closed successfully")
		}
	}
}
