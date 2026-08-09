package httpx

import (
	"net/http"
	"strings"
)

type CORS struct {
	AllowedOrigins []string
}

func (c CORS) Middleware(next http.Handler) http.Handler {
	allowed := map[string]bool{}
	for _, origin := range c.AllowedOrigins { allowed[origin] = true }
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && allowed[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Request-ID")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		}
		if strings.EqualFold(r.Method, http.MethodOptions) { w.WriteHeader(http.StatusNoContent); return }
		next.ServeHTTP(w,r)
	})
}
