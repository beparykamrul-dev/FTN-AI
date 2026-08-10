package mesh

import (
	"sync"
	"time"
)

type PeerHealth struct {
	PeerID string `json:"peer_id"`
	State LinkState `json:"state"`
	Score uint8 `json:"score"`
	LastHeartbeat time.Time `json:"last_heartbeat"`
	ConsecutiveMisses uint32 `json:"consecutive_misses"`
}

type HealthPolicy struct {
	HeartbeatTimeout time.Duration
	DegradedAfterMisses uint32
	DownAfterMisses uint32
}

func DefaultHealthPolicy() HealthPolicy {
	return HealthPolicy{HeartbeatTimeout: 15 * time.Second, DegradedAfterMisses: 2, DownAfterMisses: 4}
}

type HealthRegistry struct {
	mu sync.RWMutex
	peers map[string]PeerHealth
	policy HealthPolicy
}

func NewHealthRegistry(policy HealthPolicy) *HealthRegistry {
	if policy.HeartbeatTimeout <= 0 { policy = DefaultHealthPolicy() }
	if policy.DegradedAfterMisses == 0 { policy.DegradedAfterMisses = 2 }
	if policy.DownAfterMisses <= policy.DegradedAfterMisses { policy.DownAfterMisses = 4 }
	return &HealthRegistry{peers: make(map[string]PeerHealth), policy: policy}
}

func (r *HealthRegistry) Observe(peerID string, now time.Time, score uint8) PeerHealth {
	if score > 100 { score = 100 }
	r.mu.Lock(); defer r.mu.Unlock()
	p := PeerHealth{PeerID: peerID, State: LinkUp, Score: score, LastHeartbeat: now.UTC()}
	r.peers[peerID] = p
	return p
}

func (r *HealthRegistry) Evaluate(now time.Time) []PeerHealth {
	r.mu.Lock(); defer r.mu.Unlock()
	out := make([]PeerHealth, 0, len(r.peers))
	for id, p := range r.peers {
		if now.Sub(p.LastHeartbeat) > r.policy.HeartbeatTimeout {
			p.ConsecutiveMisses++
			if p.ConsecutiveMisses >= r.policy.DownAfterMisses { p.State = LinkDown
			} else if p.ConsecutiveMisses >= r.policy.DegradedAfterMisses { p.State = LinkDegraded }
		}
		r.peers[id] = p
		out = append(out, p)
	}
	return out
}

func (r *HealthRegistry) Get(peerID string) (PeerHealth, bool) {
	r.mu.RLock(); defer r.mu.RUnlock()
	p, ok := r.peers[peerID]
	return p, ok
}
