package auth

import (
	"context"
	"strings"
	"time"
)

type WalletLinkingService struct {
	challenges WalletChallengeService
	identities WalletIdentityStore
}

type WalletLinkResult struct {
	Challenge *WalletChallenge  `json:"challenge,omitempty"`
	Linked    *WalletIdentity   `json:"linked_wallet,omitempty"`
	Wallets   []*WalletIdentity `json:"wallets"`
}

func NewWalletLinkingService(challenges *WalletChallengeService, identities WalletIdentityStore) *WalletLinkingService {
	if challenges == nil {
		return &WalletLinkingService{identities: identities}
	}

	return &WalletLinkingService{
		challenges: *challenges,
		identities: identities,
	}
}

func (s *WalletLinkingService) CreateChallenge(ctx context.Context, userID, address, chain string) (*WalletChallenge, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, ErrUnauthorized
	}
	if s == nil {
		return nil, ErrUnauthorized
	}

	return s.challenges.CreateWithOptions(ctx, address, chain, WalletChallengeOptions{
		Purpose:           WalletChallengePurposeLinkWallet,
		RequestedByUserID: userID,
	})
}

func (s *WalletLinkingService) VerifyAndLink(ctx context.Context, userID, challengeID, address, signature string) (*WalletLinkResult, error) {
	userID = strings.TrimSpace(userID)
	if s == nil || s.identities == nil || userID == "" {
		return nil, ErrUnauthorized
	}

	address = normalizeWalletAddress(address)
	if !evmAddressRE.MatchString(address) {
		return nil, ErrInvalidWalletAddress
	}

	challenge, err := s.challenges.Get(ctx, strings.TrimSpace(challengeID))
	if err != nil {
		return nil, err
	}
	if challenge == nil {
		return nil, ErrWalletChallengeNotFound
	}
	if normalizeWalletChallengePurpose(challenge.Purpose) != WalletChallengePurposeLinkWallet {
		return nil, ErrWalletChallengePurpose
	}
	if strings.TrimSpace(challenge.RequestedByUserID) != userID {
		return nil, ErrWalletLinkChallengeMismatch
	}
	if normalizeWalletAddress(challenge.Address) != address {
		return nil, ErrInvalidWalletSignature
	}

	recoveredAddress, err := recoverWalletAddress(challenge.Message, signature)
	if err != nil {
		return nil, err
	}
	if normalizeWalletAddress(recoveredAddress) != address {
		return nil, ErrInvalidWalletSignature
	}

	identity, err := s.identities.GetOrCreate(ctx, address)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(identity.UserID) == userID {
		return nil, ErrWalletAlreadyLinkedToUser
	}
	if strings.TrimSpace(identity.UserID) != "" && strings.TrimSpace(identity.UserID) != userID {
		return nil, ErrWalletIdentityAlreadyLinked
	}

	usedAt := time.Now().UTC()
	challenge, err = s.challenges.MarkUsed(ctx, challenge.ID, usedAt)
	if err != nil {
		return nil, err
	}

	identity, err = s.identities.AttachUser(ctx, identity.ID, userID, false)
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

	return &WalletLinkResult{
		Challenge: challenge,
		Linked:    identity,
		Wallets:   wallets,
	}, nil
}
