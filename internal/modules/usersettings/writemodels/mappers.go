package writemodels

import (
	"encoding/json"

	usersettingsdomain "github.com/e-scavo/scavo-exchange-backend/internal/modules/usersettings/domain"
)

// ToDomainInput maps user settings update write intent into canonical domain
// input. JSON decoding is explicit so callers can preserve the current invalid
// payload behavior during handler alignment.
func (m UserSettingsUpdateWriteModel) ToDomainInput() (usersettingsdomain.UserSettingsUpdateInput, error) {
	input := usersettingsdomain.UserSettingsUpdateInput{}
	if len(m.Preferences) == 0 {
		return input, nil
	}
	if err := json.Unmarshal(m.Preferences, &input.Preferences); err != nil {
		return usersettingsdomain.UserSettingsUpdateInput{}, err
	}
	return input, nil
}
