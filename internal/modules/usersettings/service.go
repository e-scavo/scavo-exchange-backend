package usersettings

import (
	"context"
	"errors"
	"strings"
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
