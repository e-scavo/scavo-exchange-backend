package logger

import (
	"log/slog"
	"os"
)

const (
	RequestIDKey = "request_id"
	FlowEventKey = "flow_event"
)

type Logger struct{ *slog.Logger }

func New(env string) *Logger {
	level := slog.LevelInfo
	if env == "local" {
		level = slog.LevelDebug
	}

	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	return &Logger{Logger: slog.New(h)}
}

// AttrsWithRequestID returns the standard observability attributes for a
// correlated request. Empty request identifiers intentionally produce no
// attributes so callers can safely use it even before HTTP correlation exists.
func AttrsWithRequestID(requestID string) []any {
	if requestID == "" {
		return nil
	}
	return []any{RequestIDKey, requestID}
}

// WithRequestID returns a logger enriched with the standard request correlation
// attribute. Nil or empty inputs leave the logger unchanged.
func WithRequestID(log *Logger, requestID string) *Logger {
	if log == nil || requestID == "" {
		return log
	}
	return &Logger{Logger: log.With(RequestIDKey, requestID)}
}

// AttrsWithFlowEvent returns a standard observability attribute set for flow
// tracing. The event name is always included when present; request correlation
// is included only when available. Additional attributes are appended without
// interpretation so callers keep ownership of their domain-specific context.
func AttrsWithFlowEvent(requestID string, event string, attrs ...any) []any {
	base := make([]any, 0, 4+len(attrs))
	base = append(base, AttrsWithRequestID(requestID)...)
	if event != "" {
		base = append(base, FlowEventKey, event)
	}
	base = append(base, attrs...)
	return base
}
