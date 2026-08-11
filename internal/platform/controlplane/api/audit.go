package api

import (
	"sync"
	"time"
)

type EventType string

const (
	EventCommandCreated  EventType = "command.created"
	EventCommandApproved EventType = "command.approved"
	EventCommandRejected EventType = "command.rejected"
	EventCommandExecuted EventType = "command.executed"
	EventCommandFailed   EventType = "command.failed"
)

type AuditEvent struct {
	ID string `json:"id"`
	Type EventType `json:"type"`
	CommandID string `json:"command_id"`
	Target string `json:"target"`
	Actor string `json:"actor"`
	Message string `json:"message,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type AuditLog struct {
	mu sync.RWMutex
	events []AuditEvent
}

func NewAuditLog() *AuditLog { return &AuditLog{} }

func (l *AuditLog) Append(e AuditEvent) bool {
	if e.ID == "" || e.Type == "" || e.CreatedAt.IsZero() { return false }
	l.mu.Lock(); defer l.mu.Unlock()
	l.events = append(l.events, e)
	return true
}

func (l *AuditLog) Snapshot() []AuditEvent {
	l.mu.RLock(); defer l.mu.RUnlock()
	out := make([]AuditEvent, len(l.events))
	copy(out, l.events)
	return out
}
