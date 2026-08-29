package auth

type Role string

const (
	RoleSuperAdmin Role = "super-admin"
	RoleAdmin      Role = "admin"
	RoleBilling    Role = "billing"
	RoleNOC        Role = "noc"
	RoleEngineer   Role = "engineer"
	RoleSupport    Role = "support"
)

var permissions = map[Role]map[string]struct{}{
	RoleSuperAdmin: {"accounts:read": {}, "accounts:write": {}, "billing:read": {}, "billing:write": {}, "payments:read": {}, "payments:write": {}, "network:read": {}, "network:write": {}, "incidents:read": {}, "incidents:write": {}, "ai:read": {}, "ai:approve": {}, "audit:read": {}},
	RoleAdmin:      {"accounts:read": {}, "accounts:write": {}, "billing:read": {}, "billing:write": {}, "payments:read": {}, "payments:write": {}, "network:read": {}, "incidents:read": {}, "incidents:write": {}, "audit:read": {}},
	RoleBilling:    {"accounts:read": {}, "billing:read": {}, "billing:write": {}, "payments:read": {}, "payments:write": {}},
	RoleNOC:        {"network:read": {}, "network:write": {}, "incidents:read": {}, "incidents:write": {}, "ai:read": {}, "ai:approve": {}, "audit:read": {}},
	RoleEngineer:   {"network:read": {}, "network:write": {}, "incidents:read": {}, "incidents:write": {}},
	RoleSupport:    {"accounts:read": {}, "incidents:read": {}},
}

func Allowed(role Role, permission string) bool {
	_, ok := permissions[role][permission]
	return ok
}
