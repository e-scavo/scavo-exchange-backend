package readmodels

import "time"

// UserSettingsReadModel is the explicit output projection for user settings.
// It is separate from the internal domain model and from future write/update
// models.
type UserSettingsReadModel struct {
	UserID      string         `json:"user_id"`
	Version     int            `json:"version"`
	Preferences map[string]any `json:"preferences"`
	CreatedAt   *time.Time     `json:"created_at,omitempty"`
	UpdatedAt   *time.Time     `json:"updated_at,omitempty"`
}
