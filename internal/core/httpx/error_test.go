package httpx

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	coreerrs "github.com/e-scavo/scavo-exchange-backend/internal/core/errs"
)

type envelope struct {
	Error coreerrs.ResponseError `json:"error"`
}

func TestWriteAppError_UsesCanonicalEnvelope(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteAppError(rec, coreerrs.AuthUnauthorized())

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected status: got=%d want=%d", rec.Code, http.StatusUnauthorized)
	}

	var payload envelope
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if payload.Error.Code != "AUTH_UNAUTHORIZED" {
		t.Fatalf("unexpected code: %q", payload.Error.Code)
	}
	if payload.Error.Message != "authentication required" {
		t.Fatalf("unexpected message: %q", payload.Error.Message)
	}
}

func TestWriteAppError_NilFallsBackToInternal(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteAppError(rec, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected status: got=%d want=%d", rec.Code, http.StatusInternalServerError)
	}

	var payload envelope
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if payload.Error.Code != "INTERNAL_ERROR" {
		t.Fatalf("unexpected code: %q", payload.Error.Code)
	}
}

func TestTimeout_UsesCanonicalErrorEnvelope(t *testing.T) {
	h := Timeout(5 * time.Millisecond)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(20 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/timeout", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("unexpected status: got=%d want=%d body=%s", rec.Code, http.StatusServiceUnavailable, rec.Body.String())
	}

	var payload envelope
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode error: %v | body=%s", err, rec.Body.String())
	}
	if payload.Error.Code != "TIMEOUT" {
		t.Fatalf("unexpected code: %q", payload.Error.Code)
	}
	if payload.Error.Message != "request timed out" {
		t.Fatalf("unexpected message: %q", payload.Error.Message)
	}
}

func TestRequireAuth_MissingBearerToken_UsesCanonicalEnvelope(t *testing.T) {
	h := RequireAuth(nil, false)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected status: got=%d want=%d", rec.Code, http.StatusInternalServerError)
	}

	var payload envelope
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if payload.Error.Code != "AUTH_NOT_CONFIGURED" {
		t.Fatalf("unexpected code: %q", payload.Error.Code)
	}
}
