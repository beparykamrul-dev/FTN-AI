package main

import (
    "context"
    "encoding/json"
    "errors"
    "log"
    "net"
    "net/http"
    "os"
    "sort"
    "strings"
    "sync"
    "time"
)

type ManagedEndpoint struct { ServiceID string `json:"service_id"`; CIDR string `json:"cidr"`; Region string `json:"region,omitempty"`; Provider string `json:"provider,omitempty"`; ExpiresAt time.Time `json:"expires_at,omitempty"`; network *net.IPNet }
type TrafficFlowObservation struct { FlowRecord; ServiceID string `json:"service_id"`; Class TrafficClass `json:"class"`; PathID string `json:"path_id,omitempty"`; Region string `json:"region,omitempty"`; Provider string `json:"provider,omitempty"`; ObservedAt time.Time `json:"observed_at"` }
type TrafficRuntime struct { mu sync.RWMutex; endpoints []ManagedEndpoint; flows []TrafficFlowObservation; decisions map[string]TrafficDecision; controllers map[string]*TrafficPathController; quality *TrafficQualityStore; listener *FlowListener; activeProbe *TrafficActiveProbe }

func NewTrafficRuntime() *TrafficRuntime {
    runtime := &TrafficRuntime{decisions: make(map[string]TrafficDecision), controllers: make(map[string]*TrafficPathController), quality: NewTrafficQualityStore()}
    if raw := strings.TrimSpace(os.Getenv("FTN_TRAFFIC_PROBE_TARGETS")); raw != "" {
        interval := durationFromEnv("FTN_TRAFFIC_PROBE_INTERVAL", trafficProbeDefaultInterval)
        timeout := durationFromEnv("FTN_TRAFFIC_PROBE_TIMEOUT", trafficProbeDefaultTimeout)
        probe := NewTrafficActiveProbe(runtime, interval, timeout)
        if err := probe.Start(context.Background(), trafficProbeTargetsFromEnv); err != nil { log.Printf("FTN active traffic probe disabled: %v", err) } else { runtime.activeProbe = probe; log.Printf("FTN active traffic probe enabled") }
    }
    if cfg, enabled := flowListenerConfigFromEnv(); enabled {
        listener, err := NewFlowListener(cfg, NewFlowTelemetryCollector(), runtime)
        if err != nil { log.Printf("FTN flow telemetry disabled: %v", err) } else if err := listener.Start(context.Background()); err != nil { log.Printf("FTN flow telemetry listener failed to start: %v", err) } else { runtime.listener = listener; log.Printf("FTN flow telemetry listener active on %s", cfg.Address) }
    }
    return runtime
}

