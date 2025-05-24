package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Checker represents a service health checker
type Checker struct {
	URL string
}

// Result represents the result of a health check
type Result struct {
	URL       string
	Timestamp time.Time
	Success   bool
	Duration  time.Duration
	Error     error
	Metrics   map[string]any // Using 'any' for mixed-type metrics
}

// Ping checks the health of a service
func (c *Checker) Ping(ctx context.Context) (Result, error) {
	result := Result{
		URL:       c.URL,
		Timestamp: time.Now(),
		Metrics:   make(map[string]any),
	}

	// Create a request with the provided context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.URL, nil)
	if err != nil {
		result.Error = fmt.Errorf("failed to create request for %s: %w", c.URL, err)
		return result, result.Error
	}

	// Measure response time
	startTime := time.Now()
	resp, err := http.DefaultClient.Do(req)
	duration := time.Since(startTime)
	result.Duration = duration

	// Handle request errors
	if err != nil {
		result.Success = false
		result.Error = fmt.Errorf("request to %s failed: %w", c.URL, err)
		return result, result.Error
	}
	defer resp.Body.Close()

	// Check if status code indicates success (2xx)
	result.Success = resp.StatusCode >= 200 && resp.StatusCode < 300
	result.Metrics["statusCode"] = resp.StatusCode
	result.Metrics["timeMs"] = duration.Milliseconds()

	if !result.Success {
		result.Error = fmt.Errorf("service %s returned non-success status: %d", c.URL, resp.StatusCode)
	}

	return result, result.Error
}