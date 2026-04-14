package authorization

import "testing"

func TestPermissionForResolvesKnownActionResourcePairs(t *testing.T) {
	permission, ok := PermissionFor(Action(" READ "), Resource(" USER "))
	if !ok {
		t.Fatalf("expected permission to resolve")
	}
	if permission != PermissionUserRead {
		t.Fatalf("unexpected permission: %q", permission)
	}
}

func TestHasPermissionRequiresNormalizedSubjectAndPermissionMatch(t *testing.T) {
	subject := AuthorizationSubject{UserID: " user-1 ", Roles: []Role{" USER "}}
	if !HasPermission(subject, PermissionSettingsUpdate) {
		t.Fatalf("expected user role to include settings:update")
	}
	if HasPermission(AuthorizationSubject{}, PermissionSettingsUpdate) {
		t.Fatalf("expected empty subject to be denied")
	}
}

func TestPolicyEvaluatorCanEvaluatesStaticRolePermissions(t *testing.T) {
	evaluator := NewPolicyEvaluator()
	userSubject := AuthorizationSubject{UserID: "user-1", Roles: []Role{RoleUser}}
	adminSubject := AuthorizationSubject{UserID: "admin-1", Roles: []Role{RoleAdmin}}

	if !evaluator.Can(userSubject, ActionRead, ResourceUser) {
		t.Fatalf("expected user to read user resource")
	}
	if evaluator.Can(userSubject, ActionUpdate, ResourceUser) {
		t.Fatalf("expected user to be denied user:update")
	}
	if !evaluator.Can(adminSubject, ActionUpdate, ResourceUser) {
		t.Fatalf("expected admin to update user resource")
	}
}

func TestPolicyEvaluatorEvaluateReturnsDeniedDecisionForUnknownPair(t *testing.T) {
	evaluator := NewPolicyEvaluator()
	decision := evaluator.Evaluate(
		AuthorizationSubject{UserID: "user-1", Roles: []Role{RoleUser}},
		Action("delete"),
		ResourceUser,
	)

	if decision.Allowed {
		t.Fatalf("expected unknown action to be denied")
	}
	if decision.Permission != "" {
		t.Fatalf("expected unknown action to have no permission, got %q", decision.Permission)
	}
}
