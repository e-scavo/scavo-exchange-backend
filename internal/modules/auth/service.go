package auth

import (
	"context"
	"errors"
	"fmt"
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

// SessionView already has a stable canonical shape in auth/app.
// For root compatibility, keep it as an alias here.
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
	if s == nil || s.tokens == nil {
		return nil, fmt.Errorf("token service not configured")
	}

	email = normalizeEmail(email)
	if email == "" || strings.TrimSpace(password) != "dev" {
		return nil, ErrInvalidCredentials
	}

	userID := "u_" + strings.ReplaceAll(email, "@", "_")
	var user *usermod.User
	var err error

	if s.users != nil {
		user, err = s.users.ResolveOrCreateDevUser(ctx, email)
		if err != nil {
			return nil, err
		}
		userID = user.ID
	}

	token, err := s.tokens.Mint(userID, email)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(s.ttl.Seconds()),
		User:        user,
		AuthMethod:  "password_dev",
	}, nil
}

func (s *Service) LoginWallet(ctx context.Context, walletID, address, chain string) (*LoginResult, error) {
	return s.LoginWalletForUser(ctx, walletUser(address), walletID, address, chain)
}

func (s *Service) LoginWalletForUser(ctx context.Context, user *usermod.User, walletID, address, chain string) (*LoginResult, error) {
	if s == nil || s.tokens == nil {
		return nil, fmt.Errorf("token service not configured")
	}

	address = normalizeWalletAddress(address)
	if !evmAddressRE.MatchString(address) {
		return nil, ErrInvalidWalletAddress
	}
	chain = normalizeChain(chain)
	walletID = strings.TrimSpace(walletID)

	if user == nil {
		user = walletUser(address)
	}

	token, err := s.tokens.MintWithOptions(coreauth.MintOptions{
		UserID:        strings.TrimSpace(user.ID),
		Email:         normalizeEmail(user.Email),
		WalletID:      walletID,
		WalletAddress: address,
		AuthMethod:    "wallet_evm",
		Chain:         chain,
		Subject:       address,
	})
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:   token,
		TokenType:     "Bearer",
		ExpiresIn:     int64(s.ttl.Seconds()),
		User:          user,
		WalletID:      walletID,
		WalletAddress: address,
		Chain:         chain,
		AuthMethod:    "wallet_evm",
	}, nil
}

func (s *Service) ResolveCurrentUser(ctx context.Context, token string) (*usermod.User, error) {
	if s == nil || s.tokens == nil {
		return nil, ErrUnauthorized
	}

	token = strings.TrimSpace(token)
	if token == "" {
		return nil, ErrUnauthorized
	}

	claims, err := s.tokens.Parse(token)
	if err != nil || claims == nil || strings.TrimSpace(claims.UserID) == "" {
		return nil, ErrUnauthorized
	}

	return s.ResolveCurrentUserClaims(ctx, claims)
}

func (s *Service) ResolveCurrentUserClaims(ctx context.Context, claims *coreauth.Claims) (*usermod.User, error) {
	if claims == nil || strings.TrimSpace(claims.UserID) == "" {
		return nil, ErrUnauthorized
	}

	if s.users == nil {
		if strings.TrimSpace(claims.WalletAddress) != "" {
			return walletUser(claims.WalletAddress), nil
		}

		now := time.Now().UTC()
		return &usermod.User{
			ID:        claims.UserID,
			Email:     normalizeEmail(claims.Email),
			CreatedAt: now,
			UpdatedAt: now,
		}, nil
	}

	user, err := s.users.GetByID(ctx, claims.UserID, claims.Email)
	if err != nil {
		if errors.Is(err, usermod.ErrUserNotFound) {
			if strings.TrimSpace(claims.WalletAddress) != "" {
				return walletUser(claims.WalletAddress), nil
			}
			return nil, ErrUnauthorized
		}
		return nil, err
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

// Keep this helper in root for now because service.go still owns real runtime
// behavior for login/current-user paths in 0.11.4C3.1b1.
func buildSessionViewWithUser(claims *coreauth.Claims, user *usermod.User) *SessionView {
	return authapp.BuildSessionViewWithUser(claims, user)
}

// Root compatibility helper still needed by auth_context.go and current service
// flows. This remains intentionally local in 0.11.4C3.1b1.
func normalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

// Root compatibility wallet user helpers remain local for now because
// ResolveCurrentUserClaims and wallet login flows still depend on them directly.
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
