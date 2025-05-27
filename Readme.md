# 🩺 Go Health Monitor

A lightweight and concurrent health monitoring system built in Go. It pings web services at regular intervals, aggregates the results, and reports if services are slow or down — all with graceful shutdown support and concurrency best practices.

---

## 🚀 Features

- Periodic health checks on HTTP services
- Response time measurement and error reporting
- Console alerts for slow or down services
- Concurrent execution using goroutines and channels
- Graceful shutdown on `Ctrl+C`
- Configurable via JSON or environment variables

---

## 📁 Project Structure

```
health-monitor/
├── config.go       # Load URLs from JSON/env
├── checker.go      # Health check logic
├── reporter.go     # Reporting interface & console output
├── aggregator.go   # Aggregate and summarize results
├── main.go         # App orchestration and signal handling
├── config.json     # Sample configuration
├── go.mod          # Go module file
└── README.md       # This file
```

---

## 🧠 Core Go Concepts Used

| Concept              | Usage Description |
|----------------------|-------------------|
| **Error handling**   | In `checker.go`, errors are wrapped with context. |
| **Struct methods**   | `Checker` type with `Ping(ctx)` method. |
| **Interfaces**       | `HealthReporter` allows flexible output formats. |
| **Channels**         | For job distribution and result collection. |
| **Goroutines**       | For concurrent health checks. |
| **Atomic counters**  | For tracking success/failure counts. |
| **Mutex/RWMutex**    | Protects shared state (`LastResults` map). |
| **Context API**      | For timeout and cancellation in HTTP requests. |
| **Graceful Shutdown**| Done via OS signal and `done` channel. |

---

## ⚙️ Configuration

### ✅ JSON Example (`config.json`)

```json
{
  "urls": [
    "https://google.com",
    "https://github.com",
    "https://example.com",
    "https://nonexistentsite123456789.com"
  ],
  "checkIntervalSeconds": 15,
  "timeoutSeconds": 5,
  "slowThresholdMs": 300
}
```

### 🌐 Environment Variables

```bash
export HEALTH_MONITOR_URLS="https://google.com,https://github.com"
export HEALTH_MONITOR_INTERVAL="15"
export HEALTH_MONITOR_TIMEOUT="5"
export HEALTH_MONITOR_SLOW_THRESHOLD="300"
```

> 🧩 Precedence: Environment Variables > Config File > Default Values

---

## 🏃‍♂️ Running the Program

### 💻 Build & Run

```bash
# Build
go build -o health-monitor

# Run with default config
./health-monitor

# Run with custom config
./health-monitor -config=config.json
```

---

## 💬 Sample Console Output

```
Starting health monitor with 5 workers
Monitoring 4 URLs with 30s interval

[2024-01-15T10:30:00Z] https://google.com        ✅ UP    | 45.23ms
[2024-01-15T10:30:01Z] https://github.com        ✅ UP    | 123.45ms
[2024-01-15T10:30:02Z] https://example.com       ⚠️ SLOW | 678.90ms
[2024-01-15T10:30:03Z] https://httpstat.us/500   ❌ DOWN | 234.56ms 🔄 NEWLY DOWN
```

---

## 📉 Graceful Shutdown

Press `Ctrl+C` to stop:

```
^C Received termination signal. Shutting down gracefully...

=== Health Check Summary ===
Total Checks: 24
Successful: 18 (75.00%)
Failed: 6

Last Known Status:
- https://google.com: UP (45.23ms)
- https://github.com: UP (123.45ms)
- https://example.com: UP (678.90ms)
- https://httpstat.us/500: DOWN (234.56ms)

Health monitor shut down successfully.
```

---

## 🧪 Learning Takeaways

- Concurrency via goroutines and channels
- Interface-based extensibility
- Thread-safe state management (mutex, atomic)
- Context-aware timeout handling
- Graceful shutdown patterns in Go
- Clean architecture with modular design

---

## 🌱 Future Extensions

- ✅ Slack/Discord integration for alerts
- 📊 Prometheus endpoint for metrics
- 🌐 Web UI for live status
- 💾 Database for historical monitoring data
- 🔍 Custom checkers (e.g., DB, TCP, DNS)

---

## 📌 Quick Start

```bash
# 1. Initialize the module
go mod init health-monitor

# 2. Add Go files (config.go, checker.go, etc.)

# 3. Run
go run .
```

---

> Built to help you understand Go's concurrency model and monitoring principles through a practical, hands-on project.
