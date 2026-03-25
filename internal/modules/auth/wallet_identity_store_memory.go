package auth

import (
	"context"
	"sync"
)

type InMemoryWalletIdentityStore struct {
	mu    sync.RWMutex
	items map[string]*WalletIdentity
}

func NewInMemoryWalletIdentityStore() *InMemoryWalletIdentityStore {
	return &InMemoryWalletIdentityStore{
		items: make(map[string]*WalletIdentity),
	}
}

func (s *InMemoryWalletIdentityStore) GetOrCreate(ctx context.Context, address string) (*WalletIdentity, error) {
	address = normalizeWalletAddress(address)
	if !evmAddressRE.MatchString(address) {
		return nil, ErrInvalidWalletAddress
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if existing, ok := s.items[address]; ok && existing != nil {
		cp := *existing
		return &cp, nil
	}

	identity := &WalletIdentity{
		ID:      walletUserID(address),
		Address: address,
	}
	s.items[address] = identity

	cp := *identity
	return &cp, nil
}
