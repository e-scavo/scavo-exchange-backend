package usersettings

import (
	"context"
	"errors"
	"testing"
)

type stubRepository struct {
	lastUserID       string
	lastUpsertUserID string
	lastUpsertPrefs  map[string]any
	result           *UserSettings
	upsertResult     *UserSettings
	err              error
	upsertErr        error
}

func (s *stubRepository) GetByUserID(ctx context.Context, userID string) (*UserSettings, error) {
	s.lastUserID = userID
	return s.result, s.err
}

func (s *stubRepository) UpsertPreferences(ctx context.Context, userID string, preferences map[string]any) (*UserSettings, error) {
	s.lastUpsertUserID = userID
	s.lastUpsertPrefs = preferences
	return s.upsertResult, s.upsertErr
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

func TestService_UpdatePreferences_MergesPersistedAndPatch(t *testing.T) {
	repo := &stubRepository{
		result: &UserSettings{
			UserID: "u_test",
			Preferences: map[string]any{
				"compact_mode": true,
				"theme":        "light",
			},
		},
		upsertResult: &UserSettings{
			UserID: "u_test",
			Preferences: map[string]any{
				"compact_mode": true,
				"theme":        "dark",
			},
		},
	}

	svc := NewService(repo)

	settings, err := svc.UpdatePreferences(context.Background(), "  u_test  ", map[string]any{"theme": "dark"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.lastUserID != "u_test" {
		t.Fatalf("unexpected get user id: %q", repo.lastUserID)
	}
	if repo.lastUpsertUserID != "u_test" {
		t.Fatalf("unexpected upsert user id: %q", repo.lastUpsertUserID)
	}
	if repo.lastUpsertPrefs["compact_mode"] != true || repo.lastUpsertPrefs["theme"] != "dark" {
		t.Fatalf("unexpected merged preferences: %#v", repo.lastUpsertPrefs)
	}
	if settings == nil || settings.Preferences["theme"] != "dark" {
		t.Fatalf("unexpected settings: %#v", settings)
	}
}

func TestService_UpdatePreferences_UsesDefaultsWhenNoRowExists(t *testing.T) {
	repo := &stubRepository{
		upsertResult: &UserSettings{
			UserID: "u_test",
			Preferences: map[string]any{
				"compact_mode": true,
			},
		},
	}

	svc := NewService(repo)

	settings, err := svc.UpdatePreferences(context.Background(), "u_test", map[string]any{"compact_mode": true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.lastUpsertPrefs["compact_mode"] != true {
		t.Fatalf("unexpected preferences: %#v", repo.lastUpsertPrefs)
	}
	if settings == nil || settings.Preferences["compact_mode"] != true {
		t.Fatalf("unexpected settings: %#v", settings)
	}
}

func TestService_UpdatePreferences_RejectsNilPatch(t *testing.T) {
	svc := NewService(&stubRepository{})

	settings, err := svc.UpdatePreferences(context.Background(), "u_test", nil)
	if err == nil {
		t.Fatal("expected error for nil patch")
	}
	if !errors.Is(err, ErrInvalidPreferences) {
		t.Fatalf("unexpected error: %v", err)
	}
	if settings != nil {
		t.Fatalf("expected nil settings, got %#v", settings)
	}
}

func TestService_UpdatePreferences_RejectsNullValue(t *testing.T) {
	svc := NewService(&stubRepository{})

	settings, err := svc.UpdatePreferences(context.Background(), "u_test", map[string]any{"theme": nil})
	if err == nil {
		t.Fatal("expected error for null value")
	}
	if !errors.Is(err, ErrNullPreferenceValue) {
		t.Fatalf("unexpected error: %v", err)
	}
	if settings != nil {
		t.Fatalf("expected nil settings, got %#v", settings)
	}
}

func TestService_UpdatePreferences_RejectsEmptyPreferenceKey(t *testing.T) {
	svc := NewService(&stubRepository{})

	settings, err := svc.UpdatePreferences(context.Background(), "u_test", map[string]any{"   ": true})
	if err == nil {
		t.Fatal("expected error for empty preference key")
	}
	if !errors.Is(err, ErrInvalidPreferences) {
		t.Fatalf("unexpected error: %v", err)
	}
	if settings != nil {
		t.Fatalf("expected nil settings, got %#v", settings)
	}
}

func TestService_UpdatePreferences_PropagatesRepositoryError(t *testing.T) {
	repo := &stubRepository{upsertErr: errors.New("boom")}
	svc := NewService(repo)

	settings, err := svc.UpdatePreferences(context.Background(), "u_test", map[string]any{"theme": "dark"})
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
