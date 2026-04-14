package httpx

import (
	"context"
	"net/http"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
	coreauthorization "github.com/e-scavo/scavo-exchange-backend/internal/core/authorization"
)

func AuthorizationSubjectFromContext(ctx context.Context) (*coreauthorization.AuthorizationSubject, bool) {
	return coreauthorization.SubjectFromContext(ctx)
}

func HydrateAuthorization() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := coreauth.ClaimsFromContext(r.Context())
			if !ok || claims == nil {
				next.ServeHTTP(w, r)
				return
			}

			subject, ok := coreauthorization.SubjectFromClaims(claims)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			ctx := coreauthorization.WithSubject(r.Context(), subject)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
