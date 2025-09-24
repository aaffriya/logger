package console

import (
	"encoding/json"
	"log/slog"
)

func (h *prettyHandler) Json(r slog.Record) error {
	logLine := h.buildLogFirstLine(r)

	logData := map[string]any{}
	var trace []string

	r.Attrs(func(a slog.Attr) bool {
		if a.Key == "service" || a.Key == "version" ||
			a.Key == "trace_id" || a.Key == "user_id" || a.Key == "action" {
			return true
		}

		switch a.Key {
		case "trace":
			if traceVal, ok := a.Value.Any().([]string); ok {
				trace = traceVal
			}
		default:
			logData[a.Key] = h.convertValueForJSON(a.Value.Any())
		}
		return true
	})

	if len(trace) > 0 {
		logLine += NewLine + h.buildTraceSection(trace)
	}

	logLineByte := []byte(logLine)

	if len(logData) > 0 {
		logLineByte = append(logLineByte, []byte(NewLine+Gray+"Data:"+Reset+NewLine)...)
		jsonBytes, _ := json.MarshalIndent(logData, "", "  ")

		coloredJSON := applyJSONSyntaxHighlighting(string(jsonBytes))
		logLineByte = append(logLineByte, []byte(coloredJSON)...)
	}

	logLineByte = append(logLineByte, byte('\n'))

	_, err := h.writer.Write(logLineByte)
	return err
}

// convertValueForJSON recursively converts values to be JSON-serializable
func (h *prettyHandler) convertValueForJSON(value any) any {
	switch v := value.(type) {
	case error:
		// Convert error to string
		return v.Error()
	case map[string]any:
		// Recursively convert map values
		result := make(map[string]any)
		for k, val := range v {
			result[k] = h.convertValueForJSON(val)
		}
		return result
	case []any:
		// Recursively convert slice values
		result := make([]any, len(v))
		for i, val := range v {
			result[i] = h.convertValueForJSON(val)
		}
		return result
	default:
		// Return value as-is for other types
		return value
	}
}
