package repository

import rootauth "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth"

type InMemoryWalletChallengeStore = rootauth.InMemoryWalletChallengeStore

func NewInMemoryWalletChallengeStore() *InMemoryWalletChallengeStore {
	return rootauth.NewInMemoryWalletChallengeStore()
}
