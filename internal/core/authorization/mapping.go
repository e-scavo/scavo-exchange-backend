package authorization

var rolePermissions = map[Role][]Permission{
	RoleUser: {
		PermissionUserRead,
		PermissionSettingsRead,
		PermissionSettingsUpdate,
	},
	RoleAdmin: {
		PermissionUserRead,
		PermissionUserUpdate,
		PermissionSettingsRead,
		PermissionSettingsUpdate,
	},
}

func PermissionsForRole(role Role) []Permission {
	permissions, ok := rolePermissions[NormalizeRole(role)]
	if !ok {
		return nil
	}

	cloned := make([]Permission, len(permissions))
	copy(cloned, permissions)
	return cloned
}

func PermissionsForRoles(roles ...Role) []Permission {
	if len(roles) == 0 {
		return nil
	}

	seen := make(map[Permission]struct{})
	permissions := make([]Permission, 0)
	for _, role := range roles {
		for _, permission := range PermissionsForRole(role) {
			if _, exists := seen[permission]; exists {
				continue
			}
			seen[permission] = struct{}{}
			permissions = append(permissions, permission)
		}
	}

	if len(permissions) == 0 {
		return nil
	}

	return permissions
}
