package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
	authapp "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth/app"
	usermod "github.com/e-scavo/scavo-exchange-backend/internal/modules/user"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
)

type Service struct {
	tokens *coreauth.TokenService
	users  *usermod.Service
	ttl    time.Duration
}

type LoginResult struct {
	AccessToken   string
	TokenType     string
	ExpiresIn     int64
	User          *usermod.User
	WalletID      string
	WalletAddress string
	Chain         string
	AuthMethod    string
}

type SessionView = authapp.SessionView

func NewService(tokens *coreauth.TokenService, users *usermod.Service, ttl time.Duration) *Service {
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}

	return &Service{
		tokens: tokens,
		users:  users,
		ttl:    ttl,
	}
}

func (s *Service) appService() *authapp.Service {
	if s == nil {
		return authapp.NewService(nil, nil, 0)
	}
	return authapp.NewService(s.tokens, s.users, s.ttl)
}

func (s *Service) LoginDev(ctx context.Context, email, password string) (*LoginResult, error) {
	result, err := s.appService().LoginDev(ctx, email, password)
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	return mapLoginResultFromApp(result), nil
}

func (s *Service) LoginWallet(ctx context.Context, walletID, address, chain string) (*LoginResult, error) {
	result, err := s.appService().LoginWallet(ctx, walletID, address, chain)
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	return mapLoginResultFromApp(result), nil
}

func (s *Service) LoginWalletForUser(ctx context.Context, user *usermod.User, walletID, address, chain string) (*LoginResult, error) {
	result, err := s.appService().LoginWalletForUser(ctx, user, walletID, address, chain)
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	return mapLoginResultFromApp(result), nil
}

func (s *Service) ResolveCurrentUser(ctx context.Context, token string) (*usermod.User, error) {
	user, err := s.appService().ResolveCurrentUser(ctx, token)
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	return user, nil
}

func (s *Service) ResolveCurrentUserClaims(ctx context.Context, claims *coreauth.Claims) (*usermod.User, error) {
	user, err := s.appService().ResolveCurrentUserClaims(ctx, claims)
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	return user, nil
}

func (s *Service) ResolveSession(ctx context.Context, token string) (*SessionView, error) {
	session, err := s.appService().ResolveSession(ctx, token)
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	return session, nil
}

func (s *Service) ResolveSessionClaims(ctx context.Context, claims *coreauth.Claims) (*SessionView, error) {
	session, err := s.appService().ResolveSessionClaims(ctx, claims)
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	return session, nil
}

// Kept in root as compatibility helper for remaining package-level consumers.
func buildSessionViewWithUser(claims *coreauth.Claims, user *usermod.User) *SessionView {
	return authapp.BuildSessionViewWithUser(claims, user)
}

// Kept in root as compatibility helper for remaining package-level consumers.
func normalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

// Kept in root as compatibility helper for remaining package-level consumers.
// walletUserID and walletUserEmail already exist elsewhere in the auth package.
func walletUser(address string) *usermod.User {
	now := time.Now().UTC()
	address = normalizeWalletAddress(address)
	return &usermod.User{
		ID:          walletUserID(address),
		Email:       walletUserEmail(address),
		DisplayName: address,
		CreatedAt:   now,
		UpdatedAt:   now,
		LastLoginAt: &now,
	}
}

func normalizeServiceError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, authapp.ErrInvalidCredentials):
		return ErrInvalidCredentials
	case errors.Is(err, authapp.ErrUnauthorized):
		return ErrUnauthorized
	case errors.Is(err, authapp.ErrInvalidWalletAddress):
		return ErrInvalidWalletAddress
	default:
		return err
	}
}

func mapLoginResultFromApp(result *authapp.LoginResult) *LoginResult {
	if result == nil {
		return nil
	}

	return &LoginResult{
		AccessToken:   result.AccessToken,
		TokenType:     result.TokenType,
		ExpiresIn:     result.ExpiresIn,
		User:          result.User,
		WalletID:      result.WalletID,
		WalletAddress: result.WalletAddress,
		Chain:         result.Chain,
		AuthMethod:    result.AuthMethod,
	}
}
