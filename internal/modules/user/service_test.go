package user

import (
	"context"
	"testing"
)

type stubRepo struct {
	lastEmail string
	result    *User
	err       error
}

func (s *stubRepo) UpsertDevUser(ctx context.Context, email string) (*User, error) {
	s.lastEmail = email
	return s.result, s.err
}

func TestResolveOrCreateDevUser_FallbackWithoutRepo(t *testing.T) {
	svc := NewService(nil)

	u, err := svc.ResolveOrCreateDevUser(context.Background(), "  TEST.User+dev@example.com ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if u == nil {
		t.Fatal("expected user, got nil")
	}

	if u.Email != "test.user+dev@example.com" {
		t.Fatalf("unexpected normalized email: %q", u.Email)
	}

	if u.ID != "u_test_user_dev_example_com" {
		t.Fatalf("unexpected user id: %q", u.ID)
	}

	if u.LastLoginAt == nil {
		t.Fatal("expected LastLoginAt to be set")
	}
}

func TestResolveOrCreateDevUser_UsesRepository(t *testing.T) {
	repo := &stubRepo{
		result: &User{
			ID:    "u_persisted",
			Email: "dev@example.com",
		},
	}
	svc := NewService(repo)

	u, err := svc.ResolveOrCreateDevUser(context.Background(), " Dev@Example.com ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.lastEmail != "dev@example.com" {
		t.Fatalf("repository received unexpected email: %q", repo.lastEmail)
	}

	if u == nil || u.ID != "u_persisted" {
		t.Fatalf("unexpected user result: %#v", u)
	}
}

func TestResolveOrCreateDevUser_EmptyEmail(t *testing.T) {
	svc := NewService(nil)

	u, err := svc.ResolveOrCreateDevUser(context.Background(), "   ")
	if err == nil {
		t.Fatal("expected error for empty email")
	}
	if u != nil {
		t.Fatalf("expected nil user, got %#v", u)
	}
}
