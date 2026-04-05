package usersettings

import (
	"context"
	"errors"
	"testing"
)

type stubRepository struct {
	lastUserID string
	result     *UserSettings
	err        error
}

func (s *stubRepository) GetByUserID(ctx context.Context, userID string) (*UserSettings, error) {
	s.lastUserID = userID
	return s.result, s.err
}

func TestService_GetOrDefault_UsesRepositoryWithNormalizedUserID(t *testing.T) {
	repo := &stubRepository{
		result: &UserSettings{
			UserID: "u_test",
			Preferences: map[string]any{
				"example": true,
			},
		},
	}

	svc := NewService(repo)

	settings, err := svc.GetOrDefault(context.Background(), "  u_test  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.lastUserID != "u_test" {
		t.Fatalf("repository received unexpected user id: %q", repo.lastUserID)
	}

	if settings == nil {
		t.Fatal("expected settings, got nil")
	}

	if settings.UserID != "u_test" {
		t.Fatalf("unexpected user id: %q", settings.UserID)
	}

	if settings.Preferences == nil {
		t.Fatal("expected preferences map")
	}

	if value, ok := settings.Preferences["example"]; !ok || value != true {
		t.Fatalf("unexpected preferences: %#v", settings.Preferences)
	}
}

func TestService_GetOrDefault_ReturnsDefaultWhenNotPersisted(t *testing.T) {
	repo := &stubRepository{}
	svc := NewService(repo)

	settings, err := svc.GetOrDefault(context.Background(), "u_default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if settings == nil {
		t.Fatal("expected settings, got nil")
	}

	if settings.UserID != "u_default" {
		t.Fatalf("unexpected user id: %q", settings.UserID)
	}

	if settings.Preferences == nil {
		t.Fatal("expected default preferences map")
	}

	if len(settings.Preferences) != 0 {
		t.Fatalf("expected empty preferences, got %#v", settings.Preferences)
	}
}

func TestService_GetOrDefault_InitializesNilPreferencesFromRepository(t *testing.T) {
	repo := &stubRepository{
		result: &UserSettings{
			UserID:      "u_test",
			Preferences: nil,
		},
	}

	svc := NewService(repo)

	settings, err := svc.GetOrDefault(context.Background(), "u_test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if settings == nil {
		t.Fatal("expected settings, got nil")
	}

	if settings.Preferences == nil {
		t.Fatal("expected preferences to be initialized")
	}

	if len(settings.Preferences) != 0 {
		t.Fatalf("expected empty preferences, got %#v", settings.Preferences)
	}
}

func TestService_GetOrDefault_RejectsEmptyUserID(t *testing.T) {
	svc := NewService(&stubRepository{})

	settings, err := svc.GetOrDefault(context.Background(), "   ")
	if err == nil {
		t.Fatal("expected error for empty user id")
	}

	if !errors.Is(err, ErrUserIDRequired) {
		t.Fatalf("unexpected error: %v", err)
	}

	if settings != nil {
		t.Fatalf("expected nil settings, got %#v", settings)
	}
}

func TestService_GetOrDefault_RejectsMissingRepository(t *testing.T) {
	svc := NewService(nil)

	settings, err := svc.GetOrDefault(context.Background(), "u_test")
	if err == nil {
		t.Fatal("expected error for missing repository")
	}

	if err.Error() != "user settings repository is required" {
		t.Fatalf("unexpected error: %v", err)
	}

	if settings != nil {
		t.Fatalf("expected nil settings, got %#v", settings)
	}
}

func TestService_GetOrDefault_PropagatesRepositoryError(t *testing.T) {
	repo := &stubRepository{
		err: errors.New("boom"),
	}

	svc := NewService(repo)

	settings, err := svc.GetOrDefault(context.Background(), "u_test")
	if err == nil {
		t.Fatal("expected repository error")
	}

	if err.Error() != "boom" {
		t.Fatalf("unexpected error: %v", err)
	}

	if settings != nil {
		t.Fatalf("expected nil settings, got %#v", settings)
	}
}
