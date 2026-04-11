package auth

import (
	"context"
	"errors"
	"net/http"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
	usermod "github.com/e-scavo/scavo-exchange-backend/internal/modules/user"
	usersettingsmod "github.com/e-scavo/scavo-exchange-backend/internal/modules/usersettings"
)

type BootstrapWalletsView struct {
	Items []*WalletReadModel `json:"items"`
	Total int                `json:"total"`
}

func buildBootstrapWalletsView(wallets []*WalletReadModel) BootstrapWalletsView {
	if wallets == nil {
		wallets = []*WalletReadModel{}
	}

	return BootstrapWalletsView{
		Items: wallets,
		Total: len(wallets),
	}
}

type BootstrapResponse struct {
	Session  *SessionView         `json:"session"`
	User     *usermod.User        `json:"user,omitempty"`
	Profile  *ProfileView         `json:"profile,omitempty"`
	Settings usersettingsmod.View `json:"settings"`
	Wallets  BootstrapWalletsView `json:"wallets"`
}

func listWalletReadModelsForUser(ctx context.Context, userID string, store WalletIdentityStore) ([]*WalletReadModel, error) {
	if store == nil || userID == "" {
		return []*WalletReadModel{}, nil
	}

	wallets, err := store.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if wallets == nil {
		wallets = []*WalletIdentity{}
	}

	return mapWalletIdentitiesToReadModels(wallets), nil
}

func (h HTTPHandlers) Bootstrap(w http.ResponseWriter, r *http.Request) {
	claims, ok := coreauth.ClaimsFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		return
	}

	payload, err := h.Application().GetBootstrap(r.Context(), claims)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnauthorized), errors.Is(err, usersettingsmod.ErrUserIDRequired):
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		case errors.Is(err, ErrWalletIdentityStore):
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "wallet_identity_error"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "auth_service_error"})
		}
		return
	}

	writeJSON(w, http.StatusOK, payload)
}
