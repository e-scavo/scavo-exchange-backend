package auth

import (
	"errors"
	"net/http"

	coreerrs "github.com/e-scavo/scavo-exchange-backend/internal/core/errs"
	authapp "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth/app"
	usersettingsmod "github.com/e-scavo/scavo-exchange-backend/internal/modules/usersettings"
)

type BootstrapWalletsView = authapp.BootstrapWalletsView
type BootstrapResponse = authapp.BootstrapResponse

func (h HTTPHandlers) Bootstrap(w http.ResponseWriter, r *http.Request) {
	claims, ok := h.requireClaims(w, r)
	if !ok {
		return
	}

	payload, err := h.Application().GetBootstrap(r.Context(), claims)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnauthorized), errors.Is(err, usersettingsmod.ErrUserIDRequired):
			writeAppErrorJSON(w, coreerrs.AuthUnauthorized())
		case errors.Is(err, ErrWalletIdentityStore):
			writeAppErrorJSON(w, coreerrs.WalletIdentityError(nil))
		default:
			writeAppErrorJSON(w, coreerrs.AuthServiceError(nil))
		}
		return
	}

	writeJSON(w, http.StatusOK, payload)
}
