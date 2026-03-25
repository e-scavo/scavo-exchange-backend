package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type WalletChallengeRequest struct {
	Address string `json:"address"`
	Chain   string `json:"chain,omitempty"`
}

type WalletChallengeResponse struct {
	Challenge *WalletChallenge `json:"challenge"`
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
