package app

import (
	"time"

	rootauth "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth"
)

type WalletChallengeService = rootauth.WalletChallengeService
type WalletLinkingService = rootauth.WalletLinkingService
type WalletAccountMergeService = rootauth.WalletAccountMergeService
type WalletPrimaryService = rootauth.WalletPrimaryService
type WalletDetachService = rootauth.WalletDetachService

func NewWalletChallengeService(store rootauth.WalletChallengeStore, publicBaseURL string, ttl time.Duration) *WalletChallengeService {
	return rootauth.NewWalletChallengeService(store, publicBaseURL, ttl)
}

func NewWalletLinkingService(challenges *WalletChallengeService, identities rootauth.WalletIdentityStore) *WalletLinkingService {
	return rootauth.NewWalletLinkingService(challenges, identities)
}

func NewWalletAccountMergeService(challenges *WalletChallengeService, identities rootauth.WalletIdentityStore) *WalletAccountMergeService {
	return rootauth.NewWalletAccountMergeService(challenges, identities)
}

func NewWalletPrimaryService(identities rootauth.WalletIdentityStore) *WalletPrimaryService {
	return rootauth.NewWalletPrimaryService(identities)
}

func NewWalletDetachService(identities rootauth.WalletIdentityStore) *WalletDetachService {
	return rootauth.NewWalletDetachService(identities)
}
