package httpx

import (
	"context"
	"net/http"
	"strings"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
)

func AuthClaimsFromContext(ctx context.Context) (*coreauth.Claims, bool) {
	return coreauth.ClaimsFromContext(ctx)
}

func RequireAuth(tokens *coreauth.TokenService, allowQueryToken bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if tokens == nil {
				WriteJSON(w, http.StatusInternalServerError, map[string]any{"error": "auth_not_configured"})
				return
			}

			token := coreauth.ExtractTokenFromRequest(r, allowQueryToken)
			if strings.TrimSpace(token) == "" {
				WriteJSON(w, http.StatusUnauthorized, map[string]any{"error": "missing_bearer_token"})
				return
			}

			claims, err := tokens.Parse(token)
			if err != nil || claims == nil || strings.TrimSpace(claims.UserID) == "" {
				WriteJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
				return
			}

			ctx := context.WithValue(r.Context(), coreauth.ClaimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
