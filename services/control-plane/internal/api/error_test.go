package api

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteError(t *testing.T) {
	res := httptest.NewRecorder()
	WriteError(res, 400, "bad_request", "invalid request")
	if res.Code != 400 { t.Fatalf("expected 400, got %d", res.Code) }
	if !strings.Contains(res.Body.String(), `"code":"bad_request"`) { t.Fatalf("unexpected response: %s", res.Body.String()) }
}
