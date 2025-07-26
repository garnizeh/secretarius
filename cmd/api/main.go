package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/garnizeh/englog/internal/auth"
	"github.com/garnizeh/englog/internal/config"
	"github.com/garnizeh/englog/internal/database"
	"github.com/gin-gonic/gin"
)

var (
	// Version will be set during build
	Version = "dev"
)

func main() {
	fmt.Printf("EngLog API Server v%s\n", Version)
	fmt.Println("🚀 Starting up...")

	ctx := context.Background()

	// Load configuration
	dbConfig := config.LoadDBConfig()
	authConfig := config.LoadAuthConfig()

	// Database connection setup
	db, err := database.NewDB(ctx, dbConfig)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	defer func() {
		db.Close()
		fmt.Println("✅ Database connection closed")
	}()

	if err := db.Check(ctx); err != nil {
		panic(fmt.Sprintf("Failed to check database connection: %v", err))
	}

	// Initialize authentication service
	authService := auth.NewAuthService(db.RDBMS(), authConfig.JWTSecretKey)

	// Start token cleanup background process
	cleanupCtx, cleanupCancel := context.WithCancel(ctx)
	defer cleanupCancel()
	go authService.StartTokenCleanup(cleanupCtx)

	// Set Gin mode based on environment
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"version":   Version,
			"timestamp": time.Now().Unix(),
			"service":   "englog-api",
		})
	})

	// Basic welcome endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to EngLog API",
			"version": Version,
			"docs":    "/swagger/",
		})
	})

	// Authentication routes
	authRoutes := router.Group("/api/v1/auth")
	{
		authRoutes.POST("/register", authService.RegisterHandler)
		authRoutes.POST("/login", authService.LoginHandler)
		authRoutes.POST("/refresh", authService.RefreshHandler)
		authRoutes.POST("/logout", authService.LogoutHandler)
	}

	// Protected routes
	apiRoutes := router.Group("/api/v1")
	apiRoutes.Use(authService.RequireAuth())
	{
		apiRoutes.GET("/me", authService.MeHandler)
		// Add other protected routes here as needed
	}

	// Get port from environment or use default
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		fmt.Printf("🌐 API Server starting on port %s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("\n🛑 Shutting down server...")

	// Cancel cleanup process
	cleanupCancel()

	// Give outstanding requests 30 seconds to complete
	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	fmt.Println("✅ Server exited")
}
