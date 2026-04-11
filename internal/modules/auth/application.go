package auth

import (
	"context"
	"errors"
	"fmt"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
	usermod "github.com/e-scavo/scavo-exchange-backend/internal/modules/user"
	usersettingsmod "github.com/e-scavo/scavo-exchange-backend/internal/modules/usersettings"
)

var (
	ErrApplicationNotConfigured = errors.New("application not configured")
	ErrWalletIdentityStore      = errors.New("wallet identity store error")
)

type Application struct {
	Tokens           *coreauth.TokenService
	Users            *usermod.Service
	UserSettings     *usersettingsmod.Service
	WalletIdentities WalletIdentityStore
}

func NewApplication(tokens *coreauth.TokenService, users *usermod.Service, userSettings *usersettingsmod.Service, walletIdentities WalletIdentityStore) *Application {
	return &Application{
		Tokens:           tokens,
		Users:            users,
		UserSettings:     userSettings,
		WalletIdentities: walletIdentities,
	}
}

func (h HTTPHandlers) Application() *Application {
	return NewApplication(h.Tokens, h.Users, h.UserSettings, h.WalletIdentities)
}

func (a *Application) GetBootstrap(ctx context.Context, claims *coreauth.Claims) (BootstrapResponse, error) {
	if claims == nil {
		return BootstrapResponse{}, ErrUnauthorized
	}
	if a == nil || a.UserSettings == nil {
		return BootstrapResponse{}, ErrApplicationNotConfigured
	}

	svc := NewService(a.Tokens, a.Users, 0)
	user, err := svc.ResolveCurrentUserClaims(ctx, claims)
	if err != nil {
		return BootstrapResponse{}, err
	}

	session := buildSessionViewWithUser(claims, user)
	profile, err := buildProfileViewWithUser(ctx, claims, user, a.WalletIdentities)
	if err != nil {
		return BootstrapResponse{}, err
	}

	settings, err := a.UserSettings.GetOrDefault(ctx, session.UserID)
	if err != nil {
		return BootstrapResponse{}, err
	}

	wallets, err := listWalletReadModelsForUser(ctx, session.UserID, a.WalletIdentities)
	if err != nil {
		return BootstrapResponse{}, fmt.Errorf("%w: %v", ErrWalletIdentityStore, err)
	}

	return BootstrapResponse{
		Session:  session,
		User:     user,
		Profile:  profile,
		Settings: usersettingsmod.ToView(settings),
		Wallets:  buildBootstrapWalletsView(wallets),
	}, nil
}
