package auth

import "context"

type contextKey string

const ClaimsContextKey contextKey = "auth_claims"

func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	if ctx == nil {
		return nil, false
	}

	claims, ok := ctx.Value(ClaimsContextKey).(*Claims)
	if !ok || claims == nil {
		return nil, false
	}

	return claims, true
}
