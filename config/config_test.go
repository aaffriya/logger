package config

import (
	"encoding/json"
	"testing"
)

func TestLoggerConfigJSONSerialization(t *testing.T) {
	config := &LoggerConfig{
		Stack: StackConfig{
			Enabled: true,
			Skip:    5,
			Depth: StackDepths{
				Error: 10,
				Debug: 3,
				Info:  5,
				Warn:  7,
			},
		},
		DefaultFields: DefaultFieldInfo{
			Service: "TestService",
			Version: "v1.0.0",
		},
		Pretty: PrettyConfig{
			IncludeTimestamp: true,
			IsJsonOutput:     false,
		},
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("Failed to marshal config to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaledConfig LoggerConfig
	err = json.Unmarshal(jsonData, &unmarshaledConfig)
	if err != nil {
		t.Fatalf("Failed to unmarshal config from JSON: %v", err)
	}

	// Verify values
	if unmarshaledConfig.Stack.Enabled != config.Stack.Enabled {
		t.Errorf("Expected Stack.Enabled %v, got %v", config.Stack.Enabled, unmarshaledConfig.Stack.Enabled)
	}
	if unmarshaledConfig.Stack.Skip != config.Stack.Skip {
		t.Errorf("Expected Stack.Skip %d, got %d", config.Stack.Skip, unmarshaledConfig.Stack.Skip)
	}
	if unmarshaledConfig.Stack.Depth.Error != config.Stack.Depth.Error {
		t.Errorf("Expected Stack.Depth.Error %d, got %d", config.Stack.Depth.Error, unmarshaledConfig.Stack.Depth.Error)
	}
	if unmarshaledConfig.DefaultFields.Service != config.DefaultFields.Service {
		t.Errorf("Expected DefaultFields.Service %s, got %s", config.DefaultFields.Service, unmarshaledConfig.DefaultFields.Service)
	}
	if unmarshaledConfig.Pretty.IncludeTimestamp != config.Pretty.IncludeTimestamp {
		t.Errorf("Expected Pretty.IncludeTimestamp %v, got %v", config.Pretty.IncludeTimestamp, unmarshaledConfig.Pretty.IncludeTimestamp)
	}
}

func TestLoggerConfigJSONRoundTrip(t *testing.T) {
	config := &LoggerConfig{
		Stack: StackConfig{
			Enabled: true,
			Skip:    3,
			Depth: StackDepths{
				Error: 8,
				Debug: 2,
				Info:  4,
				Warn:  6,
			},
		},
		DefaultFields: DefaultFieldInfo{
			Service: "RoundTripTestService",
			Version: "v2.0.0",
		},
		Pretty: PrettyConfig{
			IncludeTimestamp: false,
			IsJsonOutput:     true,
		},
	}

	// Test JSON round trip
	jsonData, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("Failed to marshal config to JSON: %v", err)
	}

	var unmarshaledConfig LoggerConfig
	err = json.Unmarshal(jsonData, &unmarshaledConfig)
	if err != nil {
		t.Fatalf("Failed to unmarshal config from JSON: %v", err)
	}

	// Verify all values match exactly
	if unmarshaledConfig.Stack.Enabled != config.Stack.Enabled {
		t.Errorf("Expected Stack.Enabled %v, got %v", config.Stack.Enabled, unmarshaledConfig.Stack.Enabled)
	}
	if unmarshaledConfig.Stack.Skip != config.Stack.Skip {
		t.Errorf("Expected Stack.Skip %d, got %d", config.Stack.Skip, unmarshaledConfig.Stack.Skip)
	}
	if unmarshaledConfig.Stack.Depth.Debug != config.Stack.Depth.Debug {
		t.Errorf("Expected Stack.Depth.Debug %d, got %d", config.Stack.Depth.Debug, unmarshaledConfig.Stack.Depth.Debug)
	}
	if unmarshaledConfig.DefaultFields.Version != config.DefaultFields.Version {
		t.Errorf("Expected DefaultFields.Version %s, got %s", config.DefaultFields.Version, unmarshaledConfig.DefaultFields.Version)
	}
	if unmarshaledConfig.Pretty.IsJsonOutput != config.Pretty.IsJsonOutput {
		t.Errorf("Expected Pretty.IsJsonOutput %v, got %v", config.Pretty.IsJsonOutput, unmarshaledConfig.Pretty.IsJsonOutput)
	}
}

func TestStackDepthsValidation(t *testing.T) {
	depths := StackDepths{
		Error: 10,
		Debug: 3,
		Info:  5,
		Warn:  7,
	}

	if depths.Error <= 0 {
		t.Error("Expected Error depth to be positive")
	}
	if depths.Debug <= 0 {
		t.Error("Expected Debug depth to be positive")
	}
	if depths.Info <= 0 {
		t.Error("Expected Info depth to be positive")
	}
	if depths.Warn <= 0 {
		t.Error("Expected Warn depth to be positive")
	}
}

func TestDefaultFieldInfoValidation(t *testing.T) {
	fields := DefaultFieldInfo{
		Service: "TestService",
		Version: "v1.0.0",
	}

	if fields.Service == "" {
		t.Error("Expected Service to be non-empty")
	}
	if fields.Version == "" {
		t.Error("Expected Version to be non-empty")
	}
}

func TestPrettyConfigCombinations(t *testing.T) {
	testCases := []struct {
		name             string
		includeTimestamp bool
		isJsonOutput     bool
		description      string
	}{
		{
			name:             "Console with timestamp",
			includeTimestamp: true,
			isJsonOutput:     false,
			description:      "Pretty console output with timestamps",
		},
		{
			name:             "Console without timestamp",
			includeTimestamp: false,
			isJsonOutput:     false,
			description:      "Pretty console output without timestamps",
		},
		{
			name:             "JSON with timestamp",
			includeTimestamp: true,
			isJsonOutput:     true,
			description:      "JSON output with timestamps",
		},
		{
			name:             "JSON without timestamp",
			includeTimestamp: false,
			isJsonOutput:     true,
			description:      "JSON output without timestamps",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := PrettyConfig{
				IncludeTimestamp: tc.includeTimestamp,
				IsJsonOutput:     tc.isJsonOutput,
			}

			// Test that configuration is valid
			if config.IncludeTimestamp != tc.includeTimestamp {
				t.Errorf("Expected IncludeTimestamp %v, got %v", tc.includeTimestamp, config.IncludeTimestamp)
			}
			if config.IsJsonOutput != tc.isJsonOutput {
				t.Errorf("Expected IsJsonOutput %v, got %v", tc.isJsonOutput, config.IsJsonOutput)
			}
		})
	}
}

