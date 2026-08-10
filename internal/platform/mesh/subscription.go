package mesh

import "sync"

type SubscriptionPolicy struct {
	Role string `json:"role"`
	AllowedNodeIDs map[string]bool `json:"allowed_node_ids,omitempty"`
	AllowedEventTypes map[EventType]bool `json:"allowed_event_types,omitempty"`
}

// AuthorizeSubscription enforces least-privilege event visibility. An empty
// node allow-list means no node is implicitly authorized.
func AuthorizeSubscription(p SubscriptionPolicy, nodeID string, eventType EventType) bool {
	if p.Role == "" || nodeID == "" { return false }
	if len(p.AllowedNodeIDs) == 0 || !p.AllowedNodeIDs[nodeID] { return false }
	if len(p.AllowedEventTypes) == 0 || !p.AllowedEventTypes[eventType] { return false }
	return true
}

type SubscriptionRegistry struct {
	mu sync.RWMutex
	policies map[string]SubscriptionPolicy
}

func NewSubscriptionRegistry() *SubscriptionRegistry { return &SubscriptionRegistry{policies: make(map[string]SubscriptionPolicy)} }

func (r *SubscriptionRegistry) Set(sessionID string, policy SubscriptionPolicy) {
	r.mu.Lock(); r.policies[sessionID] = policy; r.mu.Unlock()
}

func (r *SubscriptionRegistry) Get(sessionID string) (SubscriptionPolicy, bool) {
	r.mu.RLock(); defer r.mu.RUnlock()
	p, ok := r.policies[sessionID]
	return p, ok
}

func (r *SubscriptionRegistry) Delete(sessionID string) { r.mu.Lock(); delete(r.policies, sessionID); r.mu.Unlock() }
