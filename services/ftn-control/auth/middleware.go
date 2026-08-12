package auth

import (
	"context"
	"net/http"
)

type contextKey string

const identityContextKey contextKey = "ftn_identity_id"

func WithIdentity(ctx context.Context, identityID string) context.Context {
	return context.WithValue(ctx, identityContextKey, identityID)
}

func IdentityFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(identityContextKey).(string)
	return v, ok && v != ""
}

func RequireService(store *SessionStore, serviceID string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("ftn_session")
		if err != nil || cookie.Value == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		identityID, err := store.Validate(r.Context(), cookie.Value)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		var allowed bool
		allowed, err = store.DB.QueryRowContext(r.Context(), `
			SELECT EXISTS(
				SELECT 1 FROM ftn_service_assignments
				WHERE identity_id = $1 AND service_id = $2 AND status = 'active'
			)`, identityID, serviceID).Scan(&allowed)
		if err != nil || !allowed {
			http.Error(w, "service access denied", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r.WithContext(WithIdentity(r.Context(), identityID)))
	})
}
