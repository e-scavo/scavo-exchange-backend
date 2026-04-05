package auth

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
	usermod "github.com/e-scavo/scavo-exchange-backend/internal/modules/user"
)

type stubUserRepo struct {
	upsertResult        *usermod.User
	upsertErr           error
	getByIDFn           func(ctx context.Context, id string) (*usermod.User, error)
	updateDisplayNameFn func(ctx context.Context, id, displayName string) (*usermod.User, error)
}

func (s *stubUserRepo) UpdateDisplayName(ctx context.Context, id, displayName string) (*usermod.User, error) {
	if s.updateDisplayNameFn != nil {
		return s.updateDisplayNameFn(ctx, id, displayName)
	}

	if s.getByIDFn != nil {
		u, err := s.getByIDFn(ctx, id)
		if err != nil {
			return nil, err
		}
		if u != nil {
			u.DisplayName = displayName
			return u, nil
		}
	}

	return &usermod.User{
		ID:          id,
		DisplayName: displayName,
	}, nil
}

func (s *stubUserRepo) UpsertDevUser(ctx context.Context, email string) (*usermod.User, error) {
	return s.upsertResult, s.upsertErr
}

func (s *stubUserRepo) UpsertWalletUser(ctx context.Context, id, email, displayName string) (*usermod.User, error) {
	return &usermod.User{
		ID:          id,
		Email:       email,
		DisplayName: displayName,
	}, nil
}

func (s *stubUserRepo) GetByID(ctx context.Context, id string) (*usermod.User, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, id)
	}
	return nil, usermod.ErrUserNotFound
}

func newTokenServiceForTest(t *testing.T) *coreauth.TokenService {
	t.Helper()

	ts, err := coreauth.NewTokenService("dev_dev_dev_dev_dev_dev_dev_dev", "scavo-exchange-backend", time.Hour)
	if err != nil {
		t.Fatalf("NewTokenService error: %v", err)
	}

	return ts
}

func TestService_LoginDev_Success(t *testing.T) {
	repo := &stubUserRepo{
		upsertResult: &usermod.User{
			ID:    "u_persisted",
			Email: "dev@example.com",
		},
	}

	svc := NewService(newTokenServiceForTest(t), usermod.NewService(repo), time.Hour)

	result, err := svc.LoginDev(context.Background(), " Dev@Example.com ", "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected login result")
	}

	if result.User == nil || result.User.ID != "u_persisted" {
		t.Fatalf("unexpected user: %#v", result.User)
	}

	if result.TokenType != "Bearer" {
		t.Fatalf("unexpected token type: %q", result.TokenType)
	}

	if result.ExpiresIn != int64(time.Hour.Seconds()) {
		t.Fatalf("unexpected expires_in: %d", result.ExpiresIn)
	}
}

func TestService_LoginDev_InvalidCredentials(t *testing.T) {
	svc := NewService(newTokenServiceForTest(t), usermod.NewService(nil), time.Hour)

	_, err := svc.LoginDev(context.Background(), "", "bad")
	if err == nil {
		t.Fatal("expected invalid credentials error")
	}

	if err != ErrInvalidCredentials {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestService_ResolveCurrentUser_Success(t *testing.T) {
	repo := &stubUserRepo{
		getByIDFn: func(ctx context.Context, id string) (*usermod.User, error) {
			if id != "u_persisted" {
				t.Fatalf("unexpected user id: %q", id)
			}

			return &usermod.User{
				ID:    id,
				Email: "dev@example.com",
			}, nil
		},
	}

	users := usermod.NewService(repo)
	svc := NewService(newTokenServiceForTest(t), users, time.Hour)

	token, err := svc.tokens.Mint("u_persisted", "dev@example.com")
	if err != nil {
		t.Fatalf("Mint error: %v", err)
	}

	user, err := svc.ResolveCurrentUser(context.Background(), token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user == nil || user.ID != "u_persisted" {
		t.Fatalf("unexpected user: %#v", user)
	}
}

func TestService_ResolveCurrentUser_UserNotFound(t *testing.T) {
	repo := &stubUserRepo{}
	svc := NewService(newTokenServiceForTest(t), usermod.NewService(repo), time.Hour)

	token, err := svc.tokens.Mint("u_missing", "missing@example.com")
	if err != nil {
		t.Fatalf("Mint error: %v", err)
	}

	_, err = svc.ResolveCurrentUser(context.Background(), token)
	if err == nil {
		t.Fatal("expected unauthorized error")
	}

	if err != ErrUnauthorized {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestService_ResolveSessionClaims_Success(t *testing.T) {
	repo := &stubUserRepo{
		getByIDFn: func(ctx context.Context, id string) (*usermod.User, error) {
			return &usermod.User{
				ID:    id,
				Email: "dev@example.com",
			}, nil
		},
	}

	svc := NewService(newTokenServiceForTest(t), usermod.NewService(repo), time.Hour)

	now := time.Now().UTC()
	claims := &coreauth.Claims{
		UserID: "u_persisted",
		Email:  "dev@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "scavo-exchange-backend",
			Subject:   "u_persisted",
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
		},
	}

	view, err := svc.ResolveSessionClaims(context.Background(), claims)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if view == nil {
		t.Fatal("expected session view")
	}
	if !view.Authenticated {
		t.Fatal("expected authenticated session view")
	}
	if view.User == nil || view.User.ID != "u_persisted" {
		t.Fatalf("unexpected user in session view: %#v", view.User)
	}
	if view.ExpiresAt == nil {
		t.Fatal("expected expires_at in session view")
	}
}
