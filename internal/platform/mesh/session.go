package mesh

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"sync"
	"time"
)

type Session struct {
	ID            string             `json:"id"`
	NodeID        string             `json:"node_id"`
	ConnectedAt   time.Time          `json:"connected_at"`
	LastHeartbeat time.Time          `json:"last_heartbeat"`
	Subscriptions map[EventType]bool `json:"subscriptions"`
}

type SessionRegistry struct {
	mu       sync.RWMutex
	sessions map[string]Session
}

func NewSessionRegistry() *SessionRegistry { return &SessionRegistry{sessions: make(map[string]Session)} }

func newSessionID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil { return hex.EncodeToString([]byte(time.Now().UTC().String())) }
	return hex.EncodeToString(b)
}

func (r *SessionRegistry) Open(nodeID string) Session {
	nodeID = strings.TrimSpace(nodeID)
	if r == nil || nodeID == "" { return Session{} }
	now := time.Now().UTC()
	s := Session{ID:newSessionID(), NodeID:nodeID, ConnectedAt:now, LastHeartbeat:now, Subscriptions:map[EventType]bool{}}
	r.mu.Lock()
	if r.sessions == nil { r.sessions = make(map[string]Session) }
	r.sessions[s.ID] = s
	r.mu.Unlock()
	return s
}

func (r *SessionRegistry) Heartbeat(id string) bool {
	if r == nil { return false }
	id = strings.TrimSpace(id); if id == "" { return false }
	r.mu.Lock(); defer r.mu.Unlock()
	s, ok := r.sessions[id]; if !ok { return false }
	s.LastHeartbeat = time.Now().UTC(); r.sessions[id] = s; return true
}

func (r *SessionRegistry) Subscribe(id string, eventType EventType) bool {
	if r == nil { return false }
	id = strings.TrimSpace(id); if id == "" { return false }
	r.mu.Lock(); defer r.mu.Unlock()
	s, ok := r.sessions[id]; if !ok { return false }
	if s.Subscriptions == nil { s.Subscriptions = map[EventType]bool{} }
	s.Subscriptions[eventType] = true; r.sessions[id] = s; return true
}

func (r *SessionRegistry) Close(id string) {
	if r == nil { return }
	id = strings.TrimSpace(id); if id == "" { return }
	r.mu.Lock(); delete(r.sessions, id); r.mu.Unlock()
}
