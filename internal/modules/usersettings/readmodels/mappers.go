package readmodels

import usersettingsdomain "github.com/e-scavo/scavo-exchange-backend/internal/modules/usersettings/domain"

// FromUserSettings maps the canonical user settings domain model into its
// output-only user settings read model projection.
func FromUserSettings(settings *usersettingsdomain.UserSettings) UserSettingsReadModel {
	if settings == nil {
		return UserSettingsReadModel{
			Version:     1,
			Preferences: map[string]any{},
		}
	}

	preferences := settings.Preferences
	if preferences == nil {
		preferences = map[string]any{}
	}

	model := UserSettingsReadModel{
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
