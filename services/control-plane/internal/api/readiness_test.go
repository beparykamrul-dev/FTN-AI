package api

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReadiness(t *testing.T) {
	req := httptest.NewRequest("GET", "/ready", nil)
	res := httptest.NewRecorder()
	Readiness(res, req)
	if res.Code != 200 { t.Fatalf("expected 200, got %d", res.Code) }
	if !strings.Contains(res.Body.String(), `"ready":true`) { t.Fatalf("unexpected response: %s", res.Body.String()) }
}
