package auth

import (
	"net/http"
	"time"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
)

type WalletReadModel struct {
	ID         string     `json:"id"`
	Address    string     `json:"address"`
	UserID     string     `json:"user_id,omitempty"`
	LinkedAt   *time.Time `json:"linked_at,omitempty"`
	DetachedAt *time.Time `json:"detached_at,omitempty"`
	IsPrimary  bool       `json:"is_primary"`
	Status     string     `json:"status"`
}

type WalletsResponse struct {
	Wallets []*WalletReadModel `json:"wallets"`
}

func mapWalletIdentityToReadModel(wallet *WalletIdentity) *WalletReadModel {
	if wallet == nil {
		return nil
	}

	status := "unlinked"
	switch {
	case wallet.UserID != "":
		status = "active"
	case wallet.DetachedAt != nil:
		status = "detached"
	}

	return &WalletReadModel{
		ID:         wallet.ID,
		Address:    wallet.Address,
		UserID:     wallet.UserID,
		LinkedAt:   wallet.LinkedAt,
		DetachedAt: wallet.DetachedAt,
		IsPrimary:  wallet.IsPrimary,
		Status:     status,
	}
}

func mapWalletIdentitiesToReadModels(wallets []*WalletIdentity) []*WalletReadModel {
	if len(wallets) == 0 {
		return []*WalletReadModel{}
	}

	out := make([]*WalletReadModel, 0, len(wallets))
	for _, wallet := range wallets {
		mapped := mapWalletIdentityToReadModel(wallet)
		if mapped != nil {
			out = append(out, mapped)
		}
	}

	if out == nil {
		return []*WalletReadModel{}
	}

	return out
}

func (h HTTPHandlers) Wallets(w http.ResponseWriter, r *http.Request) {
	claims, ok := coreauth.ClaimsFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		return
	}

	if h.WalletIdentities == nil {
		writeJSON(w, http.StatusOK, WalletsResponse{Wallets: []*WalletReadModel{}})
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

	writeJSON(w, http.StatusOK, WalletsResponse{Wallets: mapWalletIdentitiesToReadModels(wallets)})
}
