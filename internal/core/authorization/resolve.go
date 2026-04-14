package authorization

import (
	"strings"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
)

func SubjectFromClaims(claims *coreauth.Claims) (AuthorizationSubject, bool) {
	if claims == nil {
		return AuthorizationSubject{}, false
	}
	if strings.TrimSpace(claims.UserID) == "" {
		return AuthorizationSubject{}, false
	}

	subject := AuthorizationSubject{
		UserID: claims.UserID,
		Roles:  []Role{RoleUser},
	}

	return subject.Normalized(), true
}
