package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterHealthRequestID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	res := httptest.NewRecorder()
	Router().ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.Code)
	}
	if res.Header().Get("X-Request-ID") == "" {
		t.Fatal("expected request id")
	}
}
