package httpx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestID_UsesIncomingHeaderAndPropagatesContext(t *testing.T) {
	const incomingRequestID = "external-request-123"

	var got string
	h := RequestID()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = RequestIDFromContext(r.Context())
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("X-Request-Id", incomingRequestID)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if got != incomingRequestID {
		t.Fatalf("unexpected request id in context: got=%q want=%q", got, incomingRequestID)
	}
	if rec.Header().Get("X-Request-Id") != incomingRequestID {
		t.Fatalf("unexpected response request id header: got=%q want=%q", rec.Header().Get("X-Request-Id"), incomingRequestID)
	}
	if rec.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: got=%d want=%d", rec.Code, http.StatusNoContent)
	}
}

func TestRequestID_GeneratesRequestIDWhenHeaderIsMissing(t *testing.T) {
	var got string
	h := RequestID()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = RequestIDFromContext(r.Context())
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if got == "" {
		t.Fatalf("expected generated request id in context")
	}
	if rec.Header().Get("X-Request-Id") != got {
		t.Fatalf("unexpected response request id header: got=%q want=%q", rec.Header().Get("X-Request-Id"), got)
	}
	if rec.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: got=%d want=%d", rec.Code, http.StatusNoContent)
	}
}

func TestRequestIDFromContext_ReturnsEmptyForMissingOrNilContext(t *testing.T) {
	if got := RequestIDFromContext(context.Background()); got != "" {
		t.Fatalf("unexpected request id from empty context: got=%q want empty", got)
	}
	if got := RequestIDFromContext(nil); got != "" {
		t.Fatalf("unexpected request id from nil context: got=%q want empty", got)
	}
}
