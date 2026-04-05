package usersettings

import "time"

type UserSettings struct {
	UserID      string
	Preferences map[string]any
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type View struct {
	UserID      string         `json:"user_id"`
	Version     int            `json:"version"`
	Preferences map[string]any `json:"preferences"`
}

func Default(userID string) *UserSettings {
	return &UserSettings{
		UserID:      userID,
		Preferences: map[string]any{},
	}
}

func ToView(settings *UserSettings) View {
	if settings == nil {
		return View{
			Version:     1,
			Preferences: map[string]any{},
		}
	}

	preferences := settings.Preferences
	if preferences == nil {
		preferences = map[string]any{}
	}

	return View{
		UserID:      settings.UserID,
		Version:     1,
		Preferences: preferences,
	}
}
