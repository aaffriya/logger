# 🚀 Logger Package

A high-performance, structured logging package for Go applications with **zero external dependencies**. Built on Go's standard `log/slog` package, it provides beautiful console output, structured JSON file logging, context metadata, and configurable stack traces.

## ✨ Features

- 🎨 **Pretty Console Output** - Colorized, formatted console logging with timestamps
- 📄 **JSON File Logging** - Structured JSON output for file logging and log aggregation
- 🔍 **Stack Traces** - Configurable stack trace depth per log level
- 🏷️ **Context Metadata** - Built-in support for trace_id, user_id, and action context
- ⚡ **High Performance** - Optimized string building and memory allocation
- 🔧 **Highly Configurable** - Flexible configuration for different environments
- 📦 **slog Compatible** - Built on Go's standard `log/slog` package
- 🚫 **Zero Dependencies** - No external dependencies, only Go standard library

## 📦 Installation

```bash
go get github.com/aaffriya/logger
```

## 🚀 Quick Start

### Console Logging (Pretty Output)

```go
package main

import (
    "context"
    "log/slog"
    
    "github.com/aaffriya/logger"
    "github.com/aaffriya/logger/config"
    ctxmeta "github.com/aaffriya/logger/pkg/context"
)

func main() {
    // Configure logger
    loggerConfig := &config.LoggerConfig{
        DefaultFields: config.DefaultFieldInfo{
            Version: "v1.0.0",
            Service: "MyApp",
        },
        Stack: config.StackConfig{
            Enabled: true,
            Skip:    5,
            Depth: config.StackDepths{
                Error: 5,
                Info:  3,
                Debug: 4,
                Warn:  6,
            },
        },
        Pretty: config.PrettyConfig{
            IncludeTimestamp: true,
            IsJsonOutput:     false, // Pretty console output
        },
    }
    
    // Setup console logger
    logger.SetupConsolePrettyLogger(loggerConfig, nil)
    
    // Add context metadata
    ctx := ctxmeta.WithAction(context.Background(), "USER_LOGIN")
    ctx = ctxmeta.WithTraceID(ctx, "trace-123")
    ctx = ctxmeta.WithUserID(ctx, "user-456")
    
    // Log with context
    slog.InfoContext(ctx, "User logged in successfully", "email", "user@example.com")
    slog.WarnContext(ctx, "Rate limit approaching", "current", 95, "limit", 100)
    slog.ErrorContext(ctx, "Database connection failed", "error", "connection timeout")
}
```

### File Logging (JSON Output)

```go
package main

import (
    "context"
    "log/slog"
    "os"
    
    "github.com/aaffriya/logger"
    "github.com/aaffriya/logger/config"
    ctxmeta "github.com/aaffriya/logger/pkg/context"
)

func main() {
    loggerConfig := &config.LoggerConfig{
        DefaultFields: config.DefaultFieldInfo{
            Version: "v1.0.0",
            Service: "MyApp",
        },
        Stack: config.StackConfig{
            Enabled: true,
            Skip:    5,
            Depth: config.StackDepths{
                Error: 10,
                Info:  5,
                Debug: 3,
                Warn:  7,
            },
        },
        Pretty: config.PrettyConfig{
            IncludeTimestamp: true,
            IsJsonOutput:     true, // JSON output for files
        },
    }
    
    // Open log file
    file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        panic(err)
    }
    defer file.Close()
    
    // Setup file logger
    logger.SetupFileLogger(loggerConfig, nil, file)
    
    // Log with context
    ctx := ctxmeta.WithAction(context.Background(), "DATA_PROCESSING")
    slog.InfoContext(ctx, "Processing completed", "records", 1000, "duration", "2.5s")
}
```

## 📺 Sample Console Output

Here's what the console output looks like with pretty formatting:

```
2025-09-12 19:29:34.738 [INFO] User logged in successfully | trace-123 • user-456 • USER_LOGIN
Stack Trace:
  1. /app/main.go:45 (main)
  2. /usr/local/go/src/runtime/proc.go:285 (main)

  email=user@example.com

2025-09-12 19:29:34.740 [WARN] Rate limit approaching | trace-123 • user-456 • USER_LOGIN
Stack Trace:
  1. /app/main.go:46 (main)
  2. /usr/local/go/src/runtime/proc.go:285 (main)

  current=95
  limit=100

2025-09-12 19:29:34.741 [ERROR] Database connection failed | trace-123 • user-456 • USER_LOGIN
Stack Trace:
  1. /app/main.go:47 (main)
  2. /usr/local/go/src/runtime/proc.go:285 (main)
  3. /usr/local/go/src/runtime/asm_amd64.s:1268 (goexit)

  error=connection timeout
```

## 📄 Sample JSON File Output

