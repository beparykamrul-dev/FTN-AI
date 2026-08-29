package auth

import "net/http"

type contextKey struct{}

const roleHeader = "X-FTN-Role"

func RoleFromRequest(r *http.Request) Role {
	return Role(r.Header.Get(roleHeader))
}

func Require(permission string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := RoleFromRequest(r)
		if !Allowed(role, permission) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"error":"forbidden","permission":"` + permission + `"}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}
