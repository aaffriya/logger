package console

import (
	"net"
	"strings"
)

const (
	Reset          = "\033[0m"
	Bold           = "\033[1m"
	Italic         = "\033[3m"
	UnderlineColor = "\033[4m"

	Gray             = "\033[90m"
	White            = "\033[97m"
	Cyan             = "\033[36m"
	Green            = "\033[32m"
	Yellow           = "\033[33m"
	Red              = "\033[31m"
	Magenta          = "\033[35m"
	Blue             = "\033[34m"
	BrightBlackColor = "\033[90m"

	NewLine = "\n"
	Tab     = "\t"
	Space   = " "
)

var LevelColor = map[int]string{
	-4: Cyan,   // DEBUG
	0:  Green,  // INFO
	4:  Yellow, // WARN
	8:  Red,    // ERROR
}

// Field-specific color mapping based on data type and context
var FieldColors = map[string]string{
	"trace_id":    Cyan + Bold,
	"span_id":     Cyan + Bold,
	"trace_flags": Cyan + Bold,
	"user_id":     Blue + Bold,
	"action":      Magenta + Bold,
	"service":     Green + Italic,
	"version":     Gray + Italic,
	"error":       Red + Bold,
	"database":    Blue,
	"timeout":     Yellow,
	"retry_count": Yellow,
	"ip_address":  Cyan,
	"user_agent":  Gray,
	"email":       Blue,
	"status":      Green,
	"method":      Magenta,
	"url":         UnderlineColor + Blue,
	"duration":    Yellow,
	"memory":      Cyan,
	"cpu":         Green,
	"disk":        Yellow,
}

// Value-based color logic
func getValueColor(key string, value any) string {
	// Check field-specific colors first
	if color, exists := FieldColors[key]; exists {
		return color
	}

	// Logic-based coloring by value type and content
	switch v := value.(type) {
	case string:
		return getStringColor(key, v)
	case int, int32, int64, float32, float64:
		return getNumericColor(key, v)
	case bool:
		return getBoolColor(v)
	case map[string]any:
		return Magenta // Objects/maps
	case []any, []string:
		return Cyan // Arrays
	default:
		return White
	}
}

// String value color logic
func getStringColor(key, value string) string {
	// Email patterns
	if key == "email" || (len(value) > 0 && strings.Contains(value, "@")) {
		return Blue + UnderlineColor
	}

	// URL patterns
	if key == "url" || strings.HasPrefix(value, "http") || strings.HasPrefix(value, "https") {
		return Blue + UnderlineColor
	}

	// IP address patterns
	if key == "ip" || key == "ip_address" || net.ParseIP(value) != nil {
		return Cyan + Bold
	}

	// Error/failure patterns
	if key == "error" || key == "failure" || strings.Contains(value, "error") || strings.Contains(value, "failed") {
		return Red + Bold
	}

	// Success patterns
	if strings.Contains(value, "success") || strings.Contains(value, "completed") || strings.Contains(value, "ok") {
		return Green + Bold
	}

	// Warning patterns
	if strings.Contains(value, "warning") || strings.Contains(value, "timeout") || strings.Contains(value, "retry") {
		return Yellow + Bold
	}

	// File paths
	if strings.Contains(value, "/") && (strings.Contains(value, ".go") || strings.Contains(value, ".js") || strings.Contains(value, ".py")) {
		return Gray + Italic
	}

	// Database related
	if key == "database" || key == "table" || key == "query" {
		return Blue
	}

	// Time/duration patterns
	if key == "duration" || key == "timeout" || strings.Contains(value, "ms") || strings.Contains(value, "seconds") {
		return Yellow
	}

	return White
}

// Numeric value color logic
func getNumericColor(key string, value any) string {
	// Convert to float64 for comparison
	var numVal float64
	switch v := value.(type) {
	case int:
		numVal = float64(v)
	case int32:
		numVal = float64(v)
	case int64:
		numVal = float64(v)
	case float32:
		numVal = float64(v)
	case float64:
		numVal = v
	default:
		return White
	}

	// Status codes
	if key == "status" || key == "code" || key == "status_code" {
		if numVal >= 200 && numVal < 300 {
			return Green + Bold // Success
		} else if numVal >= 400 && numVal < 500 {
			return Yellow + Bold // Client error
		} else if numVal >= 500 {
			return Red + Bold // Server error
		}
		return Cyan
	}

	// Performance metrics
	if key == "duration" || key == "latency" || key == "response_time" {
		if numVal < 100 {
			return Green // Fast
		} else if numVal < 1000 {
			return Yellow // Medium
		} else {
			return Red // Slow
		}
	}

	// Memory/size metrics
	if key == "memory" || key == "size" || key == "bytes" {
		if numVal < 1024*1024 { // < 1MB
			return Green
		} else if numVal < 1024*1024*100 { // < 100MB
			return Yellow
		} else {
			return Red
		}
	}

	// Retry/attempt counts
	if key == "retry" || key == "retry_count" || key == "attempts" {
		if numVal == 0 {
			return Green
		} else if numVal < 3 {
			return Yellow
		} else {
			return Red
		}
	}

	// Port numbers
	if key == "port" {
		return Cyan + Bold
	}

	// Default numeric color
	return Cyan
}

// Boolean value color logic
func getBoolColor(value bool) string {
	if value {
		return Green + Bold
	}
	return Red + Bold
}
