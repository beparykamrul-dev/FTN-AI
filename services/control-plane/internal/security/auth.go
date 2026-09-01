package security

import (
    "crypto/subtle"
    "net/http"
    "os"
    "strings"
)

// Middleware protects API routes when FTN_API_AUTH_TOKEN is configured.
// Health/readiness endpoints remain public so load balancers can probe safely.
// Authentication is intentionally capability-neutral: authorization belongs in
// the policy layer and must not be inferred from a client supplied header.
func Middleware(next http.Handler) http.Handler {
    token := strings.TrimSpace(os.Getenv("FTN_API_AUTH_TOKEN"))
    if token == "" {
        return next
    }
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/healthz" || r.URL.Path == "/readyz" || r.URL.Path == "/metrics" {
            next.ServeHTTP(w, r)
            return
        }
        auth := strings.TrimSpace(r.Header.Get("Authorization"))
        if !strings.HasPrefix(auth, "Bearer ") {
            http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
            return
        }
        supplied := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
        if subtle.ConstantTimeCompare([]byte(supplied), []byte(token)) != 1 {
            http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}
