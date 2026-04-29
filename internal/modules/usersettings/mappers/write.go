package mappers

import (
	"encoding/json"

	usersettingsdomain "github.com/e-scavo/scavo-exchange-backend/internal/modules/usersettings/domain"
	usersettingswritemodels "github.com/e-scavo/scavo-exchange-backend/internal/modules/usersettings/writemodels"
)

// UserSettingsUpdateWriteToDomainInput maps user settings update write intent
// into canonical domain input. JSON decoding remains explicit to preserve
// current invalid-payload behavior.
func UserSettingsUpdateWriteToDomainInput(model usersettingswritemodels.UserSettingsUpdateWriteModel) (usersettingsdomain.UserSettingsUpdateInput, error) {
	input := usersettingsdomain.UserSettingsUpdateInput{}
	if len(model.Preferences) == 0 {
		return input, nil
	}
	if err := json.Unmarshal(model.Preferences, &input.Preferences); err != nil {
		return usersettingsdomain.UserSettingsUpdateInput{}, err
	}
	return input, nil
}
