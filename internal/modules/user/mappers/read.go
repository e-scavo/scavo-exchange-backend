// Package mappers centralizes explicit model transformations for the user module.
package mappers

import (
	userdomain "github.com/e-scavo/scavo-exchange-backend/internal/modules/user/domain"
	userreadmodels "github.com/e-scavo/scavo-exchange-backend/internal/modules/user/readmodels"
)

// UserToReadModel maps the canonical user domain model into its output-only
// user read projection.
func UserToReadModel(user *userdomain.User) *userreadmodels.UserReadModel {
	if user == nil {
		return nil
	}

	return &userreadmodels.UserReadModel{
		ID:          user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		LastLoginAt: user.LastLoginAt,
	}
}
