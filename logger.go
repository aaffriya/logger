package logger

import (
	"log/slog"
	"os"

	"github.com/aaffriya/logger/config"
	customhandler "github.com/aaffriya/logger/internal/handler"
)

func SetupConsolePrettyLogger(config *config.LoggerConfig, opts *slog.HandlerOptions) {
	writer := os.Stdout
	handler := customhandler.NewHandler(config, opts, writer)
	logger := slog.New(handler)
	slog.SetDefault(logger)

}

func SetupFileLogger(config *config.LoggerConfig, opts *slog.HandlerOptions, file *os.File) {
	handler := customhandler.NewHandler(config, opts, file)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
