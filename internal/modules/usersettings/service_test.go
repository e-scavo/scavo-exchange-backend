package usersettings

import (
	"context"
	"encoding/json"
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

func TestService_GetOrDefault_NormalizesRepositoryPreferences(t *testing.T) {
	repo := &stubRepository{
		result: &UserSettings{
			UserID: "u_test",
			Preferences: map[string]any{
				" theme ": json.Number("2"),
				"nested": map[string]any{
					"enabled": true,
				},
			},
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

	if _, exists := settings.Preferences[" theme "]; exists {
		t.Fatalf("expected key to be normalized, got %#v", settings.Preferences)
	}

	if value, ok := settings.Preferences["theme"]; !ok || value != float64(2) {
		t.Fatalf("unexpected normalized preferences: %#v", settings.Preferences)
	}
}

func TestService_GetOrDefault_RejectsInvalidPersistedPreferences(t *testing.T) {
	repo := &stubRepository{
		result: &UserSettings{
			UserID: "u_test",
			Preferences: map[string]any{
				"broken": nil,
			},
		},
	}

	svc := NewService(repo)

	settings, err := svc.GetOrDefault(context.Background(), "u_test")
	if err == nil {
		t.Fatal("expected error for invalid persisted preferences")
	}

	if !errors.Is(err, ErrNullPreferenceValue) {
		t.Fatalf("unexpected error: %v", err)
	}

	if settings != nil {
		t.Fatalf("expected nil settings, got %#v", settings)
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

func TestService_UpdatePreferences_NormalizesPatchBeforePersisting(t *testing.T) {
	repo := &stubRepository{
		upsertResult: &UserSettings{
			UserID: "u_test",
			Preferences: map[string]any{
				"theme": float64(1),
				"nested": map[string]any{
					"count": float64(3),
				},
			},
		},
	}

	svc := NewService(repo)

	settings, err := svc.UpdatePreferences(context.Background(), "u_test", map[string]any{
		" theme ": json.Number("1"),
		"nested": map[string]int{
			"count": 3,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, exists := repo.lastUpsertPrefs[" theme "]; exists {
		t.Fatalf("expected trimmed key, got %#v", repo.lastUpsertPrefs)
	}

	if repo.lastUpsertPrefs["theme"] != float64(1) {
		t.Fatalf("expected normalized number, got %#v", repo.lastUpsertPrefs["theme"])
	}

	nested, ok := repo.lastUpsertPrefs["nested"].(map[string]any)
	if !ok {
		t.Fatalf("expected nested map[string]any, got %#v", repo.lastUpsertPrefs["nested"])
	}

	if nested["count"] != float64(3) {
		t.Fatalf("expected normalized nested number, got %#v", nested["count"])
	}

	if settings == nil {
		t.Fatal("expected settings, got nil")
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

func TestService_UpdatePreferences_RejectsNestedNullValue(t *testing.T) {
	svc := NewService(&stubRepository{})

	settings, err := svc.UpdatePreferences(context.Background(), "u_test", map[string]any{
		"notifications": map[string]any{
			"email": nil,
		},
	})
	if err == nil {
		t.Fatal("expected error for nested null value")
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

func TestService_UpdatePreferences_RejectsInvalidPreferenceValue(t *testing.T) {
	svc := NewService(&stubRepository{})

	settings, err := svc.UpdatePreferences(context.Background(), "u_test", map[string]any{
		"theme": func() {},
	})
	if err == nil {
		t.Fatal("expected error for invalid preference value")
	}
	if !errors.Is(err, ErrInvalidPreferenceValue) {
		t.Fatalf("unexpected error: %v", err)
	}
	if settings != nil {
		t.Fatalf("expected nil settings, got %#v", settings)
	}
}

func TestService_UpdatePreferences_RejectsTopLevelShapeMismatch(t *testing.T) {
	repo := &stubRepository{
		result: &UserSettings{
			UserID: "u_test",
			Preferences: map[string]any{
				"notifications": map[string]any{
					"email": true,
				},
			},
		},
	}

	svc := NewService(repo)

	settings, err := svc.UpdatePreferences(context.Background(), "u_test", map[string]any{
		"notifications": true,
	})
	if err == nil {
		t.Fatal("expected error for incompatible top-level shape")
	}
	if !errors.Is(err, ErrIncompatiblePreference) {
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
