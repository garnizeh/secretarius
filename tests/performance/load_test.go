//go:build performance
// +build performance

package performance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

// TestAPILoadTesting performs load testing against the API
// "Performance is not about being fast, it's about being consistently reliable." ⚡
func TestAPILoadTesting(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	// Test configuration
	rate := vegeta.Rate{Freq: 100, Per: time.Second}
	duration := 30 * time.Second

	// Health endpoint load test
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    "http://localhost:8080/v1/health",
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
	})

	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Health Check Load Test") {
		metrics.Add(res)
	}
	metrics.Close()

	// Assert performance requirements
	if metrics.Success < 0.99 {
		t.Errorf("Success rate too low: %.2f", metrics.Success)
	}

	if metrics.Latencies.P99 > 500*time.Millisecond {
		t.Errorf("P99 latency too high: %v", metrics.Latencies.P99)
	}

	t.Logf("Health endpoint load test results:")
	t.Logf("  Success rate: %.2f%%", metrics.Success*100)
	t.Logf("  P50 latency: %v", metrics.Latencies.P50)
	t.Logf("  P95 latency: %v", metrics.Latencies.P95)
	t.Logf("  P99 latency: %v", metrics.Latencies.P99)
	t.Logf("  Max latency: %v", metrics.Latencies.Max)
	t.Logf("  Requests/second: %.2f", float64(metrics.Requests)/duration.Seconds())
}

// BenchmarkUserRegistration benchmarks user registration endpoint
func BenchmarkUserRegistration(b *testing.B) {
	endpoint := "http://localhost:8080/v1/auth/register"

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			counter++

			// Create unique user data for each request
			payload := map[string]string{
				"email":      fmt.Sprintf("bench-user-%d@example.com", counter),
				"username":   fmt.Sprintf("benchuser%d", counter),
				"password":   "benchmarkpassword123",
				"first_name": "Benchmark",
				"last_name":  "User",
			}

			body, err := json.Marshal(payload)
			require.NoError(b, err)

			req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
			if err != nil {
				b.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{Timeout: 10 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				b.Fatal(err)
			}
			resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				b.Fatalf("Expected status 201, got %d", resp.StatusCode)
			}
		}
	})
}

// BenchmarkUserLogin benchmarks user login endpoint
func BenchmarkUserLogin(b *testing.B) {
	// Note: This assumes there's a test user already created
	// In a real scenario, you'd set up test data first
	endpoint := "http://localhost:8080/v1/auth/login"

	payload := map[string]string{
		"email":    "benchmark@example.com",
		"password": "benchmarkpassword123",
	}
	body, err := json.Marshal(payload)
	require.NoError(b, err)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
			if err != nil {
				b.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{Timeout: 10 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				b.Fatal(err)
			}
			resp.Body.Close()
		}
	})
}

// BenchmarkHealthCheck benchmarks health check endpoint
func BenchmarkHealthCheck(b *testing.B) {
	endpoint := "http://localhost:8080/v1/health"

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		client := &http.Client{Timeout: 5 * time.Second}
		for pb.Next() {
			resp, err := client.Get(endpoint)
			if err != nil {
				b.Fatal(err)
			}
			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				b.Fatalf("Expected status 200, got %d", resp.StatusCode)
			}
		}
	})
}

// TestDatabasePerformance tests database query performance
func TestDatabasePerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database performance test in short mode")
	}

	// This is a placeholder for database performance tests
	// In a real implementation, you would:
	// 1. Set up a test database with sample data
	// 2. Run queries with different data sizes
	// 3. Measure query execution times
	// 4. Assert performance requirements

	t.Log("Database performance test placeholder")
	t.Log("TODO: Implement database query performance tests")
}

// TestMemoryUsage tests memory usage under load
func TestMemoryUsage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping memory usage test in short mode")
	}

	// This is a placeholder for memory usage tests
	// In a real implementation, you would:
	// 1. Start the application
	// 2. Generate load
	// 3. Monitor memory usage over time
	// 4. Assert memory doesn't exceed limits

	t.Log("Memory usage test placeholder")
	t.Log("TODO: Implement memory usage monitoring tests")
}
