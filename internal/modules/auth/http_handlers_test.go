package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
	usermod "github.com/e-scavo/scavo-exchange-backend/internal/modules/user"
)

func mustTokenService(t *testing.T) *coreauth.TokenService {
	t.Helper()

	ts, err := coreauth.NewTokenService("dev_dev_dev_dev_dev_dev_dev_dev", "scavo-exchange-backend", time.Hour)
	if err != nil {
		t.Fatalf("NewTokenService error: %v", err)
	}

	return ts
}

func TestHTTPHandlers_Login_Success(t *testing.T) {
	h := HTTPHandlers{
		Tokens: mustTokenService(t),
		TTL:    time.Hour,
		Users:  usermod.NewService(nil),
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"email":"test@example.com","password":"dev"}`))
	rec := httptest.NewRecorder()

	h.Login(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}

	var payload LoginResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if payload.AccessToken == "" {
		t.Fatal("expected access token")
	}

	if payload.UserID != "u_test_example_com" {
		t.Fatalf("unexpected user id: %q", payload.UserID)
	}
}

func TestHTTPHandlers_Login_InvalidBody(t *testing.T) {
	h := HTTPHandlers{
		Tokens: mustTokenService(t),
		TTL:    time.Hour,
		Users:  usermod.NewService(nil),
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"email":`))
	rec := httptest.NewRecorder()

	h.Login(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
}

func sessionClaims() *coreauth.Claims {
	now := time.Now().UTC()

	return &coreauth.Claims{
		UserID: "u_test_example_com",
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "scavo-exchange-backend",
			Subject:   "u_test_example_com",
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
		},
	}
}

func TestHTTPHandlers_Me_Success(t *testing.T) {
	ts := mustTokenService(t)

	h := HTTPHandlers{
		Tokens: ts,
		TTL:    time.Hour,
		Users:  usermod.NewService(nil),
	}

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req = req.WithContext(context.WithValue(req.Context(), coreauth.ClaimsContextKey, sessionClaims()))
	rec := httptest.NewRecorder()

	h.Me(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}

	var payload MeResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if payload.User == nil || payload.User.ID != "u_test_example_com" {
		t.Fatalf("unexpected user payload: %#v", payload.User)
	}
}

func TestHTTPHandlers_Me_MissingClaims(t *testing.T) {
	h := HTTPHandlers{
		Tokens: mustTokenService(t),
		TTL:    time.Hour,
		Users:  usermod.NewService(nil),
	}

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	rec := httptest.NewRecorder()

	h.Me(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
}

func TestHTTPHandlers_Session_Success(t *testing.T) {
	h := HTTPHandlers{
		Tokens: mustTokenService(t),
		TTL:    time.Hour,
		Users:  usermod.NewService(nil),
	}

	req := httptest.NewRequest(http.MethodGet, "/auth/session", nil)
	req = req.WithContext(context.WithValue(req.Context(), coreauth.ClaimsContextKey, sessionClaims()))
	rec := httptest.NewRecorder()

	h.Session(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}

	var payload SessionResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if payload.Session == nil {
		t.Fatal("expected session payload")
	}
	if !payload.Session.Authenticated {
		t.Fatal("expected authenticated session")
	}
	if payload.Session.UserID != "u_test_example_com" {
		t.Fatalf("unexpected session user id: %q", payload.Session.UserID)
	}
	if payload.Session.User == nil || payload.Session.User.Email != "test@example.com" {
		t.Fatalf("unexpected session user payload: %#v", payload.Session.User)
	}
}

func TestHTTPHandlers_Wallets_Success(t *testing.T) {
	store := NewInMemoryWalletIdentityStore()
	address := testWalletAddress()

	identity, err := store.GetOrCreate(context.Background(), address)
	if err != nil {
		t.Fatalf("GetOrCreate error: %v", err)
	}

	_, err = store.AttachUser(context.Background(), identity.ID, "u_test_example_com", true)
	if err != nil {
		t.Fatalf("AttachUser error: %v", err)
	}

	h := HTTPHandlers{
		Tokens:           mustTokenService(t),
		TTL:              time.Hour,
		Users:            usermod.NewService(nil),
		WalletIdentities: store,
	}

	req := httptest.NewRequest(http.MethodGet, "/auth/wallets", nil)
	req = req.WithContext(context.WithValue(req.Context(), coreauth.ClaimsContextKey, sessionClaims()))
	rec := httptest.NewRecorder()

	h.Wallets(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}

	var payload WalletsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if len(payload.Wallets) != 1 {
		t.Fatalf("expected 1 wallet, got %d", len(payload.Wallets))
	}
	if payload.Wallets[0].UserID != "u_test_example_com" {
		t.Fatalf("unexpected wallet user id: %q", payload.Wallets[0].UserID)
	}
	if !payload.Wallets[0].IsPrimary {
		t.Fatal("expected primary wallet")
	}
}
