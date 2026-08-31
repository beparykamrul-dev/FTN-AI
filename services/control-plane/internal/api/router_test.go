package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterRoutes(t *testing.T) {
	for _, path := range []string{"/health", "/ready", "/version"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		res := httptest.NewRecorder()
		Router().ServeHTTP(res, req)
		if res.Code != http.StatusOK {
			t.Fatalf("%s: expected 200, got %d", path, res.Code)
		}
	}
}
