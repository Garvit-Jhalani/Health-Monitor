package main

import (
	"encoding/json"
	"os"
	"strings"
	"time"
)

// Config holds the application configuration
type Config struct {
	URLs           []string      `json:"urls"`
	CheckInterval  time.Duration `json:"checkIntervalSeconds"`
	TimeoutSeconds int           `json:"timeoutSeconds"`
	SlowThreshold  int           `json:"slowThresholdMs"`
}

// LoadConfig loads configuration from a JSON file or environment variables
func LoadConfig(configPath string) (*Config, error) {
	// Default configuration
	config := &Config{
		// URLs:           []string{"https://nonexistentsite123456789.com", "https://intellezaai.com", "https://example.com"},
		URLs:           []string{"https://example.com", "https://google.com", "https://newperasdfja.com"},
		CheckInterval:  30 * time.Second,
		TimeoutSeconds: 5,
		SlowThreshold:  500, // 500ms
	}

	// Try to load from file if provided
	if configPath != "" {
		file, err := os.Open(configPath)
		if err == nil {
			defer file.Close()
			decoder := json.NewDecoder(file)
			if err := decoder.Decode(config); err != nil {
				return nil, err
			}
		}
	}

	// Override with environment variables if present
	if urls := os.Getenv("HEALTH_MONITOR_URLS"); urls != "" {
		config.URLs = strings.Split(urls, ",")
	}

	if interval := os.Getenv("HEALTH_MONITOR_INTERVAL"); interval != "" {
		if seconds, err := time.ParseDuration(interval); err == nil {
			config.CheckInterval = seconds
		}
	}

	if timeout := os.Getenv("HEALTH_MONITOR_TIMEOUT"); timeout != "" {
		if seconds, err := time.ParseDuration(timeout); err == nil {
			config.TimeoutSeconds = int(seconds.Seconds())
		}
	}

	if threshold := os.Getenv("HEALTH_MONITOR_SLOW_THRESHOLD"); threshold != "" {
		if ms, err := time.ParseDuration(threshold); err == nil {
			config.SlowThreshold = int(ms.Milliseconds())
		}
	}

	return config, nil
}