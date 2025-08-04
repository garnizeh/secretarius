package handlers

import (
	"os"
	"testing"

	"github.com/garnizeh/englog/internal/testutils"
)

// TestMain manages the lifecycle of shared test resources
func TestMain(m *testing.M) {
	// Run tests
	exitCode := m.Run()

	// Cleanup shared resources
	testutils.CleanupSharedResources()

	// Exit with the test result code
	os.Exit(exitCode)
}
