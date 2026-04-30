package httpx

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFrozenCoreStatusRoutesExposeJSONContract(t *testing.T) {
	router := newVersioningTestRouter(t)

	cases := []struct {
		path         string
		wantStatus   int
		requiredKeys []string
	}{
		{path: "/health", wantStatus: http.StatusOK, requiredKeys: []string{"ok", "status"}},
		{path: "/readiness", wantStatus: http.StatusOK, requiredKeys: []string{"ok", "status"}},
		{path: "/version", wantStatus: http.StatusOK, requiredKeys: []string{"version", "commit", "env"}},
		{path: "/diagnostics", wantStatus: http.StatusOK, requiredKeys: []string{"ok", "service", "env", "version", "commit", "observability", "time"}},
	}

	for _, tc := range cases {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != tc.wantStatus {
			t.Fatalf("unexpected status for %s: got=%d want=%d body=%s", tc.path, rec.Code, tc.wantStatus, rec.Body.String())
		}
		if got, want := rec.Header().Get("Content-Type"), "application/json; charset=utf-8"; got != want {
			t.Fatalf("unexpected content type for %s: got=%q want=%q", tc.path, got, want)
		}

		var payload map[string]any
		if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
			t.Fatalf("decode error for %s: %v body=%s", tc.path, err, rec.Body.String())
		}
		for _, key := range tc.requiredKeys {
			if _, ok := payload[key]; !ok {
				t.Fatalf("missing frozen response key %q for %s: %#v", key, tc.path, payload)
			}
		}
	}
}

func TestFrozenProtectedAuthRoutesExposeCanonicalErrorEnvelope(t *testing.T) {
	router := newVersioningTestRouter(t)

	for _, path := range []string{"/auth/me", "/api/v1/auth/me"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("unexpected status for %s: got=%d want=%d body=%s", path, rec.Code, http.StatusUnauthorized, rec.Body.String())
		}
		if got, want := rec.Header().Get("Content-Type"), "application/json; charset=utf-8"; got != want {
			t.Fatalf("unexpected content type for %s: got=%q want=%q", path, got, want)
		}

		var payload map[string]map[string]any
		if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
			t.Fatalf("decode error for %s: %v body=%s", path, err, rec.Body.String())
		}
		errorPayload, ok := payload["error"]
		if !ok {
			t.Fatalf("missing error envelope for %s: %#v", path, payload)
		}
		for _, key := range []string{"code", "message", "details"} {
			if _, ok := errorPayload[key]; !ok {
				t.Fatalf("missing frozen error key %q for %s: %#v", key, path, errorPayload)
			}
		}
		if errorPayload["code"] != "AUTH_MISSING_BEARER_TOKEN" {
			t.Fatalf("unexpected frozen error code for %s: got=%q", path, errorPayload["code"])
		}
		if _, ok := errorPayload["details"].(map[string]any); !ok {
			t.Fatalf("details must remain a JSON object for %s: %#v", path, errorPayload)
		}
	}
}
