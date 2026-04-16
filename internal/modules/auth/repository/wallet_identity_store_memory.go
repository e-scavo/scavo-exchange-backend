package repository

import rootauth "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth"

type InMemoryWalletIdentityStore = rootauth.InMemoryWalletIdentityStore

func NewInMemoryWalletIdentityStore() *InMemoryWalletIdentityStore {
	return rootauth.NewInMemoryWalletIdentityStore()
}
