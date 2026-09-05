package auth

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string
const identityContextKey contextKey = "ftn_identity_id"
func WithIdentity(ctx context.Context, identityID string) context.Context { if ctx==nil{return context.Background()};return context.WithValue(ctx,identityContextKey,strings.TrimSpace(identityID)) }
func IdentityFromContext(ctx context.Context)(string,bool){if ctx==nil{return "",false};v,ok:=ctx.Value(identityContextKey).(string);v=strings.TrimSpace(v);return v,ok&&v!=""}
func RequireService(store *SessionStore,serviceID string,next http.Handler)http.Handler{return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){if store==nil||store.DB==nil||next==nil||strings.TrimSpace(serviceID)==""{http.Error(w,"service unavailable",http.StatusServiceUnavailable);return};if r==nil||r.Context()==nil{http.Error(w,"bad request",http.StatusBadRequest);return};cookie,err:=r.Cookie("ftn_session");if err!=nil||strings.TrimSpace(cookie.Value)==""{http.Error(w,"unauthorized",http.StatusUnauthorized);return};identityID,err:=store.Validate(r.Context(),strings.TrimSpace(cookie.Value));if err!=nil||strings.TrimSpace(identityID)==""{http.Error(w,"unauthorized",http.StatusUnauthorized);return};var allowed bool;allowed,err=store.DB.QueryRowContext(r.Context(),`SELECT EXISTS(SELECT 1 FROM ftn_service_assignments WHERE identity_id = $1 AND service_id = $2 AND status = 'active')`,identityID,strings.TrimSpace(serviceID)).Scan(&allowed);if err!=nil||!allowed{http.Error(w,"service access denied",http.StatusForbidden);return};next.ServeHTTP(w,r.WithContext(WithIdentity(r.Context(),identityID)))})}
