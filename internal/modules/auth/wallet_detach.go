package auth

import (
	"context"
	"strings"
)

const (
	WalletDetachReasonNotOwnedByUser   = "wallet_not_owned_by_user"
	WalletDetachReasonWalletIsPrimary  = "wallet_is_primary"
	WalletDetachReasonUserWouldBeEmpty = "user_would_have_no_wallets"
)

type WalletDetachCheckResult struct {
	WalletAddress    string   `json:"wallet_address"`
	Eligible         bool     `json:"eligible"`
	IsPrimary        bool     `json:"is_primary"`
	OwnedWalletCount int      `json:"owned_wallet_count"`
	Reasons          []string `json:"reasons"`
}

type WalletDetachService struct {
	identities WalletIdentityStore
}

func NewWalletDetachService(identities WalletIdentityStore) *WalletDetachService {
	return &WalletDetachService{identities: identities}
}

func (s *WalletDetachService) CheckEligibility(ctx context.Context, userID, address string) (*WalletDetachCheckResult, error) {
	userID = strings.TrimSpace(userID)
	if s == nil || s.identities == nil || userID == "" {
		return nil, ErrUnauthorized
	}

	address = normalizeWalletAddress(address)
	if !evmAddressRE.MatchString(address) {
		return nil, ErrInvalidWalletAddress
	}

	identity, err := s.identities.GetByAddress(ctx, address)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(identity.UserID) != userID {
		return nil, ErrWalletNotOwnedByUser
	}

	wallets, err := s.identities.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if wallets == nil {
		wallets = []*WalletIdentity{}
	}

	result := &WalletDetachCheckResult{
		WalletAddress:    address,
		Eligible:         true,
		IsPrimary:        identity.IsPrimary,
		OwnedWalletCount: len(wallets),
		Reasons:          []string{},
	}

	if identity.IsPrimary {
		result.Reasons = append(result.Reasons, WalletDetachReasonWalletIsPrimary)
	}
	if len(wallets) <= 1 {
		result.Reasons = append(result.Reasons, WalletDetachReasonUserWouldBeEmpty)
	}
	if len(result.Reasons) > 0 {
		result.Eligible = false
	}

	return result, nil
}
