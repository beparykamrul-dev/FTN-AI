package httpx

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSMiddlewareRejectsWildcardOriginConfiguration(t *testing.T) {
	h := CORS{AllowedOrigins: []string{"*", "https://ftn.example"}}.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }))
	r := httptest.NewRequest(http.MethodOptions, "/", nil)
	r.Header.Set("Origin", "https://evil.example")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusForbidden { t.Fatalf("expected forbidden, got %d", w.Code) }
}

func TestCORSMiddlewareAllowsExplicitOrigin(t *testing.T) {
	h := CORS{AllowedOrigins: []string{"https://ftn.example"}}.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }))
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Origin", "https://ftn.example")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusOK { t.Fatalf("expected OK, got %d", w.Code) }
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "https://ftn.example" { t.Fatalf("unexpected origin header: %q", got) }
}
