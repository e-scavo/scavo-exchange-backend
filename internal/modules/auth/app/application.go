package app

import (
	"time"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
	rootauth "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth"
	usermod "github.com/e-scavo/scavo-exchange-backend/internal/modules/user"
	usersettingsmod "github.com/e-scavo/scavo-exchange-backend/internal/modules/usersettings"
)

var (
	ErrApplicationNotConfigured = rootauth.ErrApplicationNotConfigured
	ErrWalletIdentityStore      = rootauth.ErrWalletIdentityStore
)

type Application = rootauth.Application

func NewApplication(tokens *coreauth.TokenService, ttl time.Duration, users *usermod.Service, userSettings *usersettingsmod.Service, publicBaseURL string, challengeTTL time.Duration, challenges rootauth.WalletChallengeStore, walletIdentities rootauth.WalletIdentityStore) *Application {
	return rootauth.NewApplication(tokens, ttl, users, userSettings, publicBaseURL, challengeTTL, challenges, walletIdentities)
}
