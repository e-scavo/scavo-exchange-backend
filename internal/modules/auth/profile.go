package auth

import (
	"context"
	"strings"
	"time"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
	usermod "github.com/e-scavo/scavo-exchange-backend/internal/modules/user"
)

type ProfileWalletView struct {
	ID         string     `json:"id"`
	Address    string     `json:"address"`
	IsPrimary  bool       `json:"is_primary"`
	Status     string     `json:"status"`
	LinkedAt   *time.Time `json:"linked_at,omitempty"`
	DetachedAt *time.Time `json:"detached_at,omitempty"`
}

type ProfileView struct {
	User                *usermod.User        `json:"user,omitempty"`
	UserID              string               `json:"user_id"`
	AuthMethod          string               `json:"auth_method,omitempty"`
	WalletID            string               `json:"wallet_id,omitempty"`
	WalletAddress       string               `json:"wallet_address,omitempty"`
	Chain               string               `json:"chain,omitempty"`
	PrimaryWallet       *ProfileWalletView   `json:"primary_wallet,omitempty"`
	Wallets             []*ProfileWalletView `json:"wallets"`
	WalletCount         int                  `json:"wallet_count"`
	ActiveWalletCount   int                  `json:"active_wallet_count"`
	DetachedWalletCount int                  `json:"detached_wallet_count"`
	HasWalletSession    bool                 `json:"has_wallet_session"`
}

func buildProfileView(ctx context.Context, claims *coreauth.Claims, users *usermod.Service, walletStore WalletIdentityStore) (*ProfileView, error) {
	svc := NewService(nil, users, 24*time.Hour)
	user, err := svc.ResolveCurrentUserClaims(ctx, claims)
	if err != nil {
		return nil, err
	}

	return buildProfileViewWithUser(ctx, claims, user, walletStore)
}

func buildProfileViewWithUser(ctx context.Context, claims *coreauth.Claims, user *usermod.User, walletStore WalletIdentityStore) (*ProfileView, error) {
	view := &ProfileView{
		User:             user,
		UserID:           strings.TrimSpace(claims.UserID),
		AuthMethod:       strings.TrimSpace(claims.AuthMethod),
		WalletID:         strings.TrimSpace(claims.WalletID),
		WalletAddress:    normalizeWalletAddress(claims.WalletAddress),
		Chain:            normalizeChain(claims.Chain),
		Wallets:          []*ProfileWalletView{},
		HasWalletSession: strings.TrimSpace(claims.WalletAddress) != "",
	}

	if view.AuthMethod == "" {
		view.AuthMethod = "password_dev"
	}
	if walletStore == nil || view.UserID == "" {
		return view, nil
	}

	wallets, err := walletStore.ListByUser(ctx, view.UserID)
	if err != nil {
		return nil, err
	}

	for _, wallet := range wallets {
		mapped := mapWalletIdentityToProfileWallet(wallet)
		if mapped == nil {
			continue
		}
		view.Wallets = append(view.Wallets, mapped)
		view.WalletCount++
		switch mapped.Status {
		case "active":
			view.ActiveWalletCount++
		case "detached":
			view.DetachedWalletCount++
		}
		if mapped.IsPrimary && view.PrimaryWallet == nil {
			copy := *mapped
			view.PrimaryWallet = &copy
		}
	}

	return view, nil
}

func mapWalletIdentityToProfileWallet(wallet *WalletIdentity) *ProfileWalletView {
	mapped := mapWalletIdentityToReadModel(wallet)
	if mapped == nil {
		return nil
	}

	return &ProfileWalletView{
		ID:         mapped.ID,
		Address:    mapped.Address,
		IsPrimary:  mapped.IsPrimary,
		Status:     mapped.Status,
		LinkedAt:   mapped.LinkedAt,
		DetachedAt: mapped.DetachedAt,
	}
}
