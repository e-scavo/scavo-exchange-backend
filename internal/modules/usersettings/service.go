package usersettings

import usersettingsapp "github.com/e-scavo/scavo-exchange-backend/internal/modules/usersettings/app"

var (
	ErrInvalidPreferences     = usersettingsapp.ErrInvalidPreferences
	ErrNullPreferenceValue    = usersettingsapp.ErrNullPreferenceValue
	ErrInvalidPreferenceValue = usersettingsapp.ErrInvalidPreferenceValue
	ErrIncompatiblePreference = usersettingsapp.ErrIncompatiblePreference
	ErrRepositoryRequired     = usersettingsapp.ErrRepositoryRequired
)

type Service = usersettingsapp.Service

func NewService(repo Repository) *Service {
	return usersettingsapp.NewService(repo)
}
