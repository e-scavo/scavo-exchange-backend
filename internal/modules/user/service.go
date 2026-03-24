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

func normalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

func devUserID(email string) string {
	id := strings.ReplaceAll(email, "@", "_")
	id = strings.ReplaceAll(id, ".", "_")
	id = strings.ReplaceAll(id, "+", "_")
	id = strings.ReplaceAll(id, "-", "_")
	return "u_" + id
}
