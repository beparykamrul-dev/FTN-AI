package controlplane

import "sync"

type SessionReplayGuard struct {
	mu sync.Mutex
	seen map[string]struct{}
}

func NewSessionReplayGuard() *SessionReplayGuard { return &SessionReplayGuard{seen: make(map[string]struct{})} }

func (g *SessionReplayGuard) Accept(nonce string) bool {
	if nonce == "" { return false }
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.seen[nonce]; ok { return false }
	g.seen[nonce] = struct{}{}
	return true
}
