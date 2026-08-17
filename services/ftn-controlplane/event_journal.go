package controlplane

import (
	"errors"
	"sync"
	"time"
)

var ErrEventNotFound = errors.New("event not found")

// JournalEvent is an immutable event envelope used for durable control-plane coordination.
type JournalEvent struct {
	ID            string    `json:"id"`
	TenantID      string    `json:"tenant_id"`
	Type          string    `json:"type"`
	Sequence      uint64    `json:"sequence"`
	CorrelationID string    `json:"correlation_id"`
	CausationID   string    `json:"causation_id"`
	AggregateID   string    `json:"aggregate_id"`
	Payload       []byte    `json:"payload"`
	CreatedAt     time.Time `json:"created_at"`
}

// EventJournal defines append and replay semantics. Production adapters should
// back this interface with a transactional durable database or journal.
type EventJournal interface {
	Append(event JournalEvent) (JournalEvent, error)
	ReadAfter(tenantID string, sequence uint64, limit int) ([]JournalEvent, error)
	CommitOffset(consumerID, tenantID string, sequence uint64) error
	Offset(consumerID, tenantID string) uint64
}

// MemoryEventJournal provides deterministic semantics for tests and local development.
type MemoryEventJournal struct {
	mu      sync.RWMutex
	nextSeq map[string]uint64
	events  map[string][]JournalEvent
	offsets map[string]uint64
}

func NewMemoryEventJournal() *MemoryEventJournal {
	return &MemoryEventJournal{
		nextSeq: make(map[string]uint64),
		events:  make(map[string][]JournalEvent),
		offsets: make(map[string]uint64),
	}
}

func offsetKey(consumerID, tenantID string) string { return consumerID + "\x00" + tenantID }

func (j *MemoryEventJournal) Append(event JournalEvent) (JournalEvent, error) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if event.TenantID == "" || event.Type == "" || event.ID == "" {
		return JournalEvent{}, errors.New("event requires id, tenant_id and type")
	}
	j.nextSeq[event.TenantID]++
	event.Sequence = j.nextSeq[event.TenantID]
	event.CreatedAt = time.Now().UTC()
	event.Payload = append([]byte(nil), event.Payload...)
	j.events[event.TenantID] = append(j.events[event.TenantID], event)
	return event, nil
}

func (j *MemoryEventJournal) ReadAfter(tenantID string, sequence uint64, limit int) ([]JournalEvent, error) {
	j.mu.RLock()
	defer j.mu.RUnlock()
	if limit <= 0 {
		return []JournalEvent{}, nil
	}
	var out []JournalEvent
	for _, event := range j.events[tenantID] {
		if event.Sequence > sequence {
			out = append(out, event)
			if len(out) == limit {
				break
			}
		}
	}
	return out, nil
}

func (j *MemoryEventJournal) CommitOffset(consumerID, tenantID string, sequence uint64) error {
	j.mu.Lock()
	defer j.mu.Unlock()
	key := offsetKey(consumerID, tenantID)
	if sequence < j.offsets[key] {
		return errors.New("consumer offset cannot move backwards")
	}
	if sequence > j.nextSeq[tenantID] {
		return ErrEventNotFound
	}
	j.offsets[key] = sequence
	return nil
}

func (j *MemoryEventJournal) Offset(consumerID, tenantID string) uint64 {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.offsets[offsetKey(consumerID, tenantID)]
}
