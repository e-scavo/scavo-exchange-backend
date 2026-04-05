package user

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ResolveOrCreateDevUser(ctx context.Context, email string) (*User, error) {
	email = normalizeEmail(email)
	if email == "" {
		return nil, fmt.Errorf("empty email")
	}

	if s == nil || s.repo == nil {
		now := time.Now().UTC()
		return &User{
			ID:          devUserID(email),
			Email:       email,
			DisplayName: "",
			CreatedAt:   now,
			UpdatedAt:   now,
			LastLoginAt: &now,
		}, nil
	}

	return s.repo.UpsertDevUser(ctx, email)
}

func (s *Service) ResolveOrCreateWalletUser(ctx context.Context, address string) (*User, error) {
	address = normalizeWalletAddress(address)
	if address == "" {
		return nil, fmt.Errorf("empty wallet address")
	}

	id := walletUserID(address)
	email := walletUserEmail(address)
	displayName := address

	if s == nil || s.repo == nil {
		now := time.Now().UTC()
		return &User{
			ID:          id,
			Email:       email,
			DisplayName: displayName,
			CreatedAt:   now,
			UpdatedAt:   now,
			LastLoginAt: &now,
		}, nil
	}

	return s.repo.UpsertWalletUser(ctx, id, email, displayName)
}

func (s *Service) GetByID(ctx context.Context, id string, emailHint string) (*User, error) {
	id = strings.TrimSpace(id)
	emailHint = normalizeEmail(emailHint)

	if id == "" {
		return nil, fmt.Errorf("empty user id")
	}

	if s == nil || s.repo == nil {
		now := time.Now().UTC()
		return &User{
			ID:          id,
			Email:       emailHint,
			DisplayName: "",
			CreatedAt:   now,
			UpdatedAt:   now,
		}, nil
	}

	return s.repo.GetByID(ctx, id)
}

func (s *Service) UpdateDisplayName(ctx context.Context, id, displayName string) (*User, error) {
	id = strings.TrimSpace(id)
	displayName = normalizeDisplayName(displayName)

	if id == "" {
		return nil, fmt.Errorf("empty user id")
	}
	if displayName == "" {
		return nil, fmt.Errorf("empty display name")
	}
	if len(displayName) > 120 {
		return nil, fmt.Errorf("display name too long")
	}

	if s == nil || s.repo == nil {
		now := time.Now().UTC()
		return &User{
			ID:          id,
			DisplayName: displayName,
			CreatedAt:   now,
			UpdatedAt:   now,
		}, nil
	}

	return s.repo.UpdateDisplayName(ctx, id, displayName)
}

func normalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

func normalizeWalletAddress(address string) string {
	return strings.TrimSpace(strings.ToLower(address))
}

func normalizeDisplayName(displayName string) string {
	return strings.TrimSpace(displayName)
}

func devUserID(email string) string {
	id := strings.ReplaceAll(email, "@", "_")
	id = strings.ReplaceAll(id, ".", "_")
	id = strings.ReplaceAll(id, "+", "_")
	id = strings.ReplaceAll(id, "-", "_")
	return "u_" + id
}

func walletUserID(address string) string {
	address = normalizeWalletAddress(address)
	address = strings.TrimPrefix(address, "0x")
	return "u_wallet_" + address
}

func walletUserEmail(address string) string {
	address = normalizeWalletAddress(address)
	address = strings.TrimPrefix(address, "0x")
	return "wallet." + address + "@wallet.scavo.local"
}
