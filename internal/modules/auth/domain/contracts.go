package domain

import (
	"context"
	"time"
)

const (
	WalletDetachReasonNotOwnedByUser   = "wallet_not_owned_by_user"
	WalletDetachReasonWalletIsPrimary  = "wallet_is_primary"
	WalletDetachReasonUserWouldBeEmpty = "user_would_have_no_wallets"
)

type WalletIdentity struct {
	ID         string     `json:"id"`
	Address    string     `json:"address"`
	UserID     string     `json:"user_id,omitempty"`
	LinkedAt   *time.Time `json:"linked_at,omitempty"`
	DetachedAt *time.Time `json:"detached_at,omitempty"`
	IsPrimary  bool       `json:"is_primary"`
}

type WalletIdentityStore interface {
	GetOrCreate(ctx context.Context, address string) (*WalletIdentity, error)
	GetByAddress(ctx context.Context, address string) (*WalletIdentity, error)
	AttachUser(ctx context.Context, walletID, userID string, primary bool) (*WalletIdentity, error)
	ReassignUser(ctx context.Context, walletID, fromUserID, toUserID string, primary bool) (*WalletIdentity, error)
	MergeUsers(ctx context.Context, sourceUserID, targetUserID string) ([]*WalletIdentity, error)
	SetPrimary(ctx context.Context, userID, address string) (*WalletIdentity, error)
	DetachUser(ctx context.Context, userID, address string) (*WalletIdentity, []*WalletIdentity, error)
	ListByUser(ctx context.Context, userID string) ([]*WalletIdentity, error)
}

type WalletChallenge struct {
	ID                string     `json:"id"`
	Address           string     `json:"address"`
	Chain             string     `json:"chain"`
	Nonce             string     `json:"nonce"`
	Message           string     `json:"message"`
	Purpose           string     `json:"purpose"`
	RequestedByUserID string     `json:"requested_by_user_id,omitempty"`
	IssuedAt          time.Time  `json:"issued_at"`
	ExpiresAt         time.Time  `json:"expires_at"`
	UsedAt            *time.Time `json:"used_at,omitempty"`
}

type WalletChallengeStore interface {
	Create(ctx context.Context, address, chain string, ttl time.Duration) (*WalletChallenge, error)
	CreateWithOptions(ctx context.Context, address, chain string, opts WalletChallengeOptions) (*WalletChallenge, error)
	Get(ctx context.Context, challengeID string) (*WalletChallenge, error)
	MarkUsed(ctx context.Context, challengeID string, usedAt time.Time) error
}

type WalletChallengeOptions struct {
	Purpose           string
	RequestedByUserID string
}
