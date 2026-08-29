package rbac

import "strings"

type Role string

const (
	RoleCustomer    Role = "customer"
	RoleEmployee    Role = "employee"
	RoleEngineer    Role = "engineer"
	RoleReseller    Role = "reseller"
	RolePartner     Role = "partner"
	RoleSupport     Role = "support-agent"
	RoleBilling     Role = "billing-agent"
	RoleNOC         Role = "noc-operator"
	RoleAdmin       Role = "admin"
	RoleSuperAdmin  Role = "super-admin"
)

type Permission string

const (
	ReadOwnAccount   Permission = "account:read:self"
	UpdateOwnAccount Permission = "account:update:self"
	ReadBilling      Permission = "billing:read"
	CreatePayment    Permission = "payment:create"
	ReadNetwork      Permission = "network:read"
	CreateTicket     Permission = "support:create"
	ProposeRecovery  Permission = "recovery:propose"
	ApproveRecovery  Permission = "recovery:approve"
	ExecuteRecovery  Permission = "recovery:execute"
)

var rolePermissions = map[Role]map[Permission]struct{}{
	RoleCustomer: {ReadOwnAccount: {}, UpdateOwnAccount: {}, ReadBilling: {}, CreatePayment: {}, ReadNetwork: {}, CreateTicket: {}},
	RoleEngineer: {ReadNetwork: {}, CreateTicket: {}, ProposeRecovery: {}},
	RoleNOC:      {ReadNetwork: {}, ProposeRecovery: {}, ApproveRecovery: {}},
	RoleAdmin:    {ReadBilling: {}, ReadNetwork: {}, ProposeRecovery: {}, ApproveRecovery: {}},
	RoleSuperAdmin: {ReadBilling: {}, ReadNetwork: {}, ProposeRecovery: {}, ApproveRecovery: {}, ExecuteRecovery: {}},
}

func Allows(role Role, permission Permission) bool {
	return strings.TrimSpace(string(role)) != "" && func() bool {
		_, ok := rolePermissions[role][permission]
		return ok
	}()
}
