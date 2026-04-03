package auth

import (
	"context"
	"strings"
	"time"
)

type WalletAccountMergeService struct {
	challenges WalletChallengeService
	identities WalletIdentityStore
}

type WalletAccountMergeResult struct {
	Challenge    *WalletChallenge  `json:"challenge,omitempty"`
	MergedWallet *WalletIdentity   `json:"merged_wallet,omitempty"`
	Wallets      []*WalletIdentity `json:"wallets"`
	SourceUserID string            `json:"source_user_id"`
	TargetUserID string            `json:"target_user_id"`
}

func NewWalletAccountMergeService(challenges *WalletChallengeService, identities WalletIdentityStore) *WalletAccountMergeService {
	if challenges == nil {
		return &WalletAccountMergeService{identities: identities}
	}

	return &WalletAccountMergeService{
		challenges: *challenges,
		identities: identities,
	}
}

func (s *WalletAccountMergeService) CreateChallenge(ctx context.Context, userID, address, chain string) (*WalletChallenge, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, ErrUnauthorized
	}
	if s == nil {
		return nil, ErrUnauthorized
	}

	return s.challenges.CreateWithOptions(ctx, address, chain, WalletChallengeOptions{
		Purpose:           WalletChallengePurposeAccountMerge,
		RequestedByUserID: userID,
	})
}

func (s *WalletAccountMergeService) VerifyAndMerge(ctx context.Context, targetUserID, challengeID, address, signature string) (*WalletAccountMergeResult, error) {
	targetUserID = strings.TrimSpace(targetUserID)
	if s == nil || s.identities == nil || targetUserID == "" {
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
	if purpose, ok := canonicalWalletChallengePurpose(challenge.Purpose); !ok || purpose != WalletChallengePurposeAccountMerge {
		return nil, ErrWalletChallengePurpose
	}
	if strings.TrimSpace(challenge.RequestedByUserID) != targetUserID {
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

	identity, err := s.identities.GetByAddress(ctx, address)
	if err != nil {
		return nil, err
	}
	if identity == nil || strings.TrimSpace(identity.UserID) == "" {
		return nil, ErrWalletMergeSourceNotLinked
	}

	sourceUserID := strings.TrimSpace(identity.UserID)
	if sourceUserID == targetUserID {
		return nil, ErrWalletMergeSameUser
	}

	sourceWallets, err := s.identities.ListByUser(ctx, sourceUserID)
	if err != nil {
		return nil, err
	}
	if len(sourceWallets) == 0 {
		return nil, ErrWalletMergeSourceNotLinked
	}

	challenge, err = s.challenges.MarkUsed(ctx, challenge.ID, time.Now().UTC())
	if err != nil {
		return nil, err
	}

	wallets, err := s.identities.MergeUsers(ctx, sourceUserID, targetUserID)
	if err != nil {
		return nil, err
	}
	if wallets == nil {
		wallets = []*WalletIdentity{}
	}

	var mergedWallet *WalletIdentity
	for _, wallet := range wallets {
		if wallet != nil && strings.EqualFold(wallet.Address, address) {
			mergedWallet = wallet
			break
		}
	}

	return &WalletAccountMergeResult{
		Challenge:    challenge,
		MergedWallet: mergedWallet,
		Wallets:      wallets,
		SourceUserID: sourceUserID,
		TargetUserID: targetUserID,
	}, nil
}
