package httpx

import (
	"net/http"

	coreauthorization "github.com/e-scavo/scavo-exchange-backend/internal/core/authorization"
	coreerrs "github.com/e-scavo/scavo-exchange-backend/internal/core/errs"
)

func RequirePermission(action coreauthorization.Action, resource coreauthorization.Resource) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			subject, ok := AuthorizationSubjectFromContext(r.Context())
			if !ok || subject == nil {
				WriteAppError(w, coreerrs.AuthForbidden(map[string]any{
					"action":   action,
					"resource": resource,
				}))
				return
			}

			decision := coreauthorization.NewPolicyEvaluator().Evaluate(*subject, action, resource)
			if !decision.Allowed {
				details := map[string]any{
					"action":   decision.Action,
					"resource": decision.Resource,
				}
				if decision.Permission != "" {
					details["permission"] = decision.Permission
				}
				WriteAppError(w, coreerrs.AuthForbidden(details))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
