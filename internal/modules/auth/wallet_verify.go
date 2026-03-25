package auth

import (
	"context"
	"strings"
	"time"
)

type WalletVerificationService struct {
	challenges WalletChallengeService
	login      *Service
	identities WalletIdentityStore
}

func NewWalletVerificationService(
	challenges *WalletChallengeService,
	login *Service,
	identities WalletIdentityStore,
) *WalletVerificationService {
	return &WalletVerificationService{
		challenges: *challenges,
		login:      login,
		identities: identities,
	}
}

func (s *WalletVerificationService) VerifyAndLogin(ctx context.Context, challengeID, address, signature string) (*LoginResult, *WalletChallenge, error) {
	if s == nil || s.login == nil || s.identities == nil {
		return nil, nil, ErrUnauthorized
	}

	address = normalizeWalletAddress(address)
	if !evmAddressRE.MatchString(address) {
		return nil, nil, ErrInvalidWalletAddress
	}

	challenge, err := s.challenges.Get(ctx, strings.TrimSpace(challengeID))
	if err != nil {
		return nil, nil, err
	}
	if challenge == nil {
		return nil, nil, ErrWalletChallengeNotFound
	}
	if normalizeWalletAddress(challenge.Address) != address {
		return nil, nil, ErrInvalidWalletSignature
	}

	recoveredAddress, err := recoverWalletAddress(challenge.Message, signature)
	if err != nil {
		return nil, nil, err
	}
	if normalizeWalletAddress(recoveredAddress) != address {
		return nil, nil, ErrInvalidWalletSignature
	}

	usedAt := time.Now().UTC()
	challenge, err = s.challenges.MarkUsed(ctx, challenge.ID, usedAt)
	if err != nil {
		return nil, nil, err
	}

	identity, err := s.identities.GetOrCreate(ctx, address)
	if err != nil {
		return nil, nil, err
	}

	result, err := s.login.LoginWallet(ctx, identity.ID, address, challenge.Chain)
	if err != nil {
		return nil, nil, err
	}

	return result, challenge, nil
}
