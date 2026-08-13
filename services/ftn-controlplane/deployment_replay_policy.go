package controlplane

import "sync"

// DeploymentReplayPolicy prevents the same deployment nonce from being
// accepted twice by a running control-plane process.
type DeploymentReplayPolicy struct {
	mu   sync.Mutex
	seen map[string]struct{}
}

func NewDeploymentReplayPolicy() *DeploymentReplayPolicy {
	return &DeploymentReplayPolicy{seen: make(map[string]struct{})}
}

func (p *DeploymentReplayPolicy) Accept(nonce string) bool {
	if nonce == "" {
		return false
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.seen[nonce]; ok {
		return false
	}
	p.seen[nonce] = struct{}{}
	return true
}
