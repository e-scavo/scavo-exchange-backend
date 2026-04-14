package authorization

type Permission string

const (
	PermissionUserRead       Permission = "user:read"
	PermissionUserUpdate     Permission = "user:update"
	PermissionSettingsRead   Permission = "settings:read"
	PermissionSettingsUpdate Permission = "settings:update"
)
