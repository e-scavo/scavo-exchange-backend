package usersettings

import usersettingsdomain "github.com/e-scavo/scavo-exchange-backend/internal/modules/usersettings/domain"

type UserSettings = usersettingsdomain.UserSettings
type View = usersettingsdomain.View

func Default(userID string) *UserSettings {
	return usersettingsdomain.Default(userID)
}

func ToView(settings *UserSettings) View {
	return usersettingsdomain.ToView(settings)
}
