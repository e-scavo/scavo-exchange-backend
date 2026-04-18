package auth

import (
	"context"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
	authapp "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth/app"
	authdomain "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth/domain"
	usermod "github.com/e-scavo/scavo-exchange-backend/internal/modules/user"
)

type ProfileWalletView = authapp.ProfileWalletView
type ProfileView = authapp.ProfileView

func buildProfileView(ctx context.Context, claims *coreauth.Claims, users *usermod.Service, walletStore authdomain.WalletIdentityStore) (*ProfileView, error) {
	return authapp.BuildProfileView(ctx, claims, users, walletStore)
}

func buildProfileViewWithUser(ctx context.Context, claims *coreauth.Claims, user *usermod.User, walletStore authdomain.WalletIdentityStore) (*ProfileView, error) {
	return authapp.BuildProfileViewWithUser(ctx, claims, user, walletStore)
}

func mapWalletIdentityToProfileWallet(wallet *authdomain.WalletIdentity) *ProfileWalletView {
	return authapp.MapWalletIdentityToProfileWallet(wallet)
}
