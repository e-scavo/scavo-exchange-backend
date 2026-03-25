package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHTTPHandlers_WalletChallenge_Success(t *testing.T) {
	h := HTTPHandlers{
		PublicBaseURL: "https://api.scavo.exchange",
		ChallengeTTL:  5 * time.Minute,
		Challenges:    NewInMemoryWalletChallengeStore(),
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/wallet/challenge",
		strings.NewReader(`{"address":"0x1111111111111111111111111111111111111111","chain":"scavium"}`),
	)
	rec := httptest.NewRecorder()

	h.WalletChallenge(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}

	var payload WalletChallengeResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if payload.Challenge == nil {
		t.Fatal("expected challenge payload")
	}
	if payload.Challenge.ID == "" {
		t.Fatal("expected challenge id")
	}
	if payload.Challenge.Nonce == "" {
		t.Fatal("expected challenge nonce")
	}
}

func TestHTTPHandlers_WalletChallenge_InvalidAddress(t *testing.T) {
	h := HTTPHandlers{
		PublicBaseURL: "https://api.scavo.exchange",
		ChallengeTTL:  5 * time.Minute,
		Challenges:    NewInMemoryWalletChallengeStore(),
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/wallet/challenge",
		strings.NewReader(`{"address":"bad-address","chain":"scavium"}`),
	)
	rec := httptest.NewRecorder()

	h.WalletChallenge(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}
}
