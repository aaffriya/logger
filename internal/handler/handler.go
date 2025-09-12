package handler

import (
	"context"
	"github.com/aaffriya/logger/config"
	consolehandler "github.com/aaffriya/logger/internal/handler/console"
	filehandler "github.com/aaffriya/logger/internal/handler/file"
	"github.com/aaffriya/logger/internal/utils"
	ctxmeta "github.com/aaffriya/logger/pkg/context"
	"io"
	"log/slog"
	"os"
	"strings"
)

type LogHandler interface {
	Handle(r slog.Record) error
}

type Handler struct {
	config  *config.LoggerConfig
	writer  io.Writer
	opts    *slog.HandlerOptions
	handler LogHandler
	attrs   []slog.Attr
	groups  []string
	isFile  bool
}

func NewHandler(config *config.LoggerConfig, opts *slog.HandlerOptions, w io.Writer) slog.Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}

	// Determine if this is a file writer by checking if it's an *os.File
	isFile := false
	if f, ok := w.(*os.File); ok {
		fd := f.Fd()
		if fd != os.Stdout.Fd() && fd != os.Stderr.Fd() {
			isFile = true // Regular file, not stdout/stderr
		}
	}

	h := Handler{
		config: config,
		opts:   opts,
		writer: w,
		attrs:  make([]slog.Attr, 0),
		groups: make([]string, 0),
		isFile: isFile,
	}
	if h.isFile {
		if file, ok := h.writer.(*os.File); ok {
			h.handler = filehandler.NewFileHandler(h.writer, file)
		} else {
			h.handler = filehandler.NewFileHandler(h.writer, nil)
		}
	} else {
		h.handler = consolehandler.NewPrettyHandler(h.writer, &h.config.Pretty)
	}
	return &h
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	var recordAttrs []slog.Attr
	r.Attrs(func(a slog.Attr) bool {
		recordAttrs = append(recordAttrs, a)
		return true
	})

	allAttrs := h.prepareLogAttrs(ctx, r.Level, recordAttrs)

	newRecord := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
	newRecord.AddAttrs(allAttrs...)

	return h.handler.Handle(newRecord)
}

func (h *Handler) prepareLogAttrs(ctx context.Context, level slog.Level, recordAttrs []slog.Attr) []slog.Attr {
	attrs := make([]slog.Attr, 0, 10+len(recordAttrs)+len(h.attrs))

	if ctx != nil {
		contextData := ctxmeta.FromContext(ctx)
		if contextData.TraceID != "" {
			attrs = append(attrs, slog.String("trace_id", contextData.TraceID))
		}
		if contextData.UserID != "" {
			attrs = append(attrs, slog.String("user_id", contextData.UserID))
		}
		if contextData.Action != "" {
			attrs = append(attrs, slog.String("action", contextData.Action))
		}
	}

	if h.config.Stack.Enabled {
		var stackDepth int
		switch level {
		case slog.LevelError:
			stackDepth = h.config.Stack.Depth.Error
		case slog.LevelWarn:
			stackDepth = h.config.Stack.Depth.Warn
		case slog.LevelInfo:
			stackDepth = h.config.Stack.Depth.Info
		case slog.LevelDebug:
			stackDepth = h.config.Stack.Depth.Debug
		}

		if stackDepth > 0 {
			attrs = append(attrs, slog.Any("trace", utils.GetStackTrace(h.config.Stack.Skip, stackDepth)))
		}
	}

	attrs = append(attrs,
		slog.String("service", h.config.DefaultFields.Service),
		slog.String("version", h.config.DefaultFields.Version),
	)

	attrs = append(attrs, h.attrs...)

	if len(h.groups) > 0 {
		attrs = append(attrs, h.applyGroups(recordAttrs)...)
	} else {
		attrs = append(attrs, recordAttrs...)
	}

	return attrs
}

func (h *Handler) applyGroups(attrs []slog.Attr) []slog.Attr {
	if len(h.groups) == 0 {
		return attrs
	}

	groupPrefix := strings.Join(h.groups, ".")
	groupedAttrs := make([]slog.Attr, len(attrs))

	for i, attr := range attrs {
		groupedAttrs[i] = slog.Attr{
			Key:   groupPrefix + "." + attr.Key,
			Value: attr.Value,
		}
	}

	return groupedAttrs
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	copy(newAttrs[len(h.attrs):], attrs)

	handler := h.clone()
	handler.attrs = newAttrs
	return handler
}

func (h *Handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	newGroups := make([]string, len(h.groups)+1)
	copy(newGroups, h.groups)
	newGroups[len(h.groups)] = name

	handler := h.clone()
	handler.groups = newGroups

	return handler
}

func (h *Handler) clone() *Handler {
	return &Handler{
		writer:  h.writer,
		opts:    h.opts,
		config:  h.config,
		handler: h.handler,
		attrs:   append([]slog.Attr(nil), h.attrs...),
		groups:  append([]string(nil), h.groups...),
		isFile:  h.isFile,
	}
}

func (h *Handler) Close() error {
	if h.isFile {
		if file, ok := h.writer.(*os.File); ok {
			return file.Close()
		}
	}
	return nil
}
