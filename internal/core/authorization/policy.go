package authorization

import "strings"

type Action string

type Resource string

const (
	ActionRead   Action = "read"
	ActionUpdate Action = "update"
)

const (
	ResourceUser     Resource = "user"
	ResourceSettings Resource = "settings"
)

type Decision struct {
	Action     Action     `json:"action"`
	Resource   Resource   `json:"resource"`
	Permission Permission `json:"permission,omitempty"`
	Allowed    bool       `json:"allowed"`
}

type PolicyEvaluator struct{}

func NewPolicyEvaluator() PolicyEvaluator {
	return PolicyEvaluator{}
}

func (PolicyEvaluator) Evaluate(subject AuthorizationSubject, action Action, resource Resource) Decision {
	normalizedAction := NormalizeAction(action)
	normalizedResource := NormalizeResource(resource)
	decision := Decision{
		Action:   normalizedAction,
		Resource: normalizedResource,
	}

	permission, ok := PermissionFor(normalizedAction, normalizedResource)
	if !ok {
		return decision
	}

	decision.Permission = permission
	decision.Allowed = HasPermission(subject, permission)
	return decision
}

func (e PolicyEvaluator) Can(subject AuthorizationSubject, action Action, resource Resource) bool {
	return e.Evaluate(subject, action, resource).Allowed
}

func Can(subject AuthorizationSubject, action Action, resource Resource) bool {
	return NewPolicyEvaluator().Can(subject, action, resource)
}

func HasPermission(subject AuthorizationSubject, permission Permission) bool {
	normalizedSubject := subject.Normalized()
	if normalizedSubject.UserID == "" || len(normalizedSubject.Roles) == 0 {
		return false
	}

	for _, candidate := range PermissionsForRoles(normalizedSubject.Roles...) {
		if candidate == permission {
			return true
		}
	}

	return false
}

func PermissionFor(action Action, resource Resource) (Permission, bool) {
	switch NormalizeResource(resource) {
	case ResourceUser:
		switch NormalizeAction(action) {
		case ActionRead:
			return PermissionUserRead, true
		case ActionUpdate:
			return PermissionUserUpdate, true
		}
	case ResourceSettings:
		switch NormalizeAction(action) {
		case ActionRead:
			return PermissionSettingsRead, true
		case ActionUpdate:
			return PermissionSettingsUpdate, true
		}
	}

	return "", false
}

func NormalizeAction(action Action) Action {
	return Action(strings.ToLower(strings.TrimSpace(string(action))))
}

func NormalizeResource(resource Resource) Resource {
	return Resource(strings.ToLower(strings.TrimSpace(string(resource))))
}
