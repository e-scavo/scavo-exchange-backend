// Package mappers centralizes explicit model transformations for the
// usersettings module.
package mappers

import (
	usersettingsdomain "github.com/e-scavo/scavo-exchange-backend/internal/modules/usersettings/domain"
	usersettingsreadmodels "github.com/e-scavo/scavo-exchange-backend/internal/modules/usersettings/readmodels"
)

// UserSettingsToReadModel maps the canonical user settings domain model into
// its output-only settings read projection.
func UserSettingsToReadModel(settings *usersettingsdomain.UserSettings) usersettingsreadmodels.UserSettingsReadModel {
	if settings == nil {
		return usersettingsreadmodels.UserSettingsReadModel{
			Version:     1,
			Preferences: map[string]any{},
		}
	}

	preferences := settings.Preferences
	if preferences == nil {
		preferences = map[string]any{}
	}

	model := usersettingsreadmodels.UserSettingsReadModel{
		UserID:      settings.UserID,
		Version:     1,
		Preferences: preferences,
	}

	if !settings.CreatedAt.IsZero() {
		createdAt := settings.CreatedAt.UTC()
		model.CreatedAt = &createdAt
	}
	if !settings.UpdatedAt.IsZero() {
		updatedAt := settings.UpdatedAt.UTC()
		model.UpdatedAt = &updatedAt
	}

	return model
}
