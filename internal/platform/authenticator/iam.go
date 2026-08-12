package authenticator

import (
	"errors"
	"sync"
	"time"
)

// PrincipalType identifies an FTN identity without tying it to a human account.
type PrincipalType string

const (
	PrincipalUser   PrincipalType = "user"
	PrincipalService PrincipalType = "service"
	PrincipalAgent  PrincipalType = "agent"
)

type Principal struct {
	ID        string
	Type      PrincipalType
	Enabled   bool
	CreatedAt time.Time
}

type Role struct {
	ID           string
	Capabilities []string
}

type Policy struct {
	ID         string
	Effect     string // allow or deny
	Capability string
	Resource   string
}

// IAMRegistry is an in-process control-plane boundary. Persistence and distributed
// replication are intentionally injected later; authentication secrets do not live here.
type IAMRegistry struct {
	mu          sync.RWMutex
	principals  map[string]Principal
	roles       map[string]Role
	assignments map[string]map[string]struct{}
	policies    map[string]Policy
}

func NewIAMRegistry() *IAMRegistry {
	return &IAMRegistry{
		principals:  make(map[string]Principal),
		roles:       make(map[string]Role),
		assignments: make(map[string]map[string]struct{}),
		policies:    make(map[string]Policy),
	}
}

func (r *IAMRegistry) UpsertPrincipal(p Principal) error {
	if p.ID == "" || (p.Type != PrincipalUser && p.Type != PrincipalService && p.Type != PrincipalAgent) {
		return errors.New("invalid principal")
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now().UTC()
	}
	r.mu.Lock()
	r.principals[p.ID] = p
	r.mu.Unlock()
	return nil
}

func (r *IAMRegistry) UpsertRole(role Role) error {
	if role.ID == "" {
		return errors.New("invalid role")
	}
	r.mu.Lock()
	role.Capabilities = append([]string(nil), role.Capabilities...)
	r.roles[role.ID] = role
	r.mu.Unlock()
	return nil
}

func (r *IAMRegistry) AssignRole(principalID, roleID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.principals[principalID]; !ok {
		return errors.New("principal not found")
	}
	if _, ok := r.roles[roleID]; !ok {
		return errors.New("role not found")
	}
	if r.assignments[principalID] == nil {
		r.assignments[principalID] = make(map[string]struct{})
	}
	r.assignments[principalID][roleID] = struct{}{}
	return nil
}

func (r *IAMRegistry) UpsertPolicy(policy Policy) error {
	if policy.ID == "" || policy.Capability == "" || policy.Resource == "" {
		return errors.New("invalid policy")
	}
	if policy.Effect != "allow" && policy.Effect != "deny" {
		return errors.New("invalid policy effect")
	}
	r.mu.Lock()
	r.policies[policy.ID] = policy
	r.mu.Unlock()
	return nil
}

// Capabilities returns the effective role-derived capability set. Policy evaluation
// remains a separate boundary so deny/condition logic cannot be accidentally bypassed.
func (r *IAMRegistry) Capabilities(principalID string) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.principals[principalID]
	if !ok {
		return nil, errors.New("principal not found")
	}
	if !p.Enabled {
		return nil, errors.New("principal disabled")
	}
	seen := make(map[string]struct{})
	for roleID := range r.assignments[principalID] {
		for _, capability := range r.roles[roleID].Capabilities {
			seen[capability] = struct{}{}
		}
	}
	out := make([]string, 0, len(seen))
	for capability := range seen {
		out = append(out, capability)
	}
	return out, nil
}
