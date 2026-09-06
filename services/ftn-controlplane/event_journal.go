package controlplane

import (
	"errors"
	"strings"
	"sync"
	"time"
)

var ErrEventNotFound = errors.New("event not found")

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

type EventJournal interface {
	Append(event JournalEvent) (JournalEvent, error)
	ReadAfter(tenantID string, sequence uint64, limit int) ([]JournalEvent, error)
	CommitOffset(consumerID, tenantID string, sequence uint64) error
	Offset(consumerID, tenantID string) uint64
}

type MemoryEventJournal struct {
	mu       sync.RWMutex
	nextSeq  map[string]uint64
	events   map[string][]JournalEvent
	offsets  map[string]uint64
}

func NewMemoryEventJournal() *MemoryEventJournal {
	return &MemoryEventJournal{nextSeq: make(map[string]uint64), events: make(map[string][]JournalEvent), offsets: make(map[string]uint64)}
}

func offsetKey(consumerID, tenantID string) string {
	return strings.TrimSpace(consumerID) + "\x00" + strings.TrimSpace(tenantID)
}

func (j *MemoryEventJournal) Append(event JournalEvent) (JournalEvent, error) {
	if j == nil {
		return JournalEvent{}, errors.New("event journal is required")
	}
	event.ID = strings.TrimSpace(event.ID)
	event.TenantID = strings.TrimSpace(event.TenantID)
	event.Type = strings.TrimSpace(event.Type)
	if event.TenantID == "" || event.Type == "" || event.ID == "" {
		return JournalEvent{}, errors.New("event requires id, tenant_id and type")
	}
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.nextSeq == nil {
		j.nextSeq = make(map[string]uint64)
	}
	if j.events == nil {
		j.events = make(map[string][]JournalEvent)
	}
	for _, existing := range j.events[event.TenantID] {
		if existing.ID == event.ID {
			return existing, nil
		}
	}
	j.nextSeq[event.TenantID]++
	event.Sequence = j.nextSeq[event.TenantID]
	event.CreatedAt = time.Now().UTC()
	event.Payload = append([]byte(nil), event.Payload...)
	j.events[event.TenantID] = append(j.events[event.TenantID], event)
	return event, nil
}

func (j *MemoryEventJournal) ReadAfter(tenantID string, sequence uint64, limit int) ([]JournalEvent, error) {
	if j == nil {
		return nil, errors.New("event journal is required")
	}
	tenantID = strings.TrimSpace(tenantID)
	if tenantID == "" {
		return nil, errors.New("tenant is required")
	}
	if limit <= 0 {
		return []JournalEvent{}, nil
	}
	j.mu.RLock()
	defer j.mu.RUnlock()
	out := make([]JournalEvent, 0)
	for _, event := range j.events[tenantID] {
		if event.Sequence > sequence {
			event.Payload = append([]byte(nil), event.Payload...)
			out = append(out, event)
			if len(out) == limit {
				break
			}
		}
	}
	return out, nil
}

func (j *MemoryEventJournal) CommitOffset(consumerID, tenantID string, sequence uint64) error {
	if j == nil {
		return errors.New("event journal is required")
	}
	consumerID = strings.TrimSpace(consumerID)
	tenantID = strings.TrimSpace(tenantID)
	if consumerID == "" || tenantID == "" {
		return errors.New("consumer and tenant are required")
	}
	j.mu.Lock()
	defer j.mu.Unlock()
	key := offsetKey(consumerID, tenantID)
	if sequence < j.offsets[key] {
		return errors.New("consumer offset cannot move backwards")
	}
	if sequence > j.nextSeq[tenantID] {
		return ErrEventNotFound
	}
	if j.offsets == nil {
		j.offsets = make(map[string]uint64)
	}
	j.offsets[key] = sequence
	return nil
}

func (j *MemoryEventJournal) Offset(consumerID, tenantID string) uint64 {
	if j == nil {
		return 0
	}
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.offsets[offsetKey(consumerID, tenantID)]
}
