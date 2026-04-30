package httpx

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/google/uuid"

	coreerrs "github.com/e-scavo/scavo-exchange-backend/internal/core/errs"

	"github.com/e-scavo/scavo-exchange-backend/internal/core/logger"
)

type ctxKey string

const requestIDKey ctxKey = "request_id"

// RequestIDFromContext returns the correlation identifier attached to ctx by the
// RequestID middleware. An empty string means the request correlation middleware
// did not run yet or the context does not carry a valid request identifier.
func RequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	rid, _ := ctx.Value(requestIDKey).(string)
	return rid
}

func RequestID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rid := r.Header.Get("X-Request-Id")
			if rid == "" {
				rid = uuid.NewString()
			}
			w.Header().Set("X-Request-Id", rid)
			ctx := context.WithValue(r.Context(), requestIDKey, rid)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AccessLog(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := &wrapWriter{ResponseWriter: w, status: 200}
			next.ServeHTTP(ww, r)
			dur := time.Since(start)

			rid := RequestIDFromContext(r.Context())
			log.Info("http_request",
				"rid", rid,
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.status,
				"bytes", ww.bytes,
				"dur_ms", dur.Milliseconds(),
				"remote", r.RemoteAddr,
			)
		})
	}
}

type wrapWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *wrapWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *wrapWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}

func Recoverer(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					rid := RequestIDFromContext(r.Context())
					log.Error("panic",
						"rid", rid,
						"recover", rec,
						"stack", string(debug.Stack()),
					)
					WriteAppError(w, coreerrs.InternalError(nil))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func Timeout(d time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, d, string(mustJSON(coreerrs.NewErrorEnvelope(coreerrs.Timeout().ToResponseError()))))
	}
}

func mustJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		return []byte(`{"error":{"code":"INTERNAL_ERROR","message":"internal server error"}}`)
	}
	return b
}

func WriteJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func (w *wrapWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("hijacker not supported")
	}
	return h.Hijack()
}

func (w *wrapWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}
