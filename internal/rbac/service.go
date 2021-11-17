package rbac

import (
	"ghoul/internal/model"
	"ghoul/pkg/rbac"
)

// New returns new RBAC service
func New(enableLog bool) *rbac.RBAC {
	r := rbac.NewWithConfig(rbac.Config{EnableLog: enableLog})

	// Add permission for user role
	r.AddPolicy(model.RoleUser, model.ObjectUser, model.ActionViewAll)
	r.AddPolicy(model.RoleUser, model.ObjectCountry, model.ActionViewAll)

	// Add permission for admin role
	r.AddPolicy(model.RoleAdmin, model.ObjectUser, model.ActionAny)
	r.AddPolicy(model.RoleAdmin, model.ObjectCountry, model.ActionAny)

	// Add permission for superadmin role
	r.AddPolicy(model.RoleSuperAdmin, model.ObjectAny, model.ActionAny)

	// Roles inheritance
	r.AddGroupingPolicy(model.RoleAdmin, model.RoleUser)
	r.AddGroupingPolicy(model.RoleSuperAdmin, model.RoleAdmin)

	r.GetModel().PrintPolicy()

	return r
}
