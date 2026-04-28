package domain

// UserSettingsUpdateInput represents canonical settings update intent after
// transport-level write models have been decoded and mapped.
type UserSettingsUpdateInput struct {
	Preferences map[string]any
}
