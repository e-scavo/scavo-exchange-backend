package auth

import (
	"context"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
	authdomain "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth/domain"
)

// AuthSessionProvider defines the session-oriented auth application boundary
// consumed by HTTP handlers.
type AuthSessionProvider interface {
	Login(ctx context.Context, email, password string) (LoginResponse, error)
	GetMe(ctx context.Context, claims *coreauth.Claims) (MeResponse, error)
	GetSession(ctx context.Context, claims *coreauth.Claims) (SessionResponse, error)
	GetBootstrap(ctx context.Context, claims *coreauth.Claims) (BootstrapResponse, error)
}

// AuthenticatedAccountProvider defines authenticated account-surface operations
// that must stay outside the transport handler implementation.
type AuthenticatedAccountProvider interface {
	UpdateProfile(ctx context.Context, claims *coreauth.Claims, input authdomain.ProfileUpdateInput) (MeResponse, error)
	GetSettings(ctx context.Context, claims *coreauth.Claims) (MeSettingsResponse, error)
	UpdateSettings(ctx context.Context, claims *coreauth.Claims, input authdomain.SettingsUpdateInput) (MeSettingsResponse, error)
}

// AuthWalletProvider defines wallet bootstrap and authenticated wallet-management
// operations exposed to HTTP handlers through one provider boundary.
type AuthWalletProvider interface {
	CreateWalletChallenge(ctx context.Context, address, chain string) (WalletChallengeResponse, error)
	VerifyWallet(ctx context.Context, challengeID, address, signature string) (WalletVerifyResponse, error)
	ListWallets(ctx context.Context, userID string, query WalletsQuery) (WalletsResponse, error)
	CreateWalletLinkChallenge(ctx context.Context, userID, address, chain string) (WalletLinkChallengeResponse, error)
	VerifyWalletLink(ctx context.Context, userID, challengeID, address, signature string) (WalletLinkVerifyResponse, error)
	CreateWalletAccountMergeChallenge(ctx context.Context, userID, address, chain string) (WalletAccountMergeChallengeResponse, error)
	VerifyWalletAccountMerge(ctx context.Context, userID, challengeID, address, signature string) (WalletAccountMergeVerifyResponse, error)
	SetPrimaryWallet(ctx context.Context, userID, address string) (WalletPrimarySetResponse, error)
	CheckWalletDetach(ctx context.Context, userID, address string) (WalletDetachCheckResponse, error)
	ExecuteWalletDetach(ctx context.Context, userID, address string) (WalletDetachExecuteResponse, error)
}

// AuthProvider is the composite handler-facing provider contract introduced by
// Phase 0.13.3. Smaller provider interfaces remain available for focused tests.
type AuthProvider interface {
	AuthSessionProvider
	AuthenticatedAccountProvider
	AuthWalletProvider
}
