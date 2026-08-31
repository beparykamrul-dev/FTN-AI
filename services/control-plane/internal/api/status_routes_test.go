package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterStatusRoutes(t *testing.T) {
	r := Router()
	paths := []string{"/health", "/ready", "/version", "/api/v1/accounts/status", "/api/v1/billing/status", "/api/v1/noc/status", "/api/v1/service/status"}
	for _, path := range paths {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		if res.Code != http.StatusOK { t.Fatalf("%s: expected 200, got %d", path, res.Code) }
	}
}
