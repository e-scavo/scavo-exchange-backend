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
	if h.UserSettings == nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "auth_service_error"})
		return
	}

	svc := NewService(h.Tokens, h.Users, h.TTL)
	user, err := svc.ResolveCurrentUserClaims(r.Context(), claims)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnauthorized):
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "auth_service_error"})
		}
		return
	}

	session := buildSessionViewWithUser(claims, user)
	profile, err := buildProfileViewWithUser(r.Context(), claims, user, h.WalletIdentities)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "auth_service_error"})
		return
	}

	settings, err := h.UserSettings.GetOrDefault(r.Context(), session.UserID)
	if err != nil {
		switch {
		case errors.Is(err, usersettingsmod.ErrUserIDRequired):
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "auth_service_error"})
		}
		return
	}

	wallets, err := listWalletReadModelsForUser(r.Context(), session.UserID, h.WalletIdentities)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "wallet_identity_error"})
		return
	}

	writeJSON(w, http.StatusOK, BootstrapResponse{
		Session:  session,
		User:     user,
		Profile:  profile,
		Settings: usersettingsmod.ToView(settings),
		Wallets:  buildBootstrapWalletsView(wallets),
	})
}
