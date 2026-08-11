package mesh

import (
	"sort"
	"sync"
	"time"
)

type RouteState struct {
	Destination string    `json:"destination"`
	NextHop     string    `json:"next_hop"`
	Metric      uint32    `json:"metric"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RouteTable struct {
	mu     sync.RWMutex
	routes map[string]RouteState
}

func NewRouteTable() *RouteTable { return &RouteTable{routes: make(map[string]RouteState)} }

func (r *RouteTable) Upsert(route RouteState) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if route.UpdatedAt.IsZero() { route.UpdatedAt = time.Now().UTC() }
	r.routes[route.Destination] = route
}

func (r *RouteTable) Get(destination string) (RouteState, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.routes[destination]
	return v, ok
}

func (r *RouteTable) Snapshot() []RouteState {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]RouteState, 0, len(r.routes))
	for _, v := range r.routes { out = append(out, v) }
	sort.Slice(out, func(i, j int) bool { return out[i].Destination < out[j].Destination })
	return out
}