func TestCompleteLoggerConfig(t *testing.T) {
	config := &LoggerConfig{
		Stack: StackConfig{
			Enabled: true,
			Skip:    5,
			Depth: StackDepths{
				Error: 10,
				Debug: 3,
				Info:  5,
				Warn:  7,
			},
		},
		DefaultFields: DefaultFieldInfo{
			Service: "CompleteTestService",
			Version: "v1.2.3",
		},
		Pretty: PrettyConfig{
			IncludeTimestamp: true,
			IsJsonOutput:     false,
		},
	}

	// Test that all fields are properly set
	if !config.Stack.Enabled {
		t.Error("Expected Stack to be enabled")
	}
	if config.Stack.Skip != 5 {
		t.Errorf("Expected Stack.Skip to be 5, got %d", config.Stack.Skip)
	}
	if config.DefaultFields.Service != "CompleteTestService" {
		t.Errorf("Expected Service to be 'CompleteTestService', got '%s'", config.DefaultFields.Service)
	}
	if config.DefaultFields.Version != "v1.2.3" {
		t.Errorf("Expected Version to be 'v1.2.3', got '%s'", config.DefaultFields.Version)
	}
	if !config.Pretty.IncludeTimestamp {
		t.Error("Expected IncludeTimestamp to be true")
	}
	if config.Pretty.IsJsonOutput {
		t.Error("Expected IsJsonOutput to be false")
	}
}