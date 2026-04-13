package httpx

import (
	"context"
	"net/http"
	"strings"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
	coreerrs "github.com/e-scavo/scavo-exchange-backend/internal/core/errs"
)

func AuthClaimsFromContext(ctx context.Context) (*coreauth.Claims, bool) {
	return coreauth.ClaimsFromContext(ctx)
}

func RequireAuth(tokens *coreauth.TokenService, allowQueryToken bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if tokens == nil {
				WriteError(w, http.StatusInternalServerError, coreerrs.NewResponseError("AUTH_NOT_CONFIGURED", "authentication is not configured", nil))
				return
			}

			token := coreauth.ExtractTokenFromRequest(r, allowQueryToken)
			if strings.TrimSpace(token) == "" {
				WriteError(w, http.StatusUnauthorized, coreerrs.NewResponseError("AUTH_MISSING_BEARER_TOKEN", "missing bearer token", nil))
				return
			}

			claims, err := tokens.Parse(token)
			if err != nil || claims == nil || strings.TrimSpace(claims.UserID) == "" {
				WriteError(w, http.StatusUnauthorized, coreerrs.NewResponseError("AUTH_UNAUTHORIZED", "authentication required", nil))
				return
			}

			ctx := context.WithValue(r.Context(), coreauth.ClaimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
