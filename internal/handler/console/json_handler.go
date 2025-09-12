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
			logData[a.Key] = a.Value.Any()
		}
		return true
	})

	if len(trace) > 0 {
		logLine += NewLine + h.buildTraceSection(trace)
	}

	logLineByte := []byte(logLine)

	if len(logData) > 0 {
		logLineByte = append(logLineByte, []byte(Gray+"Data:"+Reset+NewLine)...)
		jsonBytes, _ := json.MarshalIndent(logData, "", "  ")

		coloredJSON := applyJSONSyntaxHighlighting(string(jsonBytes))
		logLineByte = append(logLineByte, []byte(coloredJSON)...)
	}

	logLineByte = append(logLineByte, byte('\n'))

	_, err := h.writer.Write(logLineByte)
	return err
}
