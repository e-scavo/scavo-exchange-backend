package auth

import "context"

type WalletIdentity struct {
	ID      string `json:"id"`
	Address string `json:"address"`
}

type WalletIdentityStore interface {
	GetOrCreate(ctx context.Context, address string) (*WalletIdentity, error)
}
