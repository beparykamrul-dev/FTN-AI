package agent

import (
	"context"
	"time"
)

type AuditEvent struct {
	Time       time.Time
	Principal  string
	Scope      Scope
	Category   Category
	LayerID    string
	Capability string
	Outcome    string
}

type AuditSink interface {
	Write(context.Context, AuditEvent) error
}

// Audit records routing decisions and outcomes without storing model prompts by default.
type Audit struct {
	Sink AuditSink
	Now  func() time.Time
}

func (a *Audit) Record(ctx context.Context, principal string, scope Scope, category Category, layerID, capability, outcome string) error {
	if a == nil || a.Sink == nil {
		return nil
	}
	now := time.Now()
	if a.Now != nil {
		now = a.Now()
	}
	return a.Sink.Write(ctx, AuditEvent{
		Time:       now.UTC(),
		Principal:  principal,
		Scope:      scope,
		Category:   category,
		LayerID:    layerID,
		Capability: capability,
		Outcome:    outcome,
	})
}
