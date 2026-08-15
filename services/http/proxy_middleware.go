package http

import "net/http"

// ProxyPolicy defines the safe boundary for FTN HTTP proxy middleware.
type ProxyPolicy struct {
	AllowedHosts []string
	RequireTLS   bool
}

func HostAllowed(host string, p ProxyPolicy) bool {
	for _, allowed := range p.AllowedHosts {
		if host == allowed {
			return true
		}
	}
	return false
}

// Middleware exposes policy validation as standard net/http middleware.
func Middleware(p ProxyPolicy, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if p.RequireTLS && r.TLS == nil {
			http.Error(w, "TLS required", http.StatusUpgradeRequired)
			return
		}
		if !HostAllowed(r.Host, p) {
			http.Error(w, "host not allowed", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
