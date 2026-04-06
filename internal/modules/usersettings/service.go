package usersettings

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrInvalidPreferences  = errors.New("user settings preferences must be an object")
	ErrNullPreferenceValue = errors.New("user settings preference value cannot be null")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetOrDefault(ctx context.Context, userID string) (*UserSettings, error) {
	normalizedUserID := strings.TrimSpace(userID)
	if normalizedUserID == "" {
		return nil, ErrUserIDRequired
	}

	if s.repo == nil {
		return nil, errors.New("user settings repository is required")
	}

	settings, err := s.repo.GetByUserID(ctx, normalizedUserID)
	if err != nil {
		return nil, err
	}

	if settings == nil {
		return Default(normalizedUserID), nil
	}

	if settings.UserID == "" {
		settings.UserID = normalizedUserID
	}

	if settings.Preferences == nil {
		settings.Preferences = map[string]any{}
	}

	return settings, nil
}

func (s *Service) UpdatePreferences(ctx context.Context, userID string, patch map[string]any) (*UserSettings, error) {
	normalizedUserID := strings.TrimSpace(userID)
	if normalizedUserID == "" {
		return nil, ErrUserIDRequired
	}

	if s.repo == nil {
		return nil, errors.New("user settings repository is required")
	}

	if patch == nil {
		return nil, ErrInvalidPreferences
	}

	for key, value := range patch {
		if strings.TrimSpace(key) == "" {
			return nil, ErrInvalidPreferences
		}
		if value == nil {
			return nil, ErrNullPreferenceValue
		}
	}

	current, err := s.GetOrDefault(ctx, normalizedUserID)
	if err != nil {
		return nil, err
	}

	merged := make(map[string]any, len(current.Preferences)+len(patch))
	for key, value := range current.Preferences {
		merged[key] = value
	}
	for key, value := range patch {
		merged[key] = value
	}

	updated, err := s.repo.UpsertPreferences(ctx, normalizedUserID, merged)
	if err != nil {
		return nil, err
	}

	if updated == nil {
		updated = &UserSettings{
			UserID:      normalizedUserID,
			Preferences: merged,
		}
	}
	if updated.UserID == "" {
		updated.UserID = normalizedUserID
	}
	if updated.Preferences == nil {
		updated.Preferences = map[string]any{}
	}

	return updated, nil
}
