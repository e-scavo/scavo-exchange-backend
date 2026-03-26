package auth

import (
	"context"
	"time"
)

type WalletIdentity struct {
	ID        string     `json:"id"`
	Address   string     `json:"address"`
	UserID    string     `json:"user_id,omitempty"`
	LinkedAt  *time.Time `json:"linked_at,omitempty"`
	IsPrimary bool       `json:"is_primary"`
}

type WalletIdentityStore interface {
	GetOrCreate(ctx context.Context, address string) (*WalletIdentity, error)
	GetByAddress(ctx context.Context, address string) (*WalletIdentity, error)
	AttachUser(ctx context.Context, walletID, userID string, primary bool) (*WalletIdentity, error)
	ReassignUser(ctx context.Context, walletID, fromUserID, toUserID string, primary bool) (*WalletIdentity, error)
	MergeUsers(ctx context.Context, sourceUserID, targetUserID string) ([]*WalletIdentity, error)
	SetPrimary(ctx context.Context, userID, address string) (*WalletIdentity, error)
	ListByUser(ctx context.Context, userID string) ([]*WalletIdentity, error)
}
