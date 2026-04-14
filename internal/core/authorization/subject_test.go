package authorization

import "testing"

func TestAuthorizationSubjectNormalized(t *testing.T) {
	subject := AuthorizationSubject{
		UserID: "  user-123  ",
		Roles:  []Role{" USER ", "admin", "Admin", ""},
	}

	normalized := subject.Normalized()
	if normalized.UserID != "user-123" {
		t.Fatalf("expected trimmed user id, got %q", normalized.UserID)
	}

	expectedRoles := []Role{RoleUser, RoleAdmin}
	if len(normalized.Roles) != len(expectedRoles) {
		t.Fatalf("expected %d roles, got %d", len(expectedRoles), len(normalized.Roles))
	}

	for i := range expectedRoles {
		if normalized.Roles[i] != expectedRoles[i] {
			t.Fatalf("role %d: expected %q, got %q", i, expectedRoles[i], normalized.Roles[i])
		}
	}
}