func (t *TrafficRuntime) StartActiveProbes(ctx context.Context, targets func(context.Context) []TrafficProbeTarget) error {
    interval := durationFromEnv("FTN_TRAFFIC_PROBE_INTERVAL", trafficProbeDefaultInterval); timeout := durationFromEnv("FTN_TRAFFIC_PROBE_TIMEOUT", trafficProbeDefaultTimeout); probe := NewTrafficActiveProbe(t, interval, timeout)
    if err := probe.Start(ctx, func() []TrafficProbeTarget { return targets(ctx) }); err != nil { return err }
    t.mu.Lock(); t.activeProbe = probe; t.mu.Unlock(); return nil
}
func (t *TrafficRuntime) Close() error { t.mu.RLock(); listener, probe := t.listener, t.activeProbe; t.mu.RUnlock(); var first error; if probe != nil { if err:=probe.Close(); err!=nil { first=err } }; if listener!=nil { if err:=listener.Close(); first==nil { first=err } }; return first }
func durationFromEnv(name string, fallback time.Duration) time.Duration { raw:=strings.TrimSpace(os.Getenv(name)); if raw=="" { return fallback }; d,err:=time.ParseDuration(raw); if err!=nil||d<=0{return fallback}; return d }
func (t *TrafficRuntime) FlowCount() int { t.mu.RLock(); defer t.mu.RUnlock(); return len(t.flows) }
func (t *TrafficRuntime) UpsertQuality(o TrafficQualityObservation, now time.Time) error { if t.quality==nil{return errors.New("traffic_quality_store_required")}; return t.quality.Upsert(o,now) }
func (t *TrafficRuntime) QualitySnapshot(serviceID string, now time.Time) []TrafficQualityObservation { if t.quality==nil{return nil}; return t.quality.Snapshot(serviceID,now) }
func (t *TrafficRuntime) UpsertEndpoint(e ManagedEndpoint) error { e.ServiceID=strings.TrimSpace(e.ServiceID); e.CIDR=strings.TrimSpace(e.CIDR); if e.ServiceID==""||e.CIDR==""{return errors.New("service_id_and_cidr_required")}; if _,ok:=trafficPolicyByID(e.ServiceID);!ok{return errors.New("unknown_service_id")}; _,n,err:=net.ParseCIDR(e.CIDR);if err!=nil{return errors.New("invalid_cidr")};e.network=n;t.mu.Lock();defer t.mu.Unlock();for i:=range t.endpoints{if t.endpoints[i].ServiceID==e.ServiceID&&t.endpoints[i].CIDR==e.CIDR{t.endpoints[i]=e;return nil}};if len(t.endpoints)>=4096{return errors.New("endpoint_registry_limit")};t.endpoints=append(t.endpoints,e);return nil }
func (t *TrafficRuntime) Classify(f FlowRecord, now time.Time) (TrafficFlowObservation,bool) { dst:=net.ParseIP(strings.TrimSpace(f.DestinationIP));src:=net.ParseIP(strings.TrimSpace(f.SourceIP));if dst==nil&&src==nil{return TrafficFlowObservation{},false};t.mu.RLock();endpoints:=append([]ManagedEndpoint(nil),t.endpoints...);t.mu.RUnlock();var best *ManagedEndpoint;bestBits:=-1;for i:=range endpoints{e:=&endpoints[i];if e.network==nil||(!e.ExpiresAt.IsZero()&&!now.Before(e.ExpiresAt)){continue};matched:=(dst!=nil&&e.network.Contains(dst))||(src!=nil&&e.network.Contains(src));if !matched{continue};ones,_:=e.network.Mask.Size();if ones>bestBits{best,bestBits=e,ones}};if best==nil{return TrafficFlowObservation{},false};policy,ok:=trafficPolicyByID(best.ServiceID);if !ok{return TrafficFlowObservation{},false};return TrafficFlowObservation{FlowRecord:f,ServiceID:best.ServiceID,Class:policy.Class,Region:best.Region,Provider:best.Provider,ObservedAt:now},true }
func trafficPolicyByID(id string)(TrafficServicePolicy,bool){id=strings.TrimSpace(id);for _,p:=range DefaultTrafficServicePolicies(){if p.ID==id{return p,true}};return TrafficServicePolicy{},false}
func (t *TrafficRuntime) Ingest(flows []FlowRecord,now time.Time)int{classified:=make([]TrafficFlowObservation,0,len(flows));for _,f:=range flows{if obs,ok:=t.Classify(f,now);ok{classified=append(classified,obs)}};if len(classified)==0{return 0};t.mu.Lock();defer t.mu.Unlock();t.flows=append(t.flows,classified...);cutoff:=now.Add(-2*time.Minute);keep:=t.flows[:0];for _,f:=range t.flows{if f.ObservedAt.After(cutoff){keep=append(keep,f)}};t.flows=keep;if len(t.flows)>4096{t.flows=t.flows[len(t.flows)-4096:]};return len(classified)}
func (t *TrafficRuntime) Decisions(now time.Time,nodes []Node)[]TrafficDecision{t.mu.RLock();flows:=append([]TrafficFlowObservation(nil),t.flows...);endpoints:=append([]ManagedEndpoint(nil),t.endpoints...);controllers:=make(map[string]*TrafficPathController,len(t.controllers));for serviceID,controller:=range t.controllers{controllers[serviceID]=controller};t.mu.RUnlock();services:=make(map[string]struct{});for _,f:=range flows{if f.ServiceID!=""{services[f.ServiceID]=struct{}{}}};for _,e:=range endpoints{if e.ServiceID!=""&&(e.ExpiresAt.IsZero()||now.Before(e.ExpiresAt)){services[e.ServiceID]=struct{}{}}};byService:=make(map[string][]TrafficPathObservation,len(services));for serviceID:=range services{p,ok:=trafficPolicyByID(serviceID);if !ok{continue};qualityByPath:=make(map[string]TrafficQualityObservation);for _,q:=range t.QualitySnapshot(serviceID,now){qualityByPath[q.PathID]=q};for _,n:=range nodes{q,measured:=qualityByPath[n.ID];if !n.Healthy&&!measured{continue};o:=TrafficPathObservation{PathID:n.ID,ServiceID:serviceID,Class:p.Class,LatencyMs:n.LatencyMs,PacketLoss:n.PacketLoss,Healthy:n.Healthy,ObservedAt:now};if measured{o.LatencyMs=q.LatencyMs;o.JitterMs=q.JitterMs;o.PacketLoss=q.PacketLoss;o.Congestion=q.Congestion;o.Healthy=q.Healthy;o.ObservedAt=q.ObservedAt};if !nodeHasService(n,serviceID)&&!measured{continue};if !usableObservation(o,p){continue};byService[serviceID]=append(byService[serviceID],o)}};out:=make([]TrafficDecision,0,len(byService));for serviceID,observations:=range byService{p,ok:=trafficPolicyByID(serviceID);if !ok||len(observations)==0{continue};c:=controllers[serviceID];if c==nil{c=&TrafficPathController{};t.mu.Lock();if existing:=t.controllers[serviceID];existing!=nil{c=existing}else{t.controllers[serviceID]=c};t.mu.Unlock()};d,ok:=c.Decide(observations,p,now);if !ok{continue};t.mu.Lock();t.decisions[serviceID]=d;t.mu.Unlock();out=append(out,d)};sort.Slice(out,func(i,j int)bool{return out[i].ServiceID<out[j].ServiceID});return out}
func (a *App) trafficEndpoints(w http.ResponseWriter,r *http.Request){if !method(w,r,http.MethodPost){return};if !requirePermission(a,"network.configure",w,r){return};var e ManagedEndpoint;dec:=json.NewDecoder(http.MaxBytesReader(w,r.Body,64<<10));if err:=dec.Decode(&e);err!=nil{jsonResponse(w,400,map[string]string{"error":"invalid_json"});return};if err:=a.traffic.UpsertEndpoint(e);err!=nil{jsonResponse(w,400,map[string]string{"error":err.Error()});return};a.audit(r,"traffic.endpoint",e.ServiceID,"accepted",e);jsonResponse(w,http.StatusAccepted,map[string]any{"status":"accepted","service_id":e.ServiceID,"cidr":e.CIDR,"source":"managed-endpoint-feed"})}
func (a *App) trafficFlowIngest(w http.ResponseWriter,r *http.Request){if !method(w,r,http.MethodPost){return};if !requirePermission(a,"network.read",w,r){return};var req struct{Flows []FlowRecord `json:"flows"`};dec:=json.NewDecoder(http.MaxBytesReader(w,r.Body,4<<20));if err:=dec.Decode(&req);err!=nil{jsonResponse(w,400,map[string]string{"error":"invalid_json"});return};if len(req.Flows)>4096{jsonResponse(w,413,map[string]string{"error":"flow_batch_limit"});return};now:=time.Now().UTC();accepted:=a.traffic.Ingest(req.Flows,now);jsonResponse(w,http.StatusAccepted,map[string]any{"accepted":accepted,"received":len(req.Flows),"observed_at":now})}
func (a *App) trafficDecisions(w http.ResponseWriter,r *http.Request){if !method(w,r,http.MethodGet){return};if !requirePermission(a,"network.read",w,r){return};nodes,err:=a.loadNodes(r.Context());if err!=nil{jsonResponse(w,500,map[string]string{"error":"node_query_failed"});return};decisions:=a.traffic.Decisions(time.Now().UTC(),nodes);jsonResponse(w,200,map[string]any{"decisions":decisions,"execution":"read-only","configuration_changes":"approval-required"})}
func flowListenerConfigFromEnv()(FlowListenerConfig,bool){address:=strings.TrimSpace(os.Getenv("FTN_FLOW_LISTEN_ADDR"));exportersRaw:=strings.TrimSpace(os.Getenv("FTN_FLOW_EXPORTERS"));if address==""||exportersRaw==""{return FlowListenerConfig{},false};exporters:=make([]string,0,16);for _,item:=range strings.Split(exportersRaw,","){item=strings.TrimSpace(item);if item!=""{exporters=append(exporters,item);if len(exporters)==256{break}}};return FlowListenerConfig{Address:address,Exporters:exporters},len(exporters)>0}