File logging produces structured JSON for easy parsing and log aggregation:

```json
{"action":"USER_LOGIN","email":"user@example.com","level":"INFO","message":"User logged in successfully","service":"MyApp","timestamp":"2025-09-12T19:29:34.738+05:30","trace":["/app/main.go:45 (main)","/usr/local/go/src/runtime/proc.go:285 (main)"],"trace_id":"trace-123","user_id":"user-456","version":"v1.0.0"}

{"action":"USER_LOGIN","current":95,"level":"WARN","limit":100,"message":"Rate limit approaching","service":"MyApp","timestamp":"2025-09-12T19:29:34.740+05:30","trace":["/app/main.go:46 (main)","/usr/local/go/src/runtime/proc.go:285 (main)"],"trace_id":"trace-123","user_id":"user-456","version":"v1.0.0"}

{"action":"USER_LOGIN","error":"connection timeout","level":"ERROR","message":"Database connection failed","service":"MyApp","timestamp":"2025-09-12T19:29:34.741+05:30","trace":["/app/main.go:47 (main)","/usr/local/go/src/runtime/proc.go:285 (main)","/usr/local/go/src/runtime/asm_amd64.s:1268 (goexit)"],"trace_id":"trace-123","user_id":"user-456","version":"v1.0.0"}
```

## ⚙️ Configuration

### Complete Configuration Structure

```go
type LoggerConfig struct {
    Stack         StackConfig      `yaml:"stack" json:"stack"`
    DefaultFields DefaultFieldInfo `yaml:"default_fields" json:"default_fields"`
    Pretty        PrettyConfig     `yaml:"pretty" json:"pretty"`
}

type StackConfig struct {
    Enabled bool        `yaml:"enabled" json:"enabled"`  // Enable/disable stack traces
    Skip    int         `yaml:"skip" json:"skip"`        // Number of stack frames to skip
    Depth   StackDepths `yaml:"depth" json:"depth"`      // Stack depth per log level
}

type StackDepths struct {
    Error int `yaml:"error" json:"error"`  // Stack depth for ERROR level
    Debug int `yaml:"debug" json:"debug"`  // Stack depth for DEBUG level
    Info  int `yaml:"info" json:"info"`    // Stack depth for INFO level
    Warn  int `yaml:"warn" json:"warn"`    // Stack depth for WARN level
}

type DefaultFieldInfo struct {
    Service string `yaml:"service" json:"service"`  // Service name (added to all logs)
    Version string `yaml:"version" json:"version"`  // Service version (added to all logs)
}

type PrettyConfig struct {
    IncludeTimestamp bool `yaml:"include_timestamp" json:"include_timestamp"`  // Include timestamp in output
    IsJsonOutput     bool `yaml:"is_json_output" json:"is_json_output"`        // JSON vs pretty format
}
```

### Configuration Examples

#### Development Configuration (Pretty Console)
```go
devConfig := &config.LoggerConfig{
    DefaultFields: config.DefaultFieldInfo{
        Version: "v1.0.0",
        Service: "MyApp",
    },
    Stack: config.StackConfig{
        Enabled: true,
        Skip:    5,
        Depth: config.StackDepths{
            Error: 10,  // More stack trace for errors
            Info:  3,   // Less for info
            Debug: 5,   // Medium for debug
            Warn:  5,   // Medium for warnings
        },
    },
    Pretty: config.PrettyConfig{
        IncludeTimestamp: true,
        IsJsonOutput:     false,  // Pretty console output
    },
}
```

#### Production Configuration (JSON File)
```go
prodConfig := &config.LoggerConfig{
    DefaultFields: config.DefaultFieldInfo{
        Version: "v1.2.3",
        Service: "ProductionAPI",
    },
    Stack: config.StackConfig{
        Enabled: true,
        Skip:    5,
        Depth: config.StackDepths{
            Error: 15,  // Deep stack trace for production errors
            Info:  2,   // Minimal for info
            Debug: 0,   // No debug stack in production
            Warn:  5,   // Medium for warnings
        },
    },
    Pretty: config.PrettyConfig{
        IncludeTimestamp: true,
        IsJsonOutput:     true,   // JSON for log aggregation
    },
}
```

#### Minimal Configuration (No Stack Traces)
```go
minimalConfig := &config.LoggerConfig{
    DefaultFields: config.DefaultFieldInfo{
        Version: "v1.0.0",
        Service: "SimpleApp",
    },
    Stack: config.StackConfig{
        Enabled: false,  // Disable all stack traces
    },
    Pretty: config.PrettyConfig{
        IncludeTimestamp: false,
        IsJsonOutput:     false,
    },
}
```

## 🏷️ Context Metadata

The package provides convenient functions to add metadata to your logs:

