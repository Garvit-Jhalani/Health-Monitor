package main

import (
	"fmt"
	"time"
)

// HealthReporter defines the interface for reporting health check results
type HealthReporter interface {
	Report(result Result, previousState *Result) error
}

// ConsoleReporter implements HealthReporter for console output
type ConsoleReporter struct {
	SlowThresholdMs int
}

// Report prints health check results to the console
func (r *ConsoleReporter) Report(result Result, previousState *Result) error {
	// Format timestamp
	timestamp := result.Timestamp.Format(time.RFC3339)

	// Basic status message
	statusMsg := "✅ UP"
	if !result.Success {
		statusMsg = "❌ DOWN"
	}

	// Response time info
	responseTime := fmt.Sprintf("%.2fms", float64(result.Duration.Milliseconds()))
	if result.Success && result.Duration.Milliseconds() > int64(r.SlowThresholdMs) {
		responseTime = fmt.Sprintf("⚠️ %s (slow)", responseTime)
	}

	// Check for state change
	stateChanged := ""
	if previousState != nil && previousState.Success != result.Success {
		if result.Success {
			stateChanged = "🔄 RECOVERED"
		} else {
			stateChanged = "🔄 NEWLY DOWN"
		}
	}

	// Print the report
	fmt.Printf("[%s] %s | %s | %s %s\n", 
		timestamp, 
		result.URL, 
		statusMsg, 
		responseTime, 
		stateChanged,
	)

	// Print error details if any
	if result.Error != nil {
		fmt.Printf("  Error: %v\n", result.Error)
	}

	return nil
}