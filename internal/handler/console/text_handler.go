package console

import (
	"fmt"
	"log/slog"
	"strings"
)

func (h *prettyHandler) Text(r slog.Record) error {
	var builder strings.Builder

	// Build the first line
	builder.WriteString(h.buildLogFirstLine(r))

	attrs := make(map[string]any)
	var trace []string

	r.Attrs(func(a slog.Attr) bool {
		if a.Key == "service" || a.Key == "version" ||
			a.Key == "trace_id" || a.Key == "user_id" || a.Key == "action" {
			return true
		}

		if a.Key == "trace" {
			if traceVal, ok := a.Value.Any().([]string); ok {
				trace = traceVal
			}
			return true
		}

		attrs[a.Key] = a.Value.Any()
		return true
	})

	if len(trace) > 0 {
		builder.WriteString(NewLine)
		builder.WriteString(h.buildTraceSection(trace))
	}

	if len(attrs) > 0 {
		builder.WriteString(NewLine)
		for key, value := range attrs {
			keyColor := getValueColor(key, key)
			valueColor := getValueColor(key, value)

			builder.WriteString(fmt.Sprintf("  %s%s%s%s=%s%v%s%s",
				Gray, keyColor, key, Reset,
				valueColor, value, Reset, NewLine))
		}
	} else if len(trace) == 0 {
		builder.WriteString(NewLine)
	}

	_, err := h.writer.Write([]byte(builder.String()))
	return err
}

func joinStrings(strs []string, separator string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += separator + strs[i]
	}
	return result
}
