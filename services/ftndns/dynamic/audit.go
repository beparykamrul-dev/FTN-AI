package dynamic

import "time"

// AuditEvent records the decision boundary around a dynamic DNS change.
type AuditEvent struct {
	Time       time.Time
	Action     string
	Name       string
	Type       string
	NodeID     string
	Version    uint64
	ApprovedBy string
	Result     string
}

// NewAuditEvent creates an immutable audit record for a planned operation.
func NewAuditEvent(action, name, typ, nodeID string, version uint64, approvedBy, result string) AuditEvent {
	return AuditEvent{
		Time:       time.Now().UTC(),
		Action:     action,
		Name:       name,
		Type:       typ,
		NodeID:     nodeID,
		Version:    version,
		ApprovedBy: approvedBy,
		Result:     result,
	}
}
