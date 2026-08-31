package api

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAggregateServiceStatus(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/status", nil)
	res := httptest.NewRecorder()
	AggregateServiceStatus(res, req)
	if res.Code != 200 { t.Fatalf("expected 200, got %d", res.Code) }
	body := res.Body.String()
	if !strings.Contains(body, `"status":"ok"`) || !strings.Contains(body, `"accounts"`) { t.Fatalf("unexpected response: %s", body) }
}
