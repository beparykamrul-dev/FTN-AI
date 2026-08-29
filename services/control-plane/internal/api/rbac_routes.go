package api

import "net/http"

func ProtectedAdmin(next http.Handler) http.Handler {
	return RequireRole(RoleSuperAdmin, RoleAdmin)(next)
}

func ProtectedNOC(next http.Handler) http.Handler {
	return RequireRole(RoleSuperAdmin, RoleAdmin, RoleEngineer, RoleNOC)(next)
}

func ProtectedBilling(next http.Handler) http.Handler {
	return RequireRole(RoleSuperAdmin, RoleAdmin, RoleBilling)(next)
}
