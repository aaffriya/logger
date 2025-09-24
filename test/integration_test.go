package logger

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/aaffriya/logger/config"
	customhandler "github.com/aaffriya/logger/internal/handler"
	ctxmeta "github.com/aaffriya/logger/pkg/context"
)

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
	handler := customhandler.NewHandler(loggerConfig, nil, file)
	logger := slog.New(handler)
	slog.SetDefault(logger)

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
			"code":        500,
			"message":     "Permission denied",
			"retry_count": 3,
		},
	)

	// Test complex data structures
	slog.InfoContext(ctx, "Complex data logging",
		"user_data", map[string]interface{}{
			"id":          12345,
			"name":        "John Doe",
			"email":       "john@example.com",
			"preferences": []string{"email", "sms", "push"},
			"metadata": map[string]interface{}{
				"last_login":  "2024-01-15T10:30:00Z",
				"login_count": 42,
				"is_premium":  true,
			},
		},
		"request_info", map[string]interface{}{
			"method": "POST",
			"url":    "/api/v1/users",
			"headers": map[string]string{
				"Content-Type": "application/json",
				"User-Agent":   "TestClient/1.0",
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
