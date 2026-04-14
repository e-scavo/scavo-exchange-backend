package httpx

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	coreauthorization "github.com/e-scavo/scavo-exchange-backend/internal/core/authorization"
)

func TestRequirePermission_AllowsAuthorizedSubject(t *testing.T) {
	h := RequirePermission(coreauthorization.ActionRead, coreauthorization.ResourceUser)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req = req.WithContext(coreauthorization.WithSubject(req.Context(), coreauthorization.AuthorizationSubject{
		UserID: "user-1",
		Roles:  []coreauthorization.Role{coreauthorization.RoleUser},
	}))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: got=%d want=%d body=%s", rec.Code, http.StatusNoContent, rec.Body.String())
	}
}

func TestRequirePermission_DeniesUnauthorizedSubject(t *testing.T) {
	h := RequirePermission(coreauthorization.ActionUpdate, coreauthorization.ResourceUser)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodPatch, "/auth/me", nil)
	req = req.WithContext(coreauthorization.WithSubject(req.Context(), coreauthorization.AuthorizationSubject{
		UserID: "user-1",
		Roles:  []coreauthorization.Role{coreauthorization.RoleUser},
	}))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("unexpected status: got=%d want=%d body=%s", rec.Code, http.StatusForbidden, rec.Body.String())
	}

	var payload struct {
		Error struct {
			Code    string         `json:"code"`
			Message string         `json:"message"`
			Details map[string]any `json:"details"`
		} `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v body=%s", err, rec.Body.String())
	}
	if payload.Error.Code != "AUTH_FORBIDDEN" {
		t.Fatalf("unexpected error code: got=%q want=%q", payload.Error.Code, "AUTH_FORBIDDEN")
	}
	if payload.Error.Details["permission"] != string(coreauthorization.PermissionUserUpdate) {
		t.Fatalf("unexpected permission detail: %#v", payload.Error.Details)
	}
}

func TestRequirePermission_DeniesMissingSubject(t *testing.T) {
	h := RequirePermission(coreauthorization.ActionRead, coreauthorization.ResourceUser)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req = req.WithContext(context.Background())
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("unexpected status: got=%d want=%d body=%s", rec.Code, http.StatusForbidden, rec.Body.String())
	}
}
