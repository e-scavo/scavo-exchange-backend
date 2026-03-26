package auth

import "context"

type WalletIdentity struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	UserID  string `json:"user_id,omitempty"`
}

type WalletIdentityStore interface {
	GetOrCreate(ctx context.Context, address string) (*WalletIdentity, error)
	AttachUser(ctx context.Context, walletID, userID string) (*WalletIdentity, error)
}
