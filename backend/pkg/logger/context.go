package logger

import (
	"context"
	"log/slog"
)

type loggerCtxKey struct{}

// WithLogger returns a new context carrying the given logger.
func WithLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey{}, l)
}

// FromContext extracts the logger from ctx; falls back to slog.Default().
func FromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(loggerCtxKey{}).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}
