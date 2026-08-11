package controlplane

import (
	"sync"
	"time"
)

type SessionState string

const (
	SessionConnecting SessionState = "connecting"
	SessionAuthenticated SessionState = "authenticated"
	SessionClosed SessionState = "closed"
)

type WSSession struct {
	ID string `json:"id"`
	Principal string `json:"principal"`
	State SessionState `json:"state"`
	ConnectedAt time.Time `json:"connected_at"`
	LastSeen time.Time `json:"last_seen"`
}

type SessionRegistry struct {
	mu sync.RWMutex
	sessions map[string]WSSession
}

func NewSessionRegistry() *SessionRegistry { return &SessionRegistry{sessions: make(map[string]WSSession)} }

func (r *SessionRegistry) Upsert(s WSSession) {
	if s.LastSeen.IsZero() { s.LastSeen = time.Now().UTC() }
	r.mu.Lock(); r.sessions[s.ID] = s; r.mu.Unlock()
}

func (r *SessionRegistry) Get(id string) (WSSession, bool) {
	r.mu.RLock(); defer r.mu.RUnlock(); s, ok := r.sessions[id]; return s, ok
}

func (r *SessionRegistry) Remove(id string) { r.mu.Lock(); delete(r.sessions, id); r.mu.Unlock() }
