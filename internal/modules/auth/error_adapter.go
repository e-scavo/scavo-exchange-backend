package auth

import coreerrs "github.com/e-scavo/scavo-exchange-backend/internal/core/errs"

type authErrorSpec struct {
	Code    string
	Message string
}

func normalizeAuthError(errCode string) authErrorSpec {
	spec := coreerrs.NormalizeLegacyAuthError(errCode)
	return authErrorSpec{
		Code:    spec.Code,
		Message: spec.Message,
	}
}
