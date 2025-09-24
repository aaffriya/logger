package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/aaffriya/logger/config"
	customhandler "github.com/aaffriya/logger/internal/handler"
	ctxmeta "github.com/aaffriya/logger/pkg/context"
)

func TestSetupConsolePrettyLogger(t *testing.T) {
	// Capture stdout
	var buf bytes.Buffer
	
	loggerConfig := &config.LoggerConfig{
		DefaultFields: config.DefaultFieldInfo{
			Version: "v1.0.0",
			Service: "TestService",
		},
		Stack: config.StackConfig{
			Enabled: false, // Disable for cleaner test output
		},
		Pretty: config.PrettyConfig{
			IncludeTimestamp: false,
			IsJsonOutput:     false,
		},
	}

	// Create a custom handler that writes to our buffer instead of stdout
	handler := customhandler.NewHandler(loggerConfig, nil, &buf)
	logger := slog.New(handler)
	
	// Test basic logging
	logger.Info("Test message", "key", "value")
	
	output := buf.String()
	if !strings.Contains(output, "Test message") {
		t.Errorf("Expected log output to contain 'Test message', got: %s", output)
	}
	if !strings.Contains(output, "key") || !strings.Contains(output, "value") {
		t.Errorf("Expected log output to contain key and value, got: %s", output)
	}
}

func TestSetupFileLogger(t *testing.T) {
	// Create temporary file
	tmpFile, err := os.CreateTemp("", "test_log_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	loggerConfig := &config.LoggerConfig{
		DefaultFields: config.DefaultFieldInfo{
			Version: "v1.0.0",
			Service: "TestService",
		},
		Stack: config.StackConfig{
			Enabled: false,
		},
		Pretty: config.PrettyConfig{
			IncludeTimestamp: true,
			IsJsonOutput:     true,
		},
	}

	handler := customhandler.NewHandler(loggerConfig, nil, tmpFile)
	logger := slog.New(handler)
	slog.SetDefault(logger)
	
	// Test logging
	slog.Info("Test file message", "test_key", "test_value")
	
	// Read file content
	tmpFile.Seek(0, 0)
	var logEntry map[string]interface{}
	decoder := json.NewDecoder(tmpFile)
	if err := decoder.Decode(&logEntry); err != nil {
		t.Fatalf("Failed to decode JSON log: %v", err)
	}
	
	// Verify JSON structure
	if logEntry["message"] != "Test file message" {
		t.Errorf("Expected message 'Test file message', got: %v", logEntry["message"])
	}
	if logEntry["test_key"] != "test_value" {
		t.Errorf("Expected test_key 'test_value', got: %v", logEntry["test_key"])
	}
	if logEntry["service"] != "TestService" {
		t.Errorf("Expected service 'TestService', got: %v", logEntry["service"])
	}
	if logEntry["version"] != "v1.0.0" {
		t.Errorf("Expected version 'v1.0.0', got: %v", logEntry["version"])
	}
}

func TestContextMetadataIntegration(t *testing.T) {
	var buf bytes.Buffer
	
	loggerConfig := &config.LoggerConfig{
		DefaultFields: config.DefaultFieldInfo{
			Version: "v1.0.0",
			Service: "TestService",
		},
		Stack: config.StackConfig{
			Enabled: false,
		},
		Pretty: config.PrettyConfig{
			IncludeTimestamp: false,
			IsJsonOutput:     false,
		},
	}

	handler := customhandler.NewHandler(loggerConfig, nil, &buf)
	logger := slog.New(handler)
	
	// Create context with metadata
	ctx := ctxmeta.WithAction(context.Background(), "TEST_ACTION")
	ctx = ctxmeta.WithTraceID(ctx, "trace-123")
	ctx = ctxmeta.WithUserID(ctx, "user-456")
	
	// Log with context
	logger.InfoContext(ctx, "Context test message")
	
	output := buf.String()
	if !strings.Contains(output, "TEST_ACTION") {
		t.Errorf("Expected output to contain 'TEST_ACTION', got: %s", output)
	}
	if !strings.Contains(output, "trace-123") {
		t.Errorf("Expected output to contain 'trace-123', got: %s", output)
	}
	if !strings.Contains(output, "user-456") {
		t.Errorf("Expected output to contain 'user-456', got: %s", output)
	}
}

func TestLogLevels(t *testing.T) {
	var buf bytes.Buffer
	
	loggerConfig := &config.LoggerConfig{
		DefaultFields: config.DefaultFieldInfo{
			Version: "v1.0.0",
			Service: "TestService",
		},
		Stack: config.StackConfig{
			Enabled: false,
		},
		Pretty: config.PrettyConfig{
			IncludeTimestamp: false,
			IsJsonOutput:     false,
		},
	}

	handler := customhandler.NewHandler(loggerConfig, nil, &buf)
	logger := slog.New(handler)
	
	// Test different log levels
	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warn message")
	logger.Error("Error message")
	
	output := buf.String()
	
	// Check that all levels are present (Debug might be filtered out depending on level)
	if !strings.Contains(output, "Info message") {
		t.Errorf("Expected output to contain 'Info message', got: %s", output)
	}
	if !strings.Contains(output, "Warn message") {
		t.Errorf("Expected output to contain 'Warn message', got: %s", output)
	}
	if !strings.Contains(output, "Error message") {
		t.Errorf("Expected output to contain 'Error message', got: %s", output)
	}
}

func TestStackTraceEnabled(t *testing.T) {
	var buf bytes.Buffer
	
	loggerConfig := &config.LoggerConfig{
		DefaultFields: config.DefaultFieldInfo{
			Version: "v1.0.0",
			Service: "TestService",
		},
		Stack: config.StackConfig{
			Enabled: true,
			Skip:    3,
			Depth: config.StackDepths{
				Error: 2,
				Info:  1,
				Debug: 1,
				Warn:  1,
			},
		},
		Pretty: config.PrettyConfig{
			IncludeTimestamp: false,
			IsJsonOutput:     false,
		},
	}

	handler := customhandler.NewHandler(loggerConfig, nil, &buf)
	logger := slog.New(handler)
	
	logger.Error("Error with stack trace")
	
	output := buf.String()
	if !strings.Contains(output, "Stack Trace:") {
		t.Errorf("Expected output to contain stack trace, got: %s", output)
	}
}

func TestErrorObjectInJson(t *testing.T) {
	jsonObj := map[string]any{
		"msg": "An error occurred",
		"error": errors.New("sample error"),
	}
	
	var buf bytes.Buffer
	loggerConfig := &config.LoggerConfig{
		DefaultFields: config.DefaultFieldInfo{
			Version: "v1.0.0",
			Service: "TestService",
		},
		Stack: config.StackConfig{
			Enabled: false,
		},
		Pretty: config.PrettyConfig{
			IncludeTimestamp: true,
			IsJsonOutput:     true,
		},
	}

	handler := customhandler.NewHandler(loggerConfig, nil, &buf)
	logger := slog.New(handler)
	
	logger.Info("Logging error object", "data", jsonObj)
	
	output := buf.String()
	
	// The output should be in hybrid format: formatted line + JSON data section
	if !strings.Contains(output, "Logging error object") {
		t.Errorf("Expected output to contain log message, got: %s", output)
	}
	if !strings.Contains(output, "Data:") {
		t.Errorf("Expected output to contain 'Data:' section, got: %s", output)
	}
	if !strings.Contains(output, "sample error") {
		t.Errorf("Expected output to contain 'sample error', got: %s", output)
	}
	if !strings.Contains(output, "An error occurred") {
		t.Errorf("Expected output to contain 'An error occurred', got: %s", output)
	}
}