package rbac

// AddRoleForUserID adds a role for a user by ID. Returns false if the user already has the role (aka not affected).
func (s *RBAC) AddRoleForUserID(uid int, role string) bool {
	return s.Enforcer.AddRoleForUser(NormalizeUser(uid), role)
}

// GetRolesForUserID gets the roles that a user has.
func (s *RBAC) GetRolesForUserID(uid int) []string {
	roles, _ := s.Enforcer.GetRolesForUser(NormalizeUser(uid))
	return roles
}

// ReplaceRoleForUserID removes all current roles then adds the new role for a user ID
func (s *RBAC) ReplaceRoleForUserID(uid int, role string) bool {
	s.DeleteRolesForUserID(uid)
	return s.AddRoleForUserID(uid, role)
}

// DeleteRoleForUserID deletes a role for a user ID. Returns false if the user does not have the role (aka not affected).
func (s *RBAC) DeleteRoleForUserID(uid int, role string) bool {
	return s.Enforcer.DeleteRoleForUser(NormalizeUser(uid), role)
}

// DeleteRolesForUserID delete all roles for a user ID. Returns false if the user does not have any roles (aka not affected).
func (s *RBAC) DeleteRolesForUserID(uid int) bool {
	return s.Enforcer.DeleteRolesForUser(NormalizeUser(uid))
}

// DeleteUserID deletes a user ID. Returns false if the user does not exist (aka not affected).
func (s *RBAC) DeleteUserID(uid int) bool {
	return s.Enforcer.DeleteUser(NormalizeUser(uid))
}

// HasRoleForUserID determines whether a user has a role.
func (s *RBAC) HasRoleForUserID(uid int, role string) bool {
	has, _ := s.Enforcer.HasRoleForUser(NormalizeUser(uid), role)
	return has
}

// EnforceUserID determines whether a user ID has permission to do stuff
func (s *RBAC) EnforceUserID(uid int, rvals ...interface{}) bool {
	rvals = append([]interface{}{NormalizeUser(uid)}, rvals...)
	return s.Enforcer.Enforce(rvals...)
}

// AddGroupingPolicy2 adds a role inheritance rule to the current policy.
// If the rule already exists, the function returns false and the rule will not be added.
// Otherwise the function returns true by adding the new rule.
func (s *RBAC) AddGroupingPolicy2(params ...interface{}) bool {
	return s.Enforcer.AddNamedGroupingPolicy("g2", params...)
}

// RemoveGroupingPolicy2 removes a role inheritance rule from the current policy.
func (s *RBAC) RemoveGroupingPolicy2(params ...interface{}) bool {
	return s.Enforcer.RemoveNamedGroupingPolicy("g2", params...)
}
