package authenticator

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

var ErrRotationInProgress = errors.New("key rotation already in progress")
var ErrRotationQuorum = errors.New("key rotation quorum not reached")

type RotationNode struct { NodeID string; Version uint64; Accepted bool; VerifiedAt time.Time; Healthy bool }
type RotationPlan struct { KeyID string; OldVersion uint64; NewVersion uint64; StartedAt time.Time; GraceUntil time.Time; RequiredHealthy int; RequiredVerified int; Nodes map[string]RotationNode; Committed bool }
type KeyRotationCoordinator struct { mu sync.RWMutex; plans map[string]*RotationPlan }

func NewKeyRotationCoordinator() *KeyRotationCoordinator { return &KeyRotationCoordinator{plans: make(map[string]*RotationPlan)} }

func (c *KeyRotationCoordinator) Start(keyID string, oldVersion, newVersion uint64, graceUntil time.Time, requiredHealthy, requiredVerified int) error {
	if c == nil { return errors.New("key rotation coordinator is required") }
	keyID = strings.TrimSpace(keyID)
	if keyID == "" || newVersion == 0 || newVersion == oldVersion || requiredHealthy < 1 || requiredVerified < 1 || graceUntil.IsZero() { return errors.New("invalid key rotation plan") }
	if !graceUntil.After(time.Now()) { return errors.New("rotation grace period must be in the future") }
	c.mu.Lock(); defer c.mu.Unlock()
	if c.plans == nil { c.plans = make(map[string]*RotationPlan) }
	if _, exists := c.plans[keyID]; exists { return ErrRotationInProgress }
	c.plans[keyID] = &RotationPlan{KeyID:keyID, OldVersion:oldVersion, NewVersion:newVersion, StartedAt:time.Now().UTC(), GraceUntil:graceUntil.UTC(), RequiredHealthy:requiredHealthy, RequiredVerified:requiredVerified, Nodes:make(map[string]RotationNode)}
	return nil
}

func (c *KeyRotationCoordinator) ReportNode(keyID string, node RotationNode) error {
	if c == nil { return errors.New("key rotation coordinator is required") }
	keyID = strings.TrimSpace(keyID); node.NodeID = strings.TrimSpace(node.NodeID)
	if keyID == "" || node.NodeID == "" || node.Version == 0 { return errors.New("invalid rotation node report") }
	c.mu.Lock(); defer c.mu.Unlock()
	plan, ok := c.plans[keyID]; if !ok { return errors.New("rotation plan not found") }
	if node.Version != plan.NewVersion { return errors.New("rotation node version does not match target version") }
	if !node.VerifiedAt.IsZero() { node.VerifiedAt = node.VerifiedAt.UTC() }
	plan.Nodes[node.NodeID] = node
	return nil
}

func (c *KeyRotationCoordinator) Commit(keyID string) error {
	if c == nil { return errors.New("key rotation coordinator is required") }
	keyID = strings.TrimSpace(keyID); if keyID == "" { return errors.New("key id is required") }
	c.mu.Lock(); defer c.mu.Unlock()
	plan, ok := c.plans[keyID]; if !ok { return errors.New("rotation plan not found") }
	if plan.Committed { return nil }
	healthy, verified := 0, 0
	for _, node := range plan.Nodes { if node.Healthy { healthy++ }; if node.Healthy && node.Version == plan.NewVersion && node.Accepted && !node.VerifiedAt.IsZero() { verified++ } }
	if healthy < plan.RequiredHealthy || verified < plan.RequiredVerified { return ErrRotationQuorum }
	plan.Committed = true
	return nil
}

func (c *KeyRotationCoordinator) Get(keyID string) (RotationPlan, bool) {
	if c == nil { return RotationPlan{}, false }
	c.mu.RLock(); defer c.mu.RUnlock()
	plan, ok := c.plans[strings.TrimSpace(keyID)]; if !ok { return RotationPlan{}, false }
	copyPlan := *plan
	copyPlan.Nodes = make(map[string]RotationNode, len(plan.Nodes))
	ids := make([]string, 0, len(plan.Nodes)); for id := range plan.Nodes { ids = append(ids, id) }; sort.Strings(ids)
	for _, id := range ids { copyPlan.Nodes[id] = plan.Nodes[id] }
	return copyPlan, true
}
