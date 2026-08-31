package api

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	t.Setenv("FTN_VERSION", "test")
	req := httptest.NewRequest("GET", "/version", nil)
	res := httptest.NewRecorder()
	Version(res, req)
	if res.Code != 200 { t.Fatalf("expected 200, got %d", res.Code) }
	if !strings.Contains(res.Body.String(), `"version":"test"`) { t.Fatalf("unexpected response: %s", res.Body.String()) }
}
