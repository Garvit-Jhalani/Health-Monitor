package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Aggregator collects and summarizes health check results
type Aggregator struct {
	TotalChecks   uint64
	SuccessCount  uint64
	FailureCount  uint64
	TotalDuration time.Duration
	LastResults   map[string]*Result // Cache of last results per URL
	mutex         sync.RWMutex       // Protects LastResults
}

// NewAggregator creates a new Aggregator
func NewAggregator() *Aggregator {
	return &Aggregator{
		LastResults: make(map[string]*Result),
	}
}

// ProcessResult processes a health check result
func (a *Aggregator) ProcessResult(result Result, reporter HealthReporter) {
	// Update atomic counters
	atomic.AddUint64(&a.TotalChecks, 1)
	if result.Success {
		atomic.AddUint64(&a.SuccessCount, 1)
	} else {
		atomic.AddUint64(&a.FailureCount, 1)
	}

	// Get previous state for comparison
	a.mutex.RLock()
	var previousState *Result
	if prev, exists := a.LastResults[result.URL]; exists {
		// Make a copy to avoid race conditions
		prevCopy := *prev
		previousState = &prevCopy
	}
	a.mutex.RUnlock()

	// Report the result
	if reporter != nil {
		reporter.Report(result, previousState)
	}

	// Update last known state
	a.mutex.Lock()
	resultCopy := result // Make a copy to store
	a.LastResults[result.URL] = &resultCopy
	a.mutex.Unlock()
}

// PrintSummary prints a summary of all health checks
func (a *Aggregator) PrintSummary() {
	total := atomic.LoadUint64(&a.TotalChecks)
	success := atomic.LoadUint64(&a.SuccessCount)
	failures := atomic.LoadUint64(&a.FailureCount)
	
	successRate := 0.0
	if total > 0 {
		successRate = float64(success) / float64(total) * 100
	}

	fmt.Println("\n=== Health Check Summary ===")
	fmt.Printf("Total Checks: %d\n", total)
	fmt.Printf("Successful: %d (%.2f%%)\n", success, successRate)
	fmt.Printf("Failed: %d\n", failures)
	
	fmt.Println("\nLast Known Status:")
	
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	
	for url, result := range a.LastResults {
		status := "UP"
		if !result.Success {
			status = "DOWN"
		}
		fmt.Printf("- %s: %s (%.2fms)\n", url, status, float64(result.Duration.Milliseconds()))
	}
}