package auth

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"
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

func (s *InMemoryWalletIdentityStore) AttachUser(ctx context.Context, walletID, userID string, primary bool) (*WalletIdentity, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	walletID = strings.TrimSpace(walletID)
	userID = strings.TrimSpace(userID)
	if walletID == "" || userID == "" {
		return nil, ErrUnauthorized
	}

	var target *WalletIdentity
	for _, identity := range s.items {
		if identity == nil || identity.ID != walletID {
			continue
		}
		target = identity
		break
	}

	if target == nil {
		return nil, ErrWalletIdentityNotFound
	}

	if strings.TrimSpace(target.UserID) != "" && strings.TrimSpace(target.UserID) != userID {
		return nil, ErrWalletIdentityAlreadyLinked
	}

	if primary {
		for _, identity := range s.items {
			if identity == nil {
				continue
			}
			if strings.TrimSpace(identity.UserID) == userID && identity.ID != walletID {
				identity.IsPrimary = false
			}
		}
	}

	now := time.Now().UTC()
	target.UserID = userID
	if target.LinkedAt == nil {
		target.LinkedAt = &now
	}
	target.IsPrimary = primary

	cp := *target
	return &cp, nil
}

func (s *InMemoryWalletIdentityStore) ListByUser(ctx context.Context, userID string) ([]*WalletIdentity, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userID = strings.TrimSpace(userID)
	if userID == "" {
		return []*WalletIdentity{}, nil
	}

	out := make([]*WalletIdentity, 0)
	for _, identity := range s.items {
		if identity == nil {
			continue
		}
		if strings.TrimSpace(identity.UserID) != userID {
			continue
		}

		cp := *identity
		out = append(out, &cp)
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].IsPrimary != out[j].IsPrimary {
			return out[i].IsPrimary
		}
		if out[i].LinkedAt != nil && out[j].LinkedAt != nil {
			if !out[i].LinkedAt.Equal(*out[j].LinkedAt) {
				return out[i].LinkedAt.Before(*out[j].LinkedAt)
			}
		}
		return out[i].Address < out[j].Address
	})

	return out, nil
}
