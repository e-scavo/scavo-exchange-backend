package auth

import (
	"context"
	"strings"
)

type WalletPrimaryService struct {
	identities WalletIdentityStore
}

type WalletPrimaryResult struct {
	Primary *WalletIdentity   `json:"primary_wallet,omitempty"`
	Wallets []*WalletIdentity `json:"wallets"`
}

func NewWalletPrimaryService(identities WalletIdentityStore) *WalletPrimaryService {
	return &WalletPrimaryService{identities: identities}
}

func (s *WalletPrimaryService) SetPrimary(ctx context.Context, userID, address string) (*WalletPrimaryResult, error) {
	userID = strings.TrimSpace(userID)
	if s == nil || s.identities == nil || userID == "" {
		return nil, ErrUnauthorized
	}

	address = normalizeWalletAddress(address)
	if !evmAddressRE.MatchString(address) {
		return nil, ErrInvalidWalletAddress
	}

	primary, err := s.identities.SetPrimary(ctx, userID, address)
	if err != nil {
		return nil, err
	}

	wallets, err := s.identities.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if wallets == nil {
		wallets = []*WalletIdentity{}
	}

	return &WalletPrimaryResult{
		Primary: primary,
		Wallets: wallets,
	}, nil
}
