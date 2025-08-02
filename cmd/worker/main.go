package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/garnizeh/englog/internal/ai"
	"github.com/garnizeh/englog/internal/config"
	"github.com/garnizeh/englog/internal/logging"
	"github.com/garnizeh/englog/internal/worker"
)

var (
	// Version will be set during build
	Version = "dev"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("❌ Worker error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	// Load unified configuration
	cfg := config.Load()

	// Setup structured logging with worker service identification
	baseLogger := logging.NewLogger(cfg.Logging)
	logger := baseLogger.WithService("worker")
	logger.LogStartup("worker", Version, map[string]any{
		"environment": cfg.Environment,
		"worker_id":   cfg.Worker.WorkerID,
		"worker_name": cfg.Worker.WorkerName,
		"api_server":  cfg.GRPC.APIServerAddress,
		"tls_enabled": cfg.GRPC.TLSEnabled,
		"health_port": cfg.Worker.HealthPort,
		"ollama_url":  cfg.Worker.OllamaURL,
	})

	// Initialize AI service with logger
	aiService, err := ai.NewOllamaService(ctx, cfg.Worker.OllamaURL, logger)
	if err != nil {
		logger.LogError(ctx, err, "Failed to initialize AI service",
			logging.OperationField, "initialize_ai_service")
		return fmt.Errorf("AI service initialization failed: %w", err)
	}

	// Test AI service connection
	if err := aiService.HealthCheck(ctx); err != nil {
		logger.LogWarn(ctx, "AI service health check failed - will retry during operation",
			logging.ErrorField, err.Error(),
			logging.OperationField, "health_check")
	} else {
		logger.LogInfo(ctx, "AI service connected successfully",
			logging.OperationField, "health_check")
	}

	// Setup gRPC connection manager to API server
	connectionConfig := &worker.ConnectionConfig{
		Target:              cfg.GRPC.APIServerAddress,
		TLSEnabled:          cfg.GRPC.TLSEnabled,
		CertFile:            cfg.GRPC.TLSCertFile,
		ServerName:          cfg.GRPC.ServerName,
		HealthCheckInterval: 30 * time.Second,
		RetryConfig:         worker.DefaultRetryConfig(),
	}

	connectionManager := worker.NewConnectionManager(ctx, logger, connectionConfig)

	// Connect to API server
	logger.LogInfo(ctx, "Connecting to API server",
		logging.OperationField, "connect_to_api_server",
		"address", cfg.GRPC.APIServerAddress)
	if err := connectionManager.Connect(ctx); err != nil {
		logger.LogError(ctx, err, "Failed to connect to API server",
			logging.OperationField, "connect_to_api_server")
		return fmt.Errorf("gRPC connection failed: %w", err)
	}
	defer func() {
		connectionManager.Close(ctx)
		logger.LogInfo(ctx, "gRPC connection closed",
			logging.OperationField, "close_connection")
	}()

	// Initialize worker service with connection manager and logger
	workerService := worker.NewClient(ctx, logger, connectionManager, aiService, cfg)

	// Setup HTTP health check server
	setupHealthHandlers(workerService)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Worker.HealthPort),
		Handler: http.DefaultServeMux,
	}

	// Start services in goroutines
	var wg sync.WaitGroup
	workerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.LogInfo(ctx, "Starting HTTP health server",
			logging.OperationField, "start_http_server",
			"port", cfg.Worker.HealthPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.LogError(ctx, err, "HTTP server error",
				logging.OperationField, "http_server_runtime")
		}
	}()

	// Worker client
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.LogInfo(ctx, "Starting worker client",
			logging.OperationField, "start_worker_client")
		if err := workerService.Start(workerCtx); err != nil {
			logger.LogError(ctx, err, "Worker client error",
				logging.OperationField, "worker_client_runtime")
		}
	}()

	logger.LogInfo(ctx, "Worker services started successfully",
		logging.OperationField, "startup_complete")

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	logger.LogInfo(ctx, "Shutting down worker...",
		logging.OperationField, "shutdown_start")
	cancel()

	// Shutdown HTTP server
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 30*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.LogError(ctx, err, "HTTP server shutdown error",
			logging.OperationField, "shutdown_http_server")
	}

	// Wait for all goroutines to finish
	wg.Wait()
	logger.LogInfo(ctx, "Worker stopped successfully",
		logging.OperationField, "shutdown_complete")

	return nil
}

// setupHealthHandlers configures HTTP health check endpoints
func setupHealthHandlers(workerService *worker.Client) {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if workerService.IsHealthy(ctx) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("UNHEALTHY"))
		}
	})

	http.HandleFunc("/readiness", func(w http.ResponseWriter, r *http.Request) {
		if workerService.IsReady() {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("READY"))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("NOT_READY"))
		}
	})
}
