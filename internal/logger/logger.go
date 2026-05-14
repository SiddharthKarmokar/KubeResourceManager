package logger

import (
	"context"
	"log/slog"
	"os"
)

// Init initializes the global slog logger with JSON formatting.
func Init() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

// FromContext extracts a logger from the context, or returns the default logger.
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

// loggerKey is a custom type for context keys to avoid collisions.
type loggerKey struct{}

// WithLogger returns a new context with the provided logger.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}
