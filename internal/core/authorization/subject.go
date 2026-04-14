package authorization

import "strings"

type AuthorizationSubject struct {
	UserID string `json:"user_id"`
	Roles  []Role `json:"roles"`
}

func (s AuthorizationSubject) Normalized() AuthorizationSubject {
	normalized := AuthorizationSubject{
		UserID: strings.TrimSpace(s.UserID),
		Roles:  make([]Role, 0, len(s.Roles)),
	}

	seen := make(map[Role]struct{})
	for _, role := range s.Roles {
		nr := NormalizeRole(role)
		if nr == "" {
			continue
		}
		if _, exists := seen[nr]; exists {
			continue
		}
		seen[nr] = struct{}{}
		normalized.Roles = append(normalized.Roles, nr)
	}

	if len(normalized.Roles) == 0 {
		normalized.Roles = nil
	}

	return normalized
}
