package authorization

import "strings"

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

func NormalizeRole(role Role) Role {
	return Role(strings.ToLower(strings.TrimSpace(string(role))))
}

func IsKnownRole(role Role) bool {
	switch NormalizeRole(role) {
	case RoleUser, RoleAdmin:
		return true
	default:
		return false
	}
}
