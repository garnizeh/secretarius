package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/garnizeh/englog/internal/auth"
	"github.com/garnizeh/englog/internal/config"
	"github.com/garnizeh/englog/internal/database"
	"github.com/garnizeh/englog/internal/grpc"
	"github.com/garnizeh/englog/internal/handlers"
	"github.com/garnizeh/englog/internal/logging"
	"github.com/garnizeh/englog/internal/services"
	"github.com/gin-gonic/gin"

	_ "github.com/garnizeh/englog/api" // swagger docs
)

// @title englog API
// @version 1.0
// @description Engineering Log Management API for tracking development activities, projects, and analytics.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

var (
	// Version will be set during build
	Version = "dev"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("❌ Server error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	// Load unified configuration
	cfg := config.Load()

	// Setup structured logging
	logger := logging.NewLogger(cfg.Logging)
	logger.LogStartup("api-server", Version, map[string]any{
		"environment": cfg.Environment,
		"port":        cfg.Port,
		"host":        cfg.Host,
	})

	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Database connection setup with structured config
	db, err := database.NewDB(ctx, cfg.DB)
	if err != nil {
		logger.LogError(ctx, err, "Failed to connect to database")
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer func() {
		db.CloseDB()
		logger.Info("Database connection closed")
	}()

	if err := db.Check(ctx); err != nil {
		logger.LogError(ctx, err, "Database health check failed")
		return fmt.Errorf("database health check failed: %w", err)
	}

	logger.Info("Database connection established and healthy")

	// Initialize Redis client for rate limiting
	redisClient, err := database.NewRedisClient(cfg.Redis, logger)
	if err != nil {
		logger.LogError(ctx, err, "Failed to connect to Redis - rate limiting will use fallback mode")
		// Continue without Redis (rate limiting will fall back to no-op)
		redisClient = nil
	}
	defer func() {
		if redisClient != nil {
			database.CloseRedisClient(redisClient, logger)
		}
	}()

	authService := auth.NewAuthService(db, logger, cfg.Auth.JWTSecretKey)

	// Start token cleanup background process
	cleanupCtx, cleanupCancel := context.WithCancel(ctx)
	defer cleanupCancel()
	go func() {
		logger.WithComponent("auth").Info("Starting token cleanup service")
		authService.StartTokenCleanup(cleanupCtx)
	}()

	// Initialize all other services with logger
	projectService := services.NewProjectService(db, logger)
	logEntryService := services.NewLogEntryService(db, logger)
	analyticsService := services.NewAnalyticsService(db, logger)
	tagService := services.NewTagService(db, logger)
	userService := services.NewUserService(db, logger)

	logger.Info("All services initialized successfully")

	// Initialize gRPC server for worker communication
	grpcManager := grpc.NewManager(cfg, logger)
	if err := grpcManager.Start(); err != nil {
		logger.LogError(ctx, err, "Failed to start gRPC server")
		return fmt.Errorf("gRPC server startup failed: %w", err)
	}
	defer func() {
		if err := grpcManager.Stop(); err != nil {
			logger.LogError(ctx, err, "Error stopping gRPC server")
		}
	}()

	logger.Info("gRPC server started successfully", "port", cfg.GRPC.ServerPort)

	// Create Gin router with structured logging
	router := handlers.SetupRoutes(
		cfg,
		logger,
		redisClient,
		authService,
		logEntryService,
		projectService,
		analyticsService,
		tagService,
		userService,
		grpcManager,
	)

	// Create HTTP server with configured timeouts
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("HTTP server starting",
			"host", cfg.Host,
			"port", cfg.Port,
			"read_timeout", cfg.Server.ReadTimeout,
			"write_timeout", cfg.Server.WriteTimeout,
			"idle_timeout", cfg.Server.IdleTimeout,
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.LogError(ctx, err, "HTTP server failed to start")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutdown signal received, initiating graceful shutdown")

	// Cancel cleanup process
	cleanupCancel()

	// Give outstanding requests configured time to complete
	shutdownCtx, cancel := context.WithTimeout(ctx, cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.LogError(ctx, err, "Server forced to shutdown")
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	logger.LogShutdown("api-server", "signal", true)
	return nil
}
