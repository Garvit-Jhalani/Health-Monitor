package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	// Load configuration
	config, err := LoadConfig(*configPath)
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Create components
	aggregator := NewAggregator()
	reporter := &ConsoleReporter{
		SlowThresholdMs: config.SlowThreshold,
	}

	// Create channels for communication
	jobs := make(chan Checker)
	results := make(chan Result)
	done := make(chan struct{})

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create a parent context that can be canceled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start worker goroutines
	var wg sync.WaitGroup
	numWorkers := 5 // Number of concurrent workers

	fmt.Printf("Starting health monitor with %d workers\n", numWorkers)
	fmt.Printf("Monitoring %d URLs with %s interval\n", len(config.URLs), config.CheckInterval)

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			worker(ctx, workerID, jobs, results, config.TimeoutSeconds)
		}(i)
	}

	// Start result processor
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case result := <-results:
				aggregator.ProcessResult(result, reporter)
			case <-done:
				return
			}
		}
	}()

	// Start scheduler
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(config.CheckInterval)
		defer ticker.Stop()

		// Initial check immediately
		scheduleChecks(config.URLs, jobs)

		for {
			select {
			case <-ticker.C:
				scheduleChecks(config.URLs, jobs)
			case <-done:
				return
			}
		}
	}()

	// Wait for termination signal
	<-sigChan
	fmt.Println("\nReceived termination signal. Shutting down gracefully...")
	
	// Signal all goroutines to stop
	close(done)
	cancel() // Cancel the context
	
	// Wait for all goroutines to finish
	wg.Wait()
	
	// Print summary before exiting
	aggregator.PrintSummary()
	fmt.Println("Health monitor shut down successfully")
}

// scheduleChecks sends all URLs to the jobs channel
func scheduleChecks(urls []string, jobs chan<- Checker) {
	for _, url := range urls {
		select {
		case jobs <- Checker{URL: url}:
			// Job sent successfully
		default:
			// Channel is full, log and continue
			fmt.Printf("Warning: Job queue is full, skipping check for %s\n", url)
		}
	}
}

// worker processes health check jobs
func worker(ctx context.Context, id int, jobs <-chan Checker, results chan<- Result, timeoutSeconds int) {
	for {
		select {
		case checker, ok := <-jobs:
			if !ok {
				return // Channel closed
			}
			
			// Create a timeout context for this specific check
			checkCtx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
			
			// Perform the health check
			result, _ := checker.Ping(checkCtx)
			
			// Send result back
			select {
			case results <- result:
				// Result sent successfully
			case <-ctx.Done():
				cancel()
				return
			}
			
			cancel() // Clean up the timeout context
			
		case <-ctx.Done():
			return // Parent context canceled
		}
	}
}