```go
import ctxmeta "github.com/aaffriya/logger/pkg/context"

// Add individual metadata
ctx = ctxmeta.WithTraceID(ctx, "trace-123")
ctx = ctxmeta.WithUserID(ctx, "user-456") 
ctx = ctxmeta.WithAction(ctx, "USER_ACTION")

// Retrieve metadata
traceID, exists := ctxmeta.Get(ctx, "trace_id")
contextData := ctxmeta.FromContext(ctx) // Gets TraceID, UserID, Action
allData := ctxmeta.GetAll(ctx)          // Gets all key-value pairs
```

## 🎯 Real-World Usage Examples

### Web API Request Logging
```go
func handleUserCreation(w http.ResponseWriter, r *http.Request) {
    // Create request context with metadata
    requestID := generateRequestID()
    ctx := ctxmeta.WithAction(r.Context(), "CREATE_USER")
    ctx = ctxmeta.WithTraceID(ctx, requestID)
    
    // Log request received
    slog.InfoContext(ctx, "Request received", 
        "method", r.Method,
        "path", r.URL.Path,
        "remote_addr", r.RemoteAddr,
    )
    
    // Process request...
    userID, err := createUser(ctx, userData)
    if err != nil {
        slog.ErrorContext(ctx, "User creation failed", 
            "error", err.Error(),
            "email", userData.Email,
        )
        http.Error(w, "Internal Server Error", 500)
        return
    }
    
    // Log success
    ctx = ctxmeta.WithUserID(ctx, userID)
    slog.InfoContext(ctx, "User created successfully", 
        "user_id", userID,
        "email", userData.Email,
        "duration_ms", time.Since(start).Milliseconds(),
    )
}
```

### Database Operation Logging
```go
func (db *Database) CreateUser(ctx context.Context, user *User) error {
    start := time.Now()
    
    slog.InfoContext(ctx, "Starting database operation", 
        "operation", "CREATE_USER",
        "table", "users",
    )
    
    result, err := db.conn.ExecContext(ctx, query, user.Email, user.Name)
    if err != nil {
        slog.ErrorContext(ctx, "Database operation failed", 
            "operation", "CREATE_USER",
            "error", err.Error(),
            "duration_ms", time.Since(start).Milliseconds(),
        )
        return err
    }
    
    userID, _ := result.LastInsertId()
    slog.InfoContext(ctx, "Database operation completed", 
        "operation", "CREATE_USER",
        "user_id", userID,
        "duration_ms", time.Since(start).Milliseconds(),
    )
    
    return nil
}
```

## 🔧 Advanced Usage

### Custom Log Levels
```go
// Setup with custom log level
opts := &slog.HandlerOptions{
    Level: slog.LevelDebug, // Set minimum log level
}
logger.SetupConsolePrettyLogger(loggerConfig, opts)

// Use different log levels
slog.Debug("Debug information", "details", debugData)
slog.Info("General information", "status", "ok")
slog.Warn("Warning message", "threshold", 0.8)
slog.Error("Error occurred", "error", err)
```

### Structured Logging with Complex Data
```go
slog.InfoContext(ctx, "Complex operation completed", 
    "user_data", map[string]interface{}{
        "id": 12345,
        "email": "user@example.com",
        "preferences": []string{"email", "sms"},
        "metadata": map[string]interface{}{
            "last_login": time.Now(),
            "login_count": 42,
        },
    },
    "performance", map[string]interface{}{
        "duration_ms": 156,
        "memory_mb": 23.5,
        "cpu_percent": 12.3,
    },
)
```

## 🚫 Zero Dependencies

This package has **zero external dependencies** and only uses Go's standard library:

- `log/slog` - Go's structured logging package
- `context` - Context handling
- `encoding/json` - JSON marshaling for file output
- `io` - I/O operations
- `os` - File operations
- `strings` - String manipulation
- `sync` - Concurrency primitives
- `time` - Time handling
- `runtime` - Stack trace generation

## 📊 Performance

The logger is optimized for performance:

- **String Building**: Uses `strings.Builder` for efficient string concatenation
- **Memory Allocation**: Pre-allocates slices with appropriate capacity
- **Concurrency**: Thread-safe file operations with minimal locking
- **Stack Traces**: Configurable depth to balance detail vs performance
- **Context Handling**: Efficient context metadata extraction

## 🧪 Testing

The package includes comprehensive tests:

```bash
# Run all tests
go test ./...

# Run with verbose output
go test ./... -v

# Run integration tests with actual output
go test -v -run TestConsoleOutputIntegration
go test -v -run TestFileOutputIntegration
go test -v -run TestRealWorldScenario
```

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📞 Support

If you have any questions or issues, please open an issue on GitHub.

---

**Built with ❤️ using only Go standard library - Zero external dependencies!**