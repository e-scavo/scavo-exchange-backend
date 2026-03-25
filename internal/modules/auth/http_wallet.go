package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	usermod "github.com/e-scavo/scavo-exchange-backend/internal/modules/user"
)

type WalletChallengeRequest struct {
	Address string `json:"address"`
	Chain   string `json:"chain,omitempty"`
}

type WalletChallengeResponse struct {
	Challenge *WalletChallenge `json:"challenge"`
}

type WalletVerifyRequest struct {
	ChallengeID string `json:"challenge_id"`
	Address     string `json:"address"`
	Signature   string `json:"signature"`
}

type WalletVerifyResponse struct {
	AccessToken   string           `json:"access_token"`
	TokenType     string           `json:"token_type"`
	ExpiresIn     int64            `json:"expires_in"`
	UserID        string           `json:"user_id"`
	WalletAddress string           `json:"wallet_address"`
	Chain         string           `json:"chain"`
	AuthMethod    string           `json:"auth_method"`
	User          *usermod.User    `json:"user,omitempty"`
	Challenge     *WalletChallenge `json:"challenge,omitempty"`
}

func (h HTTPHandlers) WalletChallenge(w http.ResponseWriter, r *http.Request) {
	var req WalletChallengeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "bad_request"})
		return
	}

	challengeTTL := h.ChallengeTTL
	if challengeTTL <= 0 {
		challengeTTL = 5 * time.Minute
	}

	svc := NewWalletChallengeService(h.Challenges, h.PublicBaseURL, challengeTTL)
	challenge, err := svc.Create(r.Context(), req.Address, req.Chain)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidWalletAddress):
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_wallet_address"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "wallet_challenge_error"})
		}
		return
	}

	writeJSON(w, http.StatusOK, WalletChallengeResponse{Challenge: challenge})
}

func (h HTTPHandlers) WalletVerify(w http.ResponseWriter, r *http.Request) {
	var req WalletVerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "bad_request"})
		return
	}

	challengeTTL := h.ChallengeTTL
	if challengeTTL <= 0 {
		challengeTTL = 5 * time.Minute
	}

	challengeSvc := NewWalletChallengeService(h.Challenges, h.PublicBaseURL, challengeTTL)
	loginSvc := NewService(h.Tokens, h.Users, h.TTL)
	verifySvc := NewWalletVerificationService(challengeSvc, loginSvc)

	result, challenge, err := verifySvc.VerifyAndLogin(r.Context(), req.ChallengeID, req.Address, req.Signature)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidWalletAddress):
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_wallet_address"})
		case errors.Is(err, ErrInvalidWalletSignature):
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "invalid_wallet_signature"})
		case errors.Is(err, ErrWalletChallengeNotFound):
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "wallet_challenge_not_found"})
		case errors.Is(err, ErrChallengeExpired):
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "wallet_challenge_expired"})
		case errors.Is(err, ErrChallengeUsed):
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "wallet_challenge_used"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "wallet_verify_error"})
		}
		return
	}

	userID := ""
	var user *usermod.User
	if result != nil && result.User != nil {
		userID = result.User.ID
		user = result.User
	}

	writeJSON(w, http.StatusOK, WalletVerifyResponse{
		AccessToken:   result.AccessToken,
		TokenType:     result.TokenType,
		ExpiresIn:     result.ExpiresIn,
		UserID:        userID,
		WalletAddress: result.WalletAddress,
		Chain:         result.Chain,
		AuthMethod:    result.AuthMethod,
		User:          user,
		Challenge:     challenge,
	})
}
