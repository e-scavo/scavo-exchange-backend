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

	"github.com/e-scavo/scavo-exchange-backend/internal/core/logger"
)

type ctxKey string

const requestIDKey ctxKey = "request_id"

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

			rid, _ := r.Context().Value(requestIDKey).(string)
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
					rid, _ := r.Context().Value(requestIDKey).(string)
					log.Error("panic",
						"rid", rid,
						"recover", rec,
						"stack", string(debug.Stack()),
					)
					WriteJSON(w, http.StatusInternalServerError, map[string]any{"error": "internal_server_error"})
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func Timeout(d time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, d, `{"error":"timeout"}`)
	}
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
