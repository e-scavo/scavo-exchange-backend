// =====================================================
// REPOSITORY LAYER (TRANSITIONAL FACADE)
// =====================================================
//
// This package is NOT the canonical owner of repository implementations yet.
//
// Current state (0.11.4):
// - Actual implementations live in auth root (wallet_*_store_*.go)
// - This package acts as a thin façade/alias layer
// - Domain contracts are defined in auth/domain
//
// Future direction:
// - This package will become the canonical repository layer
// - Implementations may be moved here in a dedicated phase
//
// IMPORTANT:
// Do NOT migrate implementations here without a full refactor plan.
// This is intentional and part of a staged architecture transition.
package repository

import rootauth "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth"

type InMemoryWalletChallengeStore = rootauth.InMemoryWalletChallengeStore

func NewInMemoryWalletChallengeStore() *InMemoryWalletChallengeStore {
	return rootauth.NewInMemoryWalletChallengeStore()
}
