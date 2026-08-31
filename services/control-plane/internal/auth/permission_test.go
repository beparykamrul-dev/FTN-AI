package auth

import "testing"

func TestAllowed(t *testing.T) {
	cases := []struct {
		role Role
		permission string
		want bool
	}{
		{RoleSuperAdmin, "anything", true},
		{RoleAdmin, "accounts.read", true},
		{RoleNOC, "billing.write", false},
		{RoleBilling, "billing.write", true},
		{RoleSupport, "network.write", false},
	}

	for _, tc := range cases {
		if got := Allowed(tc.role, tc.permission); got != tc.want {
			t.Fatalf("Allowed(%q, %q) = %v, want %v", tc.role, tc.permission, got, tc.want)
		}
	}
}
