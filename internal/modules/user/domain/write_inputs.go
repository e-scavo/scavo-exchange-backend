package domain

// UserUpdateInput represents canonical user update intent after transport-level
// write models have been decoded and mapped.
type UserUpdateInput struct {
	DisplayName string
}
