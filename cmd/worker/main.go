package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	// Version will be set during build
	Version = "dev"
)

func main() {
	fmt.Printf("EngLog Worker Server v%s\n", Version)
	fmt.Println("⚙️  Starting up...")

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start worker components
	go startTaskProcessor(ctx)
	go startGRPCClient(ctx)
	go startScheduler(ctx)

	fmt.Println("🔧 Worker Server started successfully")
	fmt.Println("📊 Ready to process tasks...")

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("\n🛑 Shutting down worker...")

	// Cancel context to stop all goroutines
	cancel()

	// Give workers time to finish current tasks
	time.Sleep(5 * time.Second)

	fmt.Println("✅ Worker exited")
}

func startTaskProcessor(ctx context.Context) {
	fmt.Println("📋 Task processor started")
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("📋 Task processor stopped")
			return
		case <-ticker.C:
			// Placeholder for task processing logic
			fmt.Println("🔄 Processing background tasks...")
		}
	}
}

func startGRPCClient(ctx context.Context) {
	fmt.Println("🌐 gRPC client started")

	// Placeholder for gRPC client logic
	for {
		select {
		case <-ctx.Done():
			fmt.Println("🌐 gRPC client stopped")
			return
		default:
			// Keep the client running
			time.Sleep(1 * time.Second)
		}
	}
}

func startScheduler(ctx context.Context) {
	fmt.Println("⏰ Scheduler started")
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("⏰ Scheduler stopped")
			return
		case <-ticker.C:
			// Placeholder for scheduled jobs
			fmt.Println("📅 Running scheduled tasks...")
		}
	}
}
