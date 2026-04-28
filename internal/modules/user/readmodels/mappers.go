package readmodels

import userdomain "github.com/e-scavo/scavo-exchange-backend/internal/modules/user/domain"

// FromUser maps the canonical user domain model into its output-only user read
// model projection.
func FromUser(user *userdomain.User) *UserReadModel {
	if user == nil {
		return nil
	}

	return &UserReadModel{
		ID:          user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		LastLoginAt: user.LastLoginAt,
	}
}
