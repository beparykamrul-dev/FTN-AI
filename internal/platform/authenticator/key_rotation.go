package authenticator

import (
	"errors"
	"sync"
	"time"
)

var ErrRotationInProgress = errors.New("key rotation already in progress")
var ErrRotationQuorum = errors.New("key rotation quorum not reached")

// RotationNode reports whether an authenticator node has accepted and verified a key version.
type RotationNode struct {
	NodeID        string
	Version       uint64
	Accepted      bool
	VerifiedAt    time.Time
	Healthy       bool
}

// RotationPlan describes a coordinated signing-key rollout.
type RotationPlan struct {
	KeyID             string
	OldVersion        uint64
	NewVersion        uint64
	StartedAt         time.Time
	GraceUntil        time.Time
	RequiredHealthy   int
	RequiredVerified  int
	Nodes             map[string]RotationNode
	Committed         bool
}

// KeyRotationCoordinator coordinates rollout state without owning private key material.
type KeyRotationCoordinator struct {
	mu    sync.RWMutex
	plans map[string]*RotationPlan
}

func NewKeyRotationCoordinator() *KeyRotationCoordinator {
	return &KeyRotationCoordinator{plans: make(map[string]*RotationPlan)}
}

func (c *KeyRotationCoordinator) Start(keyID string, oldVersion, newVersion uint64, graceUntil time.Time, requiredHealthy, requiredVerified int) error {
	if keyID == "" || newVersion == 0 || newVersion == oldVersion || requiredHealthy < 1 || requiredVerified < 1 {
		return errors.New("invalid key rotation plan")
	}
	if graceUntil.IsZero() {
		return errors.New("rotation grace period is required")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.plans[keyID]; exists {
		return ErrRotationInProgress
	}
	c.plans[keyID] = &RotationPlan{
		KeyID:            keyID,
		OldVersion:       oldVersion,
		NewVersion:       newVersion,
		StartedAt:        time.Now().UTC(),
		GraceUntil:       graceUntil,
		RequiredHealthy:  requiredHealthy,
		RequiredVerified: requiredVerified,
		Nodes:            make(map[string]RotationNode),
	}
	return nil
}

func (c *KeyRotationCoordinator) ReportNode(keyID string, node RotationNode) error {
	if keyID == "" || node.NodeID == "" || node.Version == 0 {
		return errors.New("invalid rotation node report")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	plan, ok := c.plans[keyID]
	if !ok {
		return errors.New("rotation plan not found")
	}
	plan.Nodes[node.NodeID] = node
	return nil
}

func (c *KeyRotationCoordinator) Commit(keyID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	plan, ok := c.plans[keyID]
	if !ok {
		return errors.New("rotation plan not found")
	}
	if plan.Committed {
		return nil
	}
	healthy, verified := 0, 0
	for _, node := range plan.Nodes {
		if node.Healthy {
			healthy++
		}
		if node.Healthy && node.Version == plan.NewVersion && node.Accepted && !node.VerifiedAt.IsZero() {
			verified++
		}
	}
	if healthy < plan.RequiredHealthy || verified < plan.RequiredVerified {
		return ErrRotationQuorum
	}
	plan.Committed = true
	return nil
}

func (c *KeyRotationCoordinator) Get(keyID string) (RotationPlan, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	plan, ok := c.plans[keyID]
	if !ok {
		return RotationPlan{}, false
	}
	copyPlan := *plan
	copyPlan.Nodes = make(map[string]RotationNode, len(plan.Nodes))
	for id, node := range plan.Nodes {
		copyPlan.Nodes[id] = node
	}
	return copyPlan, true
}
