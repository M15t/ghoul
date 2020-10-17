package model

// RBAC roles
const (
	RoleSuperAdmin = "superadmin"
	RoleAdmin      = "admin"
	RoleUser       = "user"
)

// AvailableRoles for validation
var AvailableRoles = []string{RoleAdmin, RoleUser}

// RBAC objects
const (
	ObjectAny     = "*"
	ObjectUser    = "user"
	ObjectCountry = "country"
)

// RBAC actions
const (
	ActionAny       = "*"
	ActionViewAll   = "view_all"
	ActionView      = "view"
	ActionCreateAll = "create_all"
	ActionCreate    = "create"
	ActionUpdateAll = "update_all"
	ActionUpdate    = "update"
	ActionDeleteAll = "delete_all"
	ActionDelete    = "delete"
)
