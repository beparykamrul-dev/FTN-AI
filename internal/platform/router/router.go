package router

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

var ErrNoRoute = errors.New("ftn router: no healthy route")

type Service struct { ID string; Name string; Endpoint string; Region string; GatewayID string; Priority int; Healthy bool; RTT time.Duration; Loss float64; UpdatedAt time.Time; Attributes map[string]string }
type Route struct { ServiceID string; Endpoint string; GatewayID string; Region string; RTT time.Duration }
type Registry interface { Services(context.Context, string) ([]Service, error) }
type cachedRoute struct { route Route; at time.Time }
type Router struct { registry Registry; mu sync.RWMutex; cache map[string]cachedRoute; cacheTTL time.Duration }

func New(reg Registry, cacheTTL time.Duration) *Router { if cacheTTL <= 0 { cacheTTL = 5*time.Second }; return &Router{registry:reg, cache:make(map[string]cachedRoute), cacheTTL:cacheTTL} }

func (r *Router) Resolve(ctx context.Context, service, preferredRegion string) (Route, error) {
	if r == nil || r.registry == nil { return Route{}, ErrNoRoute }
	if ctx == nil { return Route{}, errors.New("context is required") }
	service, preferredRegion = strings.TrimSpace(service), strings.TrimSpace(preferredRegion)
	if service == "" { return Route{}, errors.New("service is required") }
	if err := ctx.Err(); err != nil { return Route{}, err }
	cacheKey := service + "\x00" + preferredRegion; now := time.Now()
	r.mu.RLock(); cached, ok := r.cache[cacheKey]; r.mu.RUnlock()
	if ok && now.Sub(cached.at) >= 0 && now.Sub(cached.at) < r.cacheTTL { return cached.route, nil }
	items, err := r.registry.Services(ctx, service); if err != nil { return Route{}, err }
	candidates := make([]Service, 0, len(items))
	for _, s := range items {
		s.ID, s.Endpoint, s.Region = strings.TrimSpace(s.ID), strings.TrimSpace(s.Endpoint), strings.TrimSpace(s.Region)
		if s.Healthy && s.Endpoint != "" && s.RTT >= 0 && s.Loss >= 0 && s.Loss <= 1 { candidates = append(candidates, s) }
	}
	if len(candidates) == 0 { return Route{}, ErrNoRoute }
	sort.SliceStable(candidates, func(i,j int) bool {
		a,b := candidates[i], candidates[j]; ai,bi := a.Region==preferredRegion,b.Region==preferredRegion
		if ai != bi { return ai }
		sa := float64(a.RTT.Milliseconds()) + a.Loss*1000 + float64(a.Priority); sb := float64(b.RTT.Milliseconds()) + b.Loss*1000 + float64(b.Priority)
		if sa != sb { return sa < sb }; return a.ID < b.ID
	})
	best := candidates[0]; out := Route{ServiceID:best.ID, Endpoint:best.Endpoint, GatewayID:best.GatewayID, Region:best.Region, RTT:best.RTT}
	r.mu.Lock(); r.cache[cacheKey] = cachedRoute{route:out, at:now}; r.mu.Unlock(); return out,nil
}
