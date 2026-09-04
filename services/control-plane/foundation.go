package main

import (
    "context"
    "crypto/rand"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "strings"
    "time"

    "github.com/jackc/pgx/v5"
)

type requestContext struct { RequestID string; CorrelationID string; PrincipalID string; TenantID string }
type ctxKey struct{}

func withRequestContext(r *http.Request, rc requestContext) *http.Request { return r.WithContext(context.WithValue(r.Context(), ctxKey{}, rc)) }
func requestInfo(r *http.Request) requestContext { if v,ok:=r.Context().Value(ctxKey{}).(requestContext); ok{return v}; return requestContext{RequestID:r.Header.Get("X-Request-ID"),CorrelationID:r.Header.Get("X-Correlation-ID")} }
func randomID(prefix string) string { b:=make([]byte,16); if _,err:=rand.Read(b);err!=nil{return fmt.Sprintf("%s-%d",prefix,time.Now().UnixNano())}; return prefix+"-"+hex.EncodeToString(b) }
func bearerToken(r *http.Request) string { a:=strings.TrimSpace(r.Header.Get("Authorization")); if !strings.HasPrefix(a,"Bearer "){return ""}; return strings.TrimSpace(strings.TrimPrefix(a,"Bearer ")) }
func tokenHash(token string) string { s:=sha256.Sum256([]byte(token)); return hex.EncodeToString(s[:]) }

func (a *App) authenticateDB(r *http.Request) (requestContext,error) {
    if a.db==nil{return requestContext{},errors.New("database authentication unavailable")}
    token:=bearerToken(r)
    if token==""{return requestContext{},errors.New("missing bearer token")}
    var rc requestContext
    err:=a.db.QueryRow(r.Context(),`select p.id::text,p.tenant_id::text from api_credentials c join principals p on p.id=c.principal_id join tenants t on t.id=p.tenant_id where c.token_hash=$1 and c.revoked_at is null and (c.expires_at is null or c.expires_at>now()) and p.status='active' and t.status='active'`,tokenHash(token)).Scan(&rc.PrincipalID,&rc.TenantID)
    if err!=nil{return requestContext{},errors.New("invalid or inactive credential")}; return rc,nil
}
func (a *App) authorize(r *http.Request, permission string) error {
    rc:=requestInfo(r); if a.db==nil||rc.PrincipalID==""{return errors.New("authorization context unavailable")}; var allowed bool
    err:=a.db.QueryRow(r.Context(),`select exists(select 1 from principal_roles pr join roles ro on ro.id=pr.role_id join role_permissions rp on rp.role_id=ro.id join permissions p on p.id=rp.permission_id where pr.principal_id=$1 and p.key=$2)`,rc.PrincipalID,permission).Scan(&allowed)
    if err!=nil||!allowed{return errors.New("permission denied")}; return nil
}
func (a *App) audit(r *http.Request,action,resource,outcome string,metadata any){if a.db==nil{return}; rc:=requestInfo(r); b,_:=json.Marshal(metadata); _,_=a.db.Exec(r.Context(),`insert into audit_events(tenant_id,principal_id,action,resource,outcome,correlation_id,request_id,metadata) values(nullif($1,'')::uuid,nullif($2,'')::uuid,$3,$4,$5,$6,$7,$8::jsonb)`,rc.TenantID,rc.PrincipalID,action,resource,outcome,rc.CorrelationID,rc.RequestID,string(b))}

func (a *App) ensureSystemIdentity(ctx context.Context,token string) error {
    if a.db==nil||strings.TrimSpace(token)==""{return nil}; tx,err:=a.db.Begin(ctx); if err!=nil{return err}; defer tx.Rollback(ctx)
    var tenantID,principalID,roleID string
    if err=tx.QueryRow(ctx,`insert into tenants(slug,name) values('ftn-system','FTN System') on conflict(slug) do update set updated_at=now() returning id::text`).Scan(&tenantID);err!=nil{return err}
    if err=tx.QueryRow(ctx,`insert into principals(tenant_id,subject,kind,issuer) values($1,'control-plane','service','ftn') on conflict(tenant_id,issuer,subject) do update set updated_at=now() returning id::text`,tenantID).Scan(&principalID);err!=nil{return err}
    if err=tx.QueryRow(ctx,`insert into roles(tenant_id,name,description) values($1,'system-admin','FTN system control-plane role') on conflict(tenant_id,name) do update set description=excluded.description returning id::text`,tenantID).Scan(&roleID);err!=nil{return err}
    if _,err=tx.Exec(ctx,`insert into principal_roles(principal_id,role_id) values($1,$2) on conflict do nothing`,principalID,roleID);err!=nil{return err}
    if _,err=tx.Exec(ctx,`insert into role_permissions(role_id,permission_id) select $1,id from permissions on conflict do nothing`,roleID);err!=nil{return err}
    if _,err=tx.Exec(ctx,`insert into api_credentials(principal_id,name,token_hash) values($1,'bootstrap',$2) on conflict(token_hash) do nothing`,principalID,tokenHash(token));err!=nil{return err}
    for _,s:=range catalog{if _,err=tx.Exec(ctx,`insert into service_entitlements(tenant_id,principal_id,service_id,active) values($1,$2,$3,true) on conflict(principal_id,service_id) do update set active=true`,tenantID,principalID,s.ID);err!=nil{return err}}
    return tx.Commit(ctx)
}

func (a *App) idempotentResponse(r *http.Request,key,bodyHash string)([]byte,int,bool,error){if a.db==nil||key==""{return nil,0,false,nil};var status int;var body []byte;err:=a.db.QueryRow(r.Context(),`select response_status,response_body::text from idempotency_keys where key=$1 and expires_at>now()`,key).Scan(&status,&body);if errors.Is(err,pgx.ErrNoRows){return nil,0,false,nil};if err!=nil{return nil,0,false,err};var stored map[string]any;if err:=json.Unmarshal(body,&stored);err!=nil{return nil,0,false,err};if h,ok:=stored["_request_hash"].(string);ok&&h!=bodyHash{return nil,0,false,errors.New("idempotency key reused with different request")};return body,status,true,nil}
func (a *App) saveIdempotency(r *http.Request,key,bodyHash string,status int,body []byte)error{if a.db==nil||key==""{return nil};var payload map[string]any;if err:=json.Unmarshal(body,&payload);err!=nil{payload=map[string]any{"response":string(body)}};payload["_request_hash"]=bodyHash;encoded,_:=json.Marshal(payload);_,err:=a.db.Exec(r.Context(),`insert into idempotency_keys(key,principal_id,request_hash,response_status,response_body,expires_at) values($1,nullif($2,'')::uuid,$3,$4,$5::jsonb,now()+interval '24 hours') on conflict(key) do update set response_status=excluded.response_status,response_body=excluded.response_body,request_hash=excluded.request_hash,expires_at=excluded.expires_at`,key,requestInfo(r).PrincipalID,bodyHash,status,string(encoded));return err}
func writeAuthorizedJSON(w http.ResponseWriter,status int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(status);_=json.NewEncoder(w).Encode(v)}
func requirePermission(a *App,permission string,w http.ResponseWriter,r *http.Request)bool{if err:=a.authorize(r,permission);err!=nil{writeAuthorizedJSON(w,http.StatusForbidden,map[string]string{"error":"forbidden"});a.audit(r,"authorization.denied",r.URL.Path,"denied",map[string]string{"permission":permission});return false};return true}
func requestBodyHash(v any)string{b,_:=json.Marshal(v);s:=sha256.Sum256(b);return hex.EncodeToString(s[:])}
