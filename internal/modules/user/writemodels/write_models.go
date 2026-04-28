// Package writemodels contains explicit input-only models for the user module.
//
// Write models represent caller intent and must not be reused as responses or
// domain state.
package writemodels

// UserUpdateWriteModel represents user profile update input intent.
type UserUpdateWriteModel struct {
	DisplayName string `json:"display_name"`
}
