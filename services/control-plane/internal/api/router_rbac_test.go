package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterRBAC(t *testing.T) {
	r := Router()
	cases := []struct{path, role string; want int}{
		{"/api/v1/billing/status", "billing", http.StatusOK},
		{"/api/v1/billing/status", "noc", http.StatusForbidden},
		{"/api/v1/noc/status", "noc", http.StatusOK},
		{"/api/v1/noc/status", "billing", http.StatusForbidden},
		{"/api/v1/accounts/status", "admin", http.StatusOK},
	}
	for _, tc := range cases {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)
		req.Header.Set("X-FTN-Role", tc.role)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		if res.Code != tc.want { t.Fatalf("%s role=%s: got %d want %d", tc.path, tc.role, res.Code, tc.want) }
	}
}
