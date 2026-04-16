package domain

import (
	"context"
	"errors"
)

var ErrUserIDRequired = errors.New("user settings user id is required")

type Repository interface {
	GetByUserID(ctx context.Context, userID string) (*UserSettings, error)
	UpsertPreferences(ctx context.Context, userID string, preferences map[string]any) (*UserSettings, error)
}
