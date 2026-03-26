package auth

import (
	"context"
	"strings"
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

func (s *InMemoryWalletIdentityStore) AttachUser(ctx context.Context, walletID, userID string) (*WalletIdentity, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, identity := range s.items {
		if identity == nil || identity.ID != walletID {
			continue
		}

		identity.UserID = strings.TrimSpace(userID)
		cp := *identity
		return &cp, nil
	}

	return nil, ErrWalletIdentityNotFound
}
