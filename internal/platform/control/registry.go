package control

import (
	"sync"
	"time"
)

type Server struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Address   string            `json:"address"`
	Platform  string            `json:"platform"`
	AgentVer  string            `json:"agent_version,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
	Online    bool              `json:"online"`
	LastSeen  time.Time         `json:"last_seen"`
}

type Registry struct {
	mu      sync.RWMutex
	servers map[string]Server
}

func NewRegistry() *Registry {
	return &Registry{servers: make(map[string]Server)}
}

func (r *Registry) Upsert(s Server) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if s.LastSeen.IsZero() { s.LastSeen = time.Now().UTC() }
	r.servers[s.ID] = s
}

func (r *Registry) SetOnline(id string, online bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	s, ok := r.servers[id]
	if !ok { return }
	s.Online = online
	s.LastSeen = time.Now().UTC()
	r.servers[id] = s
}

func (r *Registry) List() []Server {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Server, 0, len(r.servers))
	for _, s := range r.servers { out = append(out, s) }
	return out
}
