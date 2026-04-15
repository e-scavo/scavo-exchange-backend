package user

import userapp "github.com/e-scavo/scavo-exchange-backend/internal/modules/user/app"

var (
	ErrEmptyUserID        = userapp.ErrEmptyUserID
	ErrEmptyDisplayName   = userapp.ErrEmptyDisplayName
	ErrDisplayNameTooLong = userapp.ErrDisplayNameTooLong
)

type Service = userapp.Service

func NewService(repo Repository) *Service {
	return userapp.NewService(repo)
}
