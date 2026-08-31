package auth

type Role string

const (
	RoleSuperAdmin Role = "super_admin"
	RoleAdmin      Role = "admin"
	RoleNOC        Role = "noc"
	RoleBilling    Role = "billing"
	RoleSupport    Role = "support"
)

var permissions = map[Role]map[string]bool{
	RoleSuperAdmin: {"*": true},
	RoleAdmin:      {"accounts.read": true, "accounts.write": true, "billing.read": true, "billing.write": true, "network.read": true, "noc.read": true},
	RoleNOC:        {"network.read": true, "noc.read": true, "noc.write": true},
	RoleBilling:   {"accounts.read": true, "billing.read": true, "billing.write": true},
	RoleSupport:   {"accounts.read": true, "billing.read": true, "noc.read": true},
}

func Allowed(role Role, permission string) bool {
	p := permissions[role]
	return p != nil && (p["*"] || p[permission])
}
