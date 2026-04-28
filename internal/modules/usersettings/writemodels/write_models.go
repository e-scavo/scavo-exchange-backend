// Package writemodels contains explicit input-only models for the usersettings
// module.
//
// Write models represent caller intent and must not be reused as responses or
// domain state.
package writemodels

import "encoding/json"

// UserSettingsUpdateWriteModel represents user settings update input intent.
// Raw JSON is preserved to keep compatibility with the current request surface
// until explicit Write → Domain mapping is introduced.
type UserSettingsUpdateWriteModel struct {
	Preferences json.RawMessage `json:"preferences"`
}
