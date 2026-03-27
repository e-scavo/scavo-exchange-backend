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

func (s *InMemoryWalletIdentityStore) GetByAddress(ctx context.Context, address string) (*WalletIdentity, error) {
	address = normalizeWalletAddress(address)
	if !evmAddressRE.MatchString(address) {
		return nil, ErrInvalidWalletAddress
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	identity, ok := s.items[address]
	if !ok || identity == nil {
		return nil, ErrWalletIdentityNotFound
	}

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

func (s *InMemoryWalletIdentityStore) ReassignUser(ctx context.Context, walletID, fromUserID, toUserID string, primary bool) (*WalletIdentity, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	walletID = strings.TrimSpace(walletID)
	fromUserID = strings.TrimSpace(fromUserID)
	toUserID = strings.TrimSpace(toUserID)
	if walletID == "" || fromUserID == "" || toUserID == "" {
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
	if strings.TrimSpace(target.UserID) != fromUserID {
		if strings.TrimSpace(target.UserID) == toUserID {
			return nil, ErrWalletMergeSameUser
		}
		return nil, ErrWalletIdentityAlreadyLinked
	}

	if primary {
		for _, identity := range s.items {
			if identity == nil {
				continue
			}
			if strings.TrimSpace(identity.UserID) == toUserID && identity.ID != walletID {
				identity.IsPrimary = false
			}
		}
	}

	target.UserID = toUserID
	target.IsPrimary = primary
	if target.LinkedAt == nil {
		now := time.Now().UTC()
		target.LinkedAt = &now
	}

	cp := *target
	return &cp, nil
}

func (s *InMemoryWalletIdentityStore) MergeUsers(ctx context.Context, sourceUserID, targetUserID string) ([]*WalletIdentity, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sourceUserID = strings.TrimSpace(sourceUserID)
	targetUserID = strings.TrimSpace(targetUserID)
	if sourceUserID == "" || targetUserID == "" {
		return nil, ErrUnauthorized
	}
	if sourceUserID == targetUserID {
		return s.listByUserLocked(targetUserID), nil
	}

	targetHasPrimary := false
	for _, identity := range s.items {
		if identity != nil && strings.TrimSpace(identity.UserID) == targetUserID && identity.IsPrimary {
			targetHasPrimary = true
			break
		}
	}

	for _, identity := range s.items {
		if identity == nil || strings.TrimSpace(identity.UserID) != sourceUserID {
			continue
		}
		identity.UserID = targetUserID
		if targetHasPrimary {
			identity.IsPrimary = false
		} else if identity.IsPrimary {
			targetHasPrimary = true
		}
		if identity.LinkedAt == nil {
			now := time.Now().UTC()
			identity.LinkedAt = &now
		}
	}

	return s.listByUserLocked(targetUserID), nil
}

func (s *InMemoryWalletIdentityStore) SetPrimary(ctx context.Context, userID, address string) (*WalletIdentity, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	userID = strings.TrimSpace(userID)
	address = normalizeWalletAddress(address)
	if userID == "" {
		return nil, ErrUnauthorized
	}
	if !evmAddressRE.MatchString(address) {
		return nil, ErrInvalidWalletAddress
	}

	target, ok := s.items[address]
	if !ok || target == nil {
		return nil, ErrWalletIdentityNotFound
	}
	if strings.TrimSpace(target.UserID) != userID {
		return nil, ErrWalletNotOwnedByUser
	}

	for _, identity := range s.items {
		if identity == nil {
			continue
		}
		if strings.TrimSpace(identity.UserID) == userID {
			identity.IsPrimary = identity.Address == address
		}
	}

	cp := *target
	return &cp, nil
}

func (s *InMemoryWalletIdentityStore) DetachUser(ctx context.Context, userID, address string) (*WalletIdentity, []*WalletIdentity, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	userID = strings.TrimSpace(userID)
	address = normalizeWalletAddress(address)
	if userID == "" {
		return nil, nil, ErrUnauthorized
	}
	if !evmAddressRE.MatchString(address) {
		return nil, nil, ErrInvalidWalletAddress
	}

	target, ok := s.items[address]
	if !ok || target == nil {
		return nil, nil, ErrWalletIdentityNotFound
	}
	if strings.TrimSpace(target.UserID) != userID {
		return nil, nil, ErrWalletNotOwnedByUser
	}

	now := time.Now().UTC()
	target.UserID = ""
	target.LinkedAt = nil
	target.DetachedAt = &now
	target.IsPrimary = false

	cp := *target
	remaining := s.listByUserLocked(userID)
	return &cp, remaining, nil
}

func (s *InMemoryWalletIdentityStore) ListByUser(ctx context.Context, userID string) ([]*WalletIdentity, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userID = strings.TrimSpace(userID)
	if userID == "" {
		return []*WalletIdentity{}, nil
	}

	return s.listByUserLocked(userID), nil
}

func (s *InMemoryWalletIdentityStore) listByUserLocked(userID string) []*WalletIdentity {
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

	return out
}
