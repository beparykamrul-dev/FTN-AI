package api

import "net/http"

func RequireRole(roles ...string) func(http.Handler) http.Handler {
	allowed := map[string]bool{}
	for _, role := range roles { allowed[role] = true }
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := r.Header.Get("X-FTN-Role")
			if !allowed[role] {
				WriteError(w, http.StatusForbidden, "forbidden", "role is not authorized")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
