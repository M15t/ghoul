package rbac

import (
	"strconv"
)

// NormalizeRole corrects role ID for RBAC service
func NormalizeRole(r int) string {
	return "r" + strconv.Itoa(r)
}

// UnnormalizeRole converts RBAC role back to normal
func UnnormalizeRole(r string) int {
	iRole, _ := strconv.Atoi(r[1:len(r)])
	return iRole
}

// NormalizeUser corrects user ID for RBAC service
func NormalizeUser(uid int) string {
	return strconv.Itoa(uid)
}
