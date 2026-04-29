package domain

import (
	"context"

	usersettingsmod "github.com/e-scavo/scavo-exchange-backend/internal/modules/usersettings"
)

// =====================================================
// USER SETTINGS CONTRACT (CROSS-MODULE)
// =====================================================
//
// Defines minimal contract required by auth module.
//
// IMPORTANT:
// - Uses real usersettings types
// - Does NOT redefine models
//

type UserSettingsProvider interface {
	GetOrDefault(ctx context.Context, userID string) (*usersettingsmod.UserSettings, error)
	UpdatePreferences(ctx context.Context, userID string, preferences map[string]any) (*usersettingsmod.UserSettings, error)
}
