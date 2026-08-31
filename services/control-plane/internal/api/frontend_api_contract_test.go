package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFrontendAPIContract(t *testing.T) {
	r := Router()
	for _, path := range []string{"/health", "/ready", "/version"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		if res.Code != http.StatusOK { t.Fatalf("%s: got %d", path, res.Code) }
	}
}
