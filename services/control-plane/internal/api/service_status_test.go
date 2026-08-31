package api

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServiceStatus(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/service/status", nil)
	res := httptest.NewRecorder()
	ServiceStatus(res, req)
	if res.Code != 200 { t.Fatalf("expected 200, got %d", res.Code) }
	if !strings.Contains(res.Body.String(), `"status":"ok"`) { t.Fatalf("unexpected response: %s", res.Body.String()) }
}
