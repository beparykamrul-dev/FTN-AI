package api

import (
 "net/http"
 "net/http/httptest"
 "testing"
)

func TestRequireRole(t *testing.T) {
 next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
 h := RequireRole("super_admin", "admin")(next)
 for _, tc := range []struct{ role string; want int }{{"admin",204},{"user",403},{"",403}} {
  req := httptest.NewRequest(http.MethodGet, "/", nil); req.Header.Set("X-FTN-Role", tc.role)
  res := httptest.NewRecorder(); h.ServeHTTP(res, req)
  if res.Code != tc.want { t.Fatalf("role %q: got %d want %d", tc.role, res.Code, tc.want) }
 }
}
