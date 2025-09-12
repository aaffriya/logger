package logger

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/aaffriya/logger"
	"github.com/aaffriya/logger/config"
	ctxmeta "github.com/aaffriya/logger/pkg/context"
)

// TestConsoleOutputIntegration demonstrates actual console output
func TestConsoleOutputIntegration(t *testing.T) {
	t.Log("=== CONSOLE LOGGER INTEGRATION TEST ===")
	
	loggerConfig := &config.LoggerConfig{
		DefaultFields: config.DefaultFieldInfo{
			Version: "v1.0.0",
			Service: "IntegrationTest",
		},
		Stack: config.StackConfig{
			Enabled: true,
			Skip:    5,
			Depth: config.StackDepths{
				Error: 3,
				Info:  2,
				Debug: 1,
				Warn:  2,
			},
		},
		Pretty: config.PrettyConfig{
			IncludeTimestamp: false,
			IsJsonOutput:     true, // Pretty console output
		},
	}

	// Setup console logger (outputs to stdout)
	logger.SetupConsolePrettyLogger(loggerConfig, nil)
	
	// Create context with metadata
	ctx := ctxmeta.WithAction(context.Background(), "INTEGRATION_TEST")
	ctx = ctxmeta.WithTraceID(ctx, "trace-integration-123")
	ctx = ctxmeta.WithUserID(ctx, "test-user-456")

	t.Log("Console output should appear below:")
	
	// Test different log levels with context
	slog.InfoContext(ctx, "Integration test started", "test_type", "console", "timestamp", time.Now().Format(time.RFC3339))
	slog.WarnContext(ctx, "This is a warning message", "warning_code", "W001", "severity", "medium")
	slog.ErrorContext(ctx, "This is an error message", "error_code", "E001", "details", "Integration test error")
	
	// Test without context
	slog.Info("Message without context", "standalone", true)
	
	t.Log("=== END CONSOLE OUTPUT ===")
}

// TestFileOutputIntegration demonstrates actual file logging
func TestFileOutputIntegration(t *testing.T) {
	t.Log("=== FILE LOGGER INTEGRATION TEST ===")
	
	// Create test log file
	testFile := "integration_test.log"
	file, err := os.OpenFile(testFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		t.Fatalf("Failed to create test log file: %v", err)
	}
	defer func() {
		file.Close()
		// Clean up test file
		os.Remove(testFile)
	}()

	loggerConfig := &config.LoggerConfig{
		DefaultFields: config.DefaultFieldInfo{
			Version: "v2.0.0",
			Service: "FileIntegrationTest",
		},
		Stack: config.StackConfig{
			Enabled: true,
			Skip:    5,
			Depth: config.StackDepths{
				Error: 5,
				Info:  3,
				Debug: 2,
				Warn:  4,
			},
		},
		Pretty: config.PrettyConfig{
			IncludeTimestamp: true,
			IsJsonOutput:     true, // JSON file output
		},
	}

	// Setup file logger
	logger.SetupFileLogger(loggerConfig, nil, file)
	
	// Create context with metadata
	ctx := ctxmeta.WithAction(context.Background(), "FILE_TEST")
	ctx = ctxmeta.WithTraceID(ctx, "trace-file-789")
	ctx = ctxmeta.WithUserID(ctx, "file-user-123")

	// Test various log scenarios
	slog.InfoContext(ctx, "File integration test started", 
		"test_type", "file_logging",
		"file_name", testFile,
		"timestamp", time.Now().Unix(),
	)
	
	slog.WarnContext(ctx, "File warning message", 
		"warning_type", "disk_space",
		"threshold", 85.5,
		"current_usage", 92.3,
	)
	
	slog.ErrorContext(ctx, "File error occurred", 
		"error_type", "io_error",
		"file_path", "/tmp/test.log",
		"error_details", map[string]interface{}{
			"code": 500,
			"message": "Permission denied",
			"retry_count": 3,
		},
	)

	// Test complex data structures
	slog.InfoContext(ctx, "Complex data logging", 
		"user_data", map[string]interface{}{
			"id": 12345,
			"name": "John Doe",
			"email": "john@example.com",
			"preferences": []string{"email", "sms", "push"},
			"metadata": map[string]interface{}{
				"last_login": "2024-01-15T10:30:00Z",
				"login_count": 42,
				"is_premium": true,
			},
		},
		"request_info", map[string]interface{}{
			"method": "POST",
			"url": "/api/v1/users",
			"headers": map[string]string{
				"Content-Type": "application/json",
				"User-Agent": "TestClient/1.0",
			},
			"body_size": 1024,
		},
	)

	// Force file sync
	file.Sync()
	
	// Read and display file contents
	file.Seek(0, 0)
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test log file: %v", err)
	}
	
	t.Logf("=== FILE CONTENT (%s) ===", testFile)
	t.Logf("%s", string(content))
	t.Log("=== END FILE CONTENT ===")
	
	// Verify file contains expected elements
	contentStr := string(content)
	expectedElements := []string{
		"File integration test started",
		"FILE_TEST",
		"trace-file-789",
		"file-user-123",
		"FileIntegrationTest",
		"v2.0.0",
		"Complex data logging",
		"john@example.com",
	}
	
	for _, element := range expectedElements {
		if !contains(contentStr, element) {
			t.Errorf("Expected file content to contain '%s'", element)
		}
	}
}

// TestRealWorldScenario demonstrates a realistic usage scenario
func TestRealWorldScenario(t *testing.T) {
	t.Log("=== REAL WORLD SCENARIO TEST ===")
	
	// Simulate a web application scenario
	loggerConfig := &config.LoggerConfig{
		DefaultFields: config.DefaultFieldInfo{
			Version: "v1.2.3",
			Service: "WebAPI",
		},
		Stack: config.StackConfig{
			Enabled: true,
			Skip:    5,
			Depth: config.StackDepths{
				Error: 10,
				Info:  3,
				Debug: 2,
				Warn:  5,
			},
		},
		Pretty: config.PrettyConfig{
			IncludeTimestamp: true,
			IsJsonOutput:     false,
		},
	}

	// Setup console logger for this scenario
	logger.SetupConsolePrettyLogger(loggerConfig, nil)
	
	t.Log("Simulating web request processing:")
	
	// Simulate incoming request
	requestID := "req-" + time.Now().Format("20060102-150405")
	ctx := ctxmeta.WithAction(context.Background(), "HTTP_REQUEST")
	ctx = ctxmeta.WithTraceID(ctx, requestID)
	ctx = ctxmeta.WithUserID(ctx, "user-12345")
	
	// Request received
	slog.InfoContext(ctx, "Request received", 
		"method", "POST",
		"path", "/api/v1/users",
		"remote_addr", "192.168.1.100",
		"user_agent", "Mozilla/5.0 (compatible; TestClient/1.0)",
	)
	
	// Processing
	slog.InfoContext(ctx, "Processing user creation", 
		"email", "newuser@example.com",
		"role", "standard",
	)
	
	// Database operation
	slog.InfoContext(ctx, "Database operation completed", 
		"operation", "INSERT",
		"table", "users",
		"duration_ms", 45,
	)
	
	// Warning scenario
	slog.WarnContext(ctx, "Rate limit approaching", 
		"current_requests", 95,
		"limit", 100,
		"window", "1m",
	)
	
	// Success response
	slog.InfoContext(ctx, "Request completed successfully", 
		"status_code", 201,
		"response_time_ms", 156,
		"user_id", 67890,
	)
	
	t.Log("=== END REAL WORLD SCENARIO ===")
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		containsAt(s, substr))))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}