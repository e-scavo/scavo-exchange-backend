package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	usermod "github.com/e-scavo/scavo-exchange-backend/internal/modules/user"
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
	if normalizeWalletChallengePurpose(challenge.Purpose) != WalletChallengePurposeAuthBootstrap {
		return nil, nil, ErrWalletChallengePurpose
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

	user, identity, err := s.resolveLinkedUser(ctx, identity, address)
	if err != nil {
		return nil, nil, err
	}

	result, err := s.login.LoginWalletForUser(ctx, user, identity.ID, address, challenge.Chain)
	if err != nil {
		return nil, nil, err
	}

	return result, challenge, nil
}

func (s *WalletVerificationService) resolveLinkedUser(ctx context.Context, identity *WalletIdentity, address string) (*usermod.User, *WalletIdentity, error) {
	if identity == nil {
		return nil, nil, ErrUnauthorized
	}

	if strings.TrimSpace(identity.UserID) != "" {
		if s.login == nil || s.login.users == nil {
			return walletUser(address), identity, nil
		}

		user, err := s.login.users.GetByID(ctx, identity.UserID, walletUserEmail(address))
		if err == nil {
			return user, identity, nil
		}
		if !errors.Is(err, usermod.ErrUserNotFound) {
			return nil, nil, err
		}
	}

	if s.login == nil || s.login.users == nil {
		linked := walletUser(address)
		identity.UserID = linked.ID
		now := time.Now().UTC()
		identity.LinkedAt = &now
		identity.IsPrimary = true
		return linked, identity, nil
	}

	user, err := s.login.users.ResolveOrCreateWalletUser(ctx, address)
	if err != nil {
		return nil, nil, err
	}

	identity, err = s.identities.AttachUser(ctx, identity.ID, user.ID, true)
	if err != nil {
		return nil, nil, err
	}

	return user, identity, nil
}
