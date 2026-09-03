package main

// Route wiring is maintained in the existing control-plane implementation.
// This file intentionally delegates the interactive UI to control_panel.go.

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "strings"
    "sync/atomic"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

var requestCount uint64

type App struct { db *pgxpool.Pool; catalog []Service }

func jsonResponse(w http.ResponseWriter,status int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(status);_=json.NewEncoder(w).Encode(v)}
func method(w http.ResponseWriter,r *http.Request,want string)bool{if r.Method!=want{jsonResponse(w,http.StatusMethodNotAllowed,map[string]string{"error":"method_not_allowed"});return false};return true}

func(a *App)routes()http.Handler{m:=http.NewServeMux();m.HandleFunc("/healthz",a.health);m.HandleFunc("/readyz",a.ready);m.HandleFunc("/metrics",metricsHandler);m.HandleFunc("/api/v1",a.apiRoot);m.HandleFunc("/api/v1/",a.apiRoot);m.HandleFunc("/api/v1/me",a.apiMe);m.HandleFunc("/api/v1/capabilities",a.apiCapabilities);m.HandleFunc("/api/v1/dashboard",a.apiDashboard);m.HandleFunc("/api/v1/ai/status",a.aiStatus);m.HandleFunc("/api/v1/ai/analyze",a.aiAnalyze);m.HandleFunc("/api/v1/services",a.serviceCatalog);m.HandleFunc("/api/v1/entitlements",a.entitlements);m.HandleFunc("/api/v1/service-requests",a.requests);m.HandleFunc("/api/v1/events",a.readEvents);m.HandleFunc("/api/v1/events/append",a.appendEvent);m.HandleFunc("/api/v1/events/offset",a.eventOffset);m.HandleFunc("/api/v1/events/offset/commit",a.commitEventOffset);m.HandleFunc("/api/v1/approvals",a.approvalCollection);m.HandleFunc("/api/v1/approvals/",a.approvalRouter);m.HandleFunc("/api/v1/jobs",a.jobCollection);m.HandleFunc("/api/v1/jobs/claim",a.claimJob);m.HandleFunc("/api/v1/jobs/",a.jobRouter);m.HandleFunc("/api/v1/nodes",a.nodeCatalog);m.HandleFunc("/api/v1/nodes/register",a.registerNode);m.HandleFunc("/api/v1/nodes/heartbeat",a.nodeHeartbeat);m.HandleFunc("/api/v1/placement/preview",a.placement);m.HandleFunc("/api/v1/network/health",a.networkHealth);m.HandleFunc("/api/v1/network/integration",a.networkIntegration);m.HandleFunc("/api/v1/exrouter/route",a.exRouter);m.HandleFunc("/api/v1/qkd/nodes",a.qkdNodes);m.HandleFunc("/api/v1/qkd/links",a.qkdLinks);m.HandleFunc("/api/v1/qkd/status",a.qkdStatus);m.HandleFunc("/api/v1/qkd/kms",a.qkdKMS);m.HandleFunc("/api/v1/qkd/intents",a.qkdIntent);m.HandleFunc("/api/v1/data-governor/assets",a.dataGovernor);m.HandleFunc("/api/v1/data-governor/requests",a.dataRequest);m.HandleFunc("/api/v1/dns-guard/profiles",a.dnsGuardProfiles);m.HandleFunc("/api/v1/dns-guard/bindings",a.dnsGuardBindings);m.HandleFunc("/api/v1/dns-guard/events",a.dnsGuardEvents);m.HandleFunc("/api/v1/dns-guard/summary",a.dnsGuardSummary);m.HandleFunc("/api/v1/dns-guard/requests",a.dnsGuardRequest);m.HandleFunc("/control-panel",a.controlPanel);m.HandleFunc("/",a.frontend);return security.Middleware(requestContextMiddleware(a,securityHeaders(counting(m))))}

func requestContextMiddleware(a *App,next http.Handler)http.Handler{return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){rid:=strings.TrimSpace(r.Header.Get("X-Request-ID"));if rid==""{rid=randomID("req")};cid:=strings.TrimSpace(r.Header.Get("X-Correlation-ID"));if cid==""{cid=rid};w.Header().Set("X-Request-ID",rid);w.Header().Set("X-Correlation-ID",cid);rc:=requestContext{RequestID:rid,CorrelationID:cid};if r.URL.Path!="/healthz"&&r.URL.Path!="/readyz"&&r.URL.Path!="/metrics"{var err error;rc,err=a.authenticateDB(r);if err!=nil{jsonResponse(w,http.StatusUnauthorized,map[string]string{"error":"unauthorized"});return};rc.RequestID=rid;rc.CorrelationID=cid};next.ServeHTTP(w,withRequestContext(r,rc))})}
func(a *App)health(w http.ResponseWriter,r *http.Request){if !method(w,r,http.MethodGet){return};jsonResponse(w,200,map[string]any{"status":"ok","time":time.Now().UTC()})}
func(a *App)ready(w http.ResponseWriter,r *http.Request){if !method(w,r,http.MethodGet){return};if a.db==nil{jsonResponse(w,200,map[string]string{"status":"ready","database":"disabled"});return};ctx,c:=context.WithTimeout(r.Context(),2*time.Second);defer c();if err:=a.db.Ping(ctx);err!=nil{jsonResponse(w,503,map[string]string{"status":"not_ready","database":"unavailable"});return};jsonResponse(w,200,map[string]string{"status":"ready","database":"ok"})}
func(a *App)frontend(w http.ResponseWriter,r *http.Request){if r.URL.Path!="/"{http.NotFound(w,r);return};w.Header().Set("Content-Type","text/html; charset=utf-8");_,_=w.Write([]byte(indexHTML))}
func counting(next http.Handler)http.Handler{return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){atomic.AddUint64(&requestCount,1);next.ServeHTTP(w,r)})}

func main(){if len(os.Args)>1&&os.Args[1]=="--healthcheck"{healthcheck()};addr:=os.Getenv("FTN_CONTROL_ADDR");if addr==""{addr=":8080"};var db *pgxpool.Pool;if dsn:=os.Getenv("DATABASE_URL");dsn!=""{ctx,c:=context.WithTimeout(context.Background(),10*time.Second);defer c();var err error;db,err=pgxpool.New(ctx,dsn);if err!=nil{log.Fatal(err)};if err=db.Ping(ctx);err!=nil{log.Fatal(err)};if err=(&App{db:db,catalog:catalog}).ensureSystemIdentity(ctx,os.Getenv("FTN_API_AUTH_TOKEN"));err!=nil{log.Fatal(err)};defer db.Close()};app:=&App{db:db,catalog:catalog};srv:=&http.Server{Addr:addr,Handler:app.routes(),ReadHeaderTimeout:5*time.Second,ReadTimeout:15*time.Second,WriteTimeout:15*time.Second,IdleTimeout:60*time.Second};log.Printf("FTN control plane listening on %s",addr);log.Fatal(srv.ListenAndServe())}
