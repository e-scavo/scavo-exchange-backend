package authorization

import (
	"testing"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
)

func TestSubjectFromClaims_DefaultsToUserRole(t *testing.T) {
	subject, ok := SubjectFromClaims(&coreauth.Claims{UserID: "user-1"})
	if !ok {
		t.Fatalf("expected subject to resolve")
	}
	if subject.UserID != "user-1" {
		t.Fatalf("unexpected user id: %q", subject.UserID)
	}
	if len(subject.Roles) != 1 || subject.Roles[0] != RoleUser {
		t.Fatalf("unexpected roles: %#v", subject.Roles)
	}
}

func TestSubjectFromClaims_RejectsMissingUserID(t *testing.T) {
	if _, ok := SubjectFromClaims(&coreauth.Claims{}); ok {
		t.Fatalf("expected missing user id to be rejected")
	}
}
