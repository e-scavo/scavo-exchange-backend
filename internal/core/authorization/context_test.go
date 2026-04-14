package authorization

import (
	"context"
	"testing"
)

func TestWithSubject_StoresNormalizedSubject(t *testing.T) {
	ctx := WithSubject(context.Background(), AuthorizationSubject{
		UserID: "  user-1  ",
		Roles:  []Role{" USER ", RoleUser, RoleAdmin},
	})

	subject, ok := SubjectFromContext(ctx)
	if !ok || subject == nil {
		t.Fatalf("expected subject in context")
	}
	if subject.UserID != "user-1" {
		t.Fatalf("unexpected user id: %q", subject.UserID)
	}
	if len(subject.Roles) != 2 {
		t.Fatalf("unexpected roles length: got=%d want=%d", len(subject.Roles), 2)
	}
}

func TestWithSubject_EmptySubjectLeavesContextUnchanged(t *testing.T) {
	ctx := context.Background()
	ctx = WithSubject(ctx, AuthorizationSubject{})

	if _, ok := SubjectFromContext(ctx); ok {
		t.Fatalf("expected no subject in context")
	}
}
