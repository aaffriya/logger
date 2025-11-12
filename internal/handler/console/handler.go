package console

import (
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/aaffriya/logger/config"
)

type PrettyHandler interface {
	Handle(r slog.Record) error
}

type prettyHandler struct {
	writer io.Writer
	config *config.PrettyConfig
}

func NewPrettyHandler(w io.Writer, config *config.PrettyConfig) PrettyHandler {
	return &prettyHandler{
		writer: w,
		config: config,
	}
}

func (h *prettyHandler) buildLogFirstLine(r slog.Record) string {
	var builder strings.Builder

	if h.config.IncludeTimestamp {
		builder.WriteString(Gray)
		builder.WriteString(r.Time.Format("2006-01-02 15:04:05.000"))
		builder.WriteString(Reset)
		builder.WriteString(Space)
	}

	levelColor := LevelColor[int(r.Level.Level())]
	builder.WriteString("[")
	builder.WriteString(levelColor)
	builder.WriteString(Bold)
	builder.WriteString(r.Level.String())
	builder.WriteString(Reset)
	builder.WriteString("] ")

	builder.WriteString(Bold)
	builder.WriteString(White)
	builder.WriteString(r.Message)
	builder.WriteString(Reset)

	contextAttrs := make(map[string]any)
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == "trace_id" || a.Key == "user_id" || a.Key == "action" {
			contextAttrs[a.Key] = a.Value.Any()
		}
		return true
	})

	if len(contextAttrs) > 0 {
		builder.WriteString(Space)
		builder.WriteString(Gray)
		builder.WriteString("|")
		builder.WriteString(Reset)
		builder.WriteString(Space)

		order := []string{"trace_id", "user_id", "action"}
		parts := make([]string, 0, len(order))

		for _, key := range order {
			if value, exists := contextAttrs[key]; exists {
				color := getValueColor(key, value)
				parts = append(parts, color+fmt.Sprintf("%v", value)+Reset)
			}
		}

		builder.WriteString(joinStrings(parts, Space+Gray+"•"+Reset+Space))
	}

	return builder.String()
}

func (h *prettyHandler) buildTraceSection(trace []string) string {
	if len(trace) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.WriteString(Gray)
	builder.WriteString("Stack Trace:")
	builder.WriteString(Reset)
	builder.WriteString(NewLine)

	for i, s := range trace {
		traceColor := Gray + Italic
		if strings.Contains(s, ".go:") {
			traceColor = Cyan + Italic
		}
		builder.WriteString(fmt.Sprintf("  %s%d. %s%s%s", Gray, i+1, traceColor, s, Reset+NewLine))
	}
	return builder.String()
}

func (h *prettyHandler) Handle(r slog.Record) error {
	if h.config.IsJsonOutput {
		return h.Json(r)
	}
	return h.Text(r)
}
