package auth

import (
	"strings"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
)

type authenticatedContextView struct {
	UserID           string
	Email            string
	WalletID         string
	WalletAddress    string
	AuthMethod       string
	Chain            string
	HasWalletSession bool
}

func buildAuthenticatedContextView(claims *coreauth.Claims) authenticatedContextView {
	if claims == nil {
		return authenticatedContextView{AuthMethod: "password_dev"}
	}

	view := authenticatedContextView{
		UserID:        strings.TrimSpace(claims.UserID),
		Email:         normalizeEmail(claims.Email),
		WalletID:      strings.TrimSpace(claims.WalletID),
		WalletAddress: normalizeWalletAddress(claims.WalletAddress),
		AuthMethod:    strings.TrimSpace(claims.AuthMethod),
		Chain:         normalizeChain(claims.Chain),
	}

	if view.AuthMethod == "" {
		view.AuthMethod = "password_dev"
	}
	if view.WalletAddress == "" {
		view.WalletID = ""
		view.Chain = ""
	}
	view.HasWalletSession = view.WalletAddress != ""

	return view
}
