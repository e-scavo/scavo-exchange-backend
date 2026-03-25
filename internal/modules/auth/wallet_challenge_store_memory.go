package auth

import (
	"context"
	"errors"
	"sync"
	"time"
)

type InMemoryWalletChallengeStore struct {
	mu    sync.RWMutex
	items map[string]*WalletChallenge
}

func NewInMemoryWalletChallengeStore() *InMemoryWalletChallengeStore {
	return &InMemoryWalletChallengeStore{
		items: make(map[string]*WalletChallenge),
	}
}

func (s *InMemoryWalletChallengeStore) Save(ctx context.Context, challenge *WalletChallenge) error {
	if challenge == nil || challenge.ID == "" {
		return errors.New("invalid wallet challenge")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.cleanupLocked(time.Now().UTC())
	s.items[challenge.ID] = cloneWalletChallenge(challenge)

	return nil
}

func (s *InMemoryWalletChallengeStore) GetByID(ctx context.Context, id string) (*WalletChallenge, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cleanupLocked(time.Now().UTC())

	challenge, ok := s.items[id]
	if !ok {
		return nil, ErrWalletChallengeNotFound
	}

	return cloneWalletChallenge(challenge), nil
}

func (s *InMemoryWalletChallengeStore) Use(ctx context.Context, id string, usedAt time.Time) (*WalletChallenge, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	s.cleanupLocked(now)

	challenge, ok := s.items[id]
	if !ok {
		return nil, ErrWalletChallengeNotFound
	}
	if challenge.UsedAt != nil {
		return nil, ErrChallengeUsed
	}
	if now.After(challenge.ExpiresAt) {
		delete(s.items, id)
		return nil, ErrChallengeExpired
	}

	ts := usedAt.UTC()
	challenge.UsedAt = &ts
	return cloneWalletChallenge(challenge), nil
}

func (s *InMemoryWalletChallengeStore) cleanupLocked(now time.Time) {
	for id, ch := range s.items {
		if ch == nil {
			delete(s.items, id)
			continue
		}
		if now.After(ch.ExpiresAt) {
			delete(s.items, id)
		}
	}
}

func cloneWalletChallenge(ch *WalletChallenge) *WalletChallenge {
	if ch == nil {
		return nil
	}

	cp := *ch
	if ch.UsedAt != nil {
		ts := *ch.UsedAt
		cp.UsedAt = &ts
	}

	return &cp
}
