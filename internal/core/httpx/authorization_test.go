package httpx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
	coreauthorization "github.com/e-scavo/scavo-exchange-backend/internal/core/authorization"
)

func TestHydrateAuthorization_AttachesSubjectFromClaims(t *testing.T) {
	var got *coreauthorization.AuthorizationSubject

	h := HydrateAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ok bool
		got, ok = AuthorizationSubjectFromContext(r.Context())
		if !ok {
			t.Fatalf("expected authorization subject in context")
		}
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	ctx := context.WithValue(req.Context(), coreauth.ClaimsContextKey, &coreauth.Claims{UserID: "user-1"})
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: got=%d want=%d", rec.Code, http.StatusNoContent)
	}
	if got == nil || got.UserID != "user-1" {
		t.Fatalf("unexpected subject: %#v", got)
	}
	if len(got.Roles) != 1 || got.Roles[0] != coreauthorization.RoleUser {
		t.Fatalf("unexpected roles: %#v", got.Roles)
	}
}

func TestHydrateAuthorization_WithoutClaimsLeavesContextUntouched(t *testing.T) {
	h := HydrateAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := AuthorizationSubjectFromContext(r.Context()); ok {
			t.Fatalf("did not expect authorization subject")
		}
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: got=%d want=%d", rec.Code, http.StatusNoContent)
	}
}
