package httpx

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
	"github.com/e-scavo/scavo-exchange-backend/internal/core/config"
	"github.com/e-scavo/scavo-exchange-backend/internal/core/logger"
	authmod "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth"
)

type versionedErrorEnvelope struct {
	Error struct {
		Code    string         `json:"code"`
		Message string         `json:"message"`
		Details map[string]any `json:"details,omitempty"`
	} `json:"error"`
}

func newVersioningTestRouter(t *testing.T) http.Handler {
	t.Helper()

	tokens, err := coreauth.NewTokenService("test_secret_value_1234567890", "scavo-exchange-backend", 24*time.Hour)
	if err != nil {
		t.Fatalf("token service init error: %v", err)
	}

	return NewRouter(RouterParams{
		Log:            logger.New("test"),
		Config:         config.Config{Env: "test", Version: "test", Commit: "test", CORSAllowOrigins: []string{"*"}, JWTTTLHrs: 24},
		TokenService:   tokens,
		ChallengeStore: authmod.NewInMemoryWalletChallengeStore(),
		ChallengeTTL:   5 * time.Minute,
		PublicBaseURL:  "https://api.scavo.exchange",
	})
}

func decodeVersionedError(t *testing.T, rec *httptest.ResponseRecorder) versionedErrorEnvelope {
	t.Helper()

	var payload versionedErrorEnvelope
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode error: %v body=%s", err, rec.Body.String())
	}
	return payload
}

func TestRoutePath(t *testing.T) {
	if got := routePath("", "/auth/me"); got != "/auth/me" {
		t.Fatalf("unexpected route path for legacy route: got=%q want=%q", got, "/auth/me")
	}
	if got := routePath("/api/v1", "/auth/me"); got != "/api/v1/auth/me" {
		t.Fatalf("unexpected route path for canonical route: got=%q want=%q", got, "/api/v1/auth/me")
	}
}

func TestNewRouter_ProtectedLegacyAndCanonicalRoutesShareMissingBearerContract(t *testing.T) {
	router := newVersioningTestRouter(t)

	for _, path := range []string{"/auth/me", "/api/v1/auth/me"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("unexpected status for %s: got=%d want=%d body=%s", path, rec.Code, http.StatusUnauthorized, rec.Body.String())
		}

		payload := decodeVersionedError(t, rec)
		if payload.Error.Code != "AUTH_MISSING_BEARER_TOKEN" {
			t.Fatalf("unexpected error code for %s: got=%q want=%q", path, payload.Error.Code, "AUTH_MISSING_BEARER_TOKEN")
		}
	}
}

func TestNewRouter_ProtectedLegacyAndCanonicalRoutesShareUnauthorizedContract(t *testing.T) {
	router := newVersioningTestRouter(t)

	for _, path := range []string{"/auth/me", "/api/v1/auth/me"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("unexpected status for %s: got=%d want=%d body=%s", path, rec.Code, http.StatusUnauthorized, rec.Body.String())
		}

		payload := decodeVersionedError(t, rec)
		if payload.Error.Code != "AUTH_UNAUTHORIZED" {
			t.Fatalf("unexpected error code for %s: got=%q want=%q", path, payload.Error.Code, "AUTH_UNAUTHORIZED")
		}
	}
}

func TestNewRouter_WalletChallengeLegacyAndCanonicalRoutesExposeCompatibleSuccessShape(t *testing.T) {
	router := newVersioningTestRouter(t)

	for _, path := range []string{"/auth/wallet/challenge", "/api/v1/auth/wallet/challenge"} {
		body := bytes.NewBufferString(`{"address":"0x1111111111111111111111111111111111111111","chain":"scavium"}`)
		req := httptest.NewRequest(http.MethodPost, path, body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("unexpected status for %s: got=%d want=%d body=%s", path, rec.Code, http.StatusOK, rec.Body.String())
		}

		var payload authmod.WalletChallengeResponse
		if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
			t.Fatalf("decode error for %s: %v body=%s", path, err, rec.Body.String())
		}
		if payload.Challenge == nil {
			t.Fatalf("nil challenge payload for %s", path)
		}
		if payload.Challenge.Address != "0x1111111111111111111111111111111111111111" {
			t.Fatalf("unexpected address for %s: got=%q", path, payload.Challenge.Address)
		}
		if payload.Challenge.Chain != "scavium" {
			t.Fatalf("unexpected chain for %s: got=%q", path, payload.Challenge.Chain)
		}
		if payload.Challenge.ID == "" || payload.Challenge.Nonce == "" || payload.Challenge.Message == "" {
			t.Fatalf("incomplete challenge payload for %s: %+v", path, payload.Challenge)
		}
	}
}
