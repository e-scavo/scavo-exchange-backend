package auth

import (
	"net/http"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
)

type WalletsResponse struct {
	Wallets []*WalletIdentity `json:"wallets"`
}

func (h HTTPHandlers) Wallets(w http.ResponseWriter, r *http.Request) {
	claims, ok := coreauth.ClaimsFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		return
	}

	if h.WalletIdentities == nil {
		writeJSON(w, http.StatusOK, WalletsResponse{Wallets: []*WalletIdentity{}})
		return
	}

	wallets, err := h.WalletIdentities.ListByUser(r.Context(), claims.UserID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "wallet_identity_error"})
		return
	}
	if wallets == nil {
		wallets = []*WalletIdentity{}
	}

	writeJSON(w, http.StatusOK, WalletsResponse{Wallets: wallets})
}
