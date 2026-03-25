package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
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
	AccessToken string
	TokenType   string
	ExpiresIn   int64
	User        *usermod.User
}

type SessionView struct {
	Authenticated bool          `json:"authenticated"`
	TokenType     string        `json:"token_type"`
	UserID        string        `json:"user_id"`
	Email         string        `json:"email,omitempty"`
	Subject       string        `json:"subject,omitempty"`
	Issuer        string        `json:"issuer,omitempty"`
	ExpiresAt     *time.Time    `json:"expires_at,omitempty"`
	User          *usermod.User `json:"user,omitempty"`
}

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
			return nil, ErrUnauthorized
		}
		return nil, err
	}

	return user, nil
}

func (s *Service) ResolveSession(ctx context.Context, token string) (*SessionView, error) {
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

	return s.ResolveSessionClaims(ctx, claims)
}

func (s *Service) ResolveSessionClaims(ctx context.Context, claims *coreauth.Claims) (*SessionView, error) {
	user, err := s.ResolveCurrentUserClaims(ctx, claims)
	if err != nil {
		return nil, err
	}

	var expiresAt *time.Time
	if claims != nil && claims.ExpiresAt != nil {
		ts := claims.ExpiresAt.Time.UTC()
		expiresAt = &ts
	}

	view := &SessionView{
		Authenticated: true,
		TokenType:     "Bearer",
		UserID:        strings.TrimSpace(claims.UserID),
		Email:         normalizeEmail(claims.Email),
		Subject:       strings.TrimSpace(claims.Subject),
		Issuer:        strings.TrimSpace(claims.Issuer),
		ExpiresAt:     expiresAt,
		User:          user,
	}

	if view.Subject == "" {
		view.Subject = view.UserID
	}

	return view, nil
}

func normalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}
