package app

import (
	"time"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
	rootauth "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth"
	usermod "github.com/e-scavo/scavo-exchange-backend/internal/modules/user"
)

var (
	ErrInvalidCredentials = rootauth.ErrInvalidCredentials
	ErrUnauthorized       = rootauth.ErrUnauthorized
)

type Service = rootauth.Service
type LoginResult = rootauth.LoginResult

// type SessionView = rootauth.SessionView

func NewService(tokens *coreauth.TokenService, users *usermod.Service, ttl time.Duration) *Service {
	return rootauth.NewService(tokens, users, ttl)
}
