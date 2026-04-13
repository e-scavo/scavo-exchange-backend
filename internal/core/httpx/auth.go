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
				WriteAppError(w, coreerrs.AuthNotConfigured())
				return
			}

			token := coreauth.ExtractTokenFromRequest(r, allowQueryToken)
			if strings.TrimSpace(token) == "" {
				WriteAppError(w, coreerrs.AuthMissingBearerToken())
				return
			}

			claims, err := tokens.Parse(token)
			if err != nil || claims == nil || strings.TrimSpace(claims.UserID) == "" {
				WriteAppError(w, coreerrs.AuthUnauthorized())
				return
			}

			ctx := context.WithValue(r.Context(), coreauth.ClaimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
