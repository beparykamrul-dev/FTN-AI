package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProtectedRouteGroups(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
	cases := []struct{name string; h http.Handler; role string; want int}{
		{"admin", ProtectedAdmin(next), RoleAdmin, 204},
		{"admin-deny", ProtectedAdmin(next), RoleNOC, 403},
		{"noc", ProtectedNOC(next), RoleEngineer, 204},
		{"billing", ProtectedBilling(next), RoleBilling, 204},
		{"billing-deny", ProtectedBilling(next), RoleNOC, 403},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("X-FTN-Role", tc.role)
			res := httptest.NewRecorder()
			tc.h.ServeHTTP(res, req)
			if res.Code != tc.want { t.Fatalf("got %d want %d", res.Code, tc.want) }
		})
	}
}
