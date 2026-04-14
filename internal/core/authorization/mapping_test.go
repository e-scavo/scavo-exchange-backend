package authorization

import "testing"

func TestPermissionsForRoleReturnsClone(t *testing.T) {
	permissions := PermissionsForRole(RoleUser)
	if len(permissions) != 3 {
		t.Fatalf("expected 3 permissions, got %d", len(permissions))
	}

	permissions[0] = PermissionUserUpdate
	fresh := PermissionsForRole(RoleUser)
	if fresh[0] != PermissionUserRead {
		t.Fatalf("expected original mapping to stay immutable, got %q", fresh[0])
	}
}

func TestPermissionsForRolesDeduplicatesAndNormalizes(t *testing.T) {
	permissions := PermissionsForRoles(RoleUser, Role(" ADMIN "), RoleAdmin)

	expected := []Permission{
		PermissionUserRead,
		PermissionSettingsRead,
		PermissionSettingsUpdate,
		PermissionUserUpdate,
	}

	if len(permissions) != len(expected) {
		t.Fatalf("expected %d permissions, got %d", len(expected), len(permissions))
	}

	for i := range expected {
		if permissions[i] != expected[i] {
			t.Fatalf("permission %d: expected %q, got %q", i, expected[i], permissions[i])
		}
	}
}
