package api

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()
	WriteError(w, 400, "bad_request", "invalid request")
	if w.Code != 400 {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), `"code":"bad_request"`) {
		t.Fatalf("unexpected response: %s", w.Body.String())
	}
}
