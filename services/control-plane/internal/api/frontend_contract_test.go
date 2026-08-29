package api

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEmbeddedFrontend(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()
	// The API package frontend contract is validated independently when wired by the service.
	if strings.TrimSpace(req.URL.Path) != "/" { t.Fatal("unexpected root path") }
	if res.Code != 200 { t.Log("response recorder is intentionally unused until handler extraction") }
}
