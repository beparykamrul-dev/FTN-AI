package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIRoot(t *testing.T) {
	a := &App{catalog: catalog}
	r := httptest.NewRequest(http.MethodGet, "/api/v1", nil)
	w := httptest.NewRecorder()
	a.apiRoot(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestAPIMeUsesRequestContext(t *testing.T) {
	a := &App{catalog: catalog}
	r := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	r = withRequestContext(r, requestContext{RequestID: "req-test", CorrelationID: "corr-test", PrincipalID: "principal-test", TenantID: "tenant-test"})
	w := httptest.NewRecorder()
	a.apiMe(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	if got := w.Body.String(); got == "" || !containsAll(got, "principal-test", "tenant-test", "req-test") {
		t.Fatalf("unexpected response: %s", got)
	}
}

func containsAll(s string, values ...string) bool {
	for _, value := range values {
		if !contains(s, value) {
			return false
		}
	}
	return true
}

func contains(s, value string) bool {
	return len(value) == 0 || (len(s) >= len(value) && stringIndex(s, value) >= 0)
}

func stringIndex(s, value string) int {
	for i := 0; i+len(value) <= len(s); i++ {
		if s[i:i+len(value)] == value {
			return i
		}
	}
	return -1
}
