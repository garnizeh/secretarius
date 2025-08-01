package main

import (
	"fmt"
	"log"

	"github.com/garnizeh/englog/internal/ai"
	"github.com/garnizeh/englog/internal/logging"
)

func main() {
	// Create a test logger
	logger := logging.NewTestLogger()

	// Create Ollama service
	service, err := ai.NewOllamaService("http://localhost:11434", logger)
	if err != nil {
		log.Fatal(err)
	}

	// Create a sample insight request
	req := &ai.InsightRequest{
		Prompt:      "Please analyze my productivity patterns and provide actionable insights for improvement.",
		UserID:      "user-12345",
		EntryIDs:    []string{"entry-001", "entry-002", "entry-003", "entry-004", "entry-005", "entry-006"},
		InsightType: "productivity",
		Context: map[string]any{
			"time_blocks": []string{"morning", "afternoon", "evening"},
			"focus_areas": []string{"development", "meetings", "documentation"},
			"date_range": map[string]string{
				"start": "2025-07-01",
				"end":   "2025-07-31",
			},
			"performance_metrics": map[string]float64{
				"avg_daily_hours":    8.5,
				"productivity_score": 0.85,
			},
		},
	}

	// Generate enhanced prompt
	enhancedPrompt := service.BuildEnhancedPrompt(req)

	fmt.Println("=== Enhanced Prompt Example ===")
	fmt.Println(enhancedPrompt)
	fmt.Println("\n=== End of Example ===")
}
