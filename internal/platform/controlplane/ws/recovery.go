package ws

import (
	"sync"
	"time"
)

// RecoveryBuffer is a bounded, in-memory stream used to bridge short network
// gaps. It is intentionally not a source of truth; callers must resync from
// authoritative state when a requested sequence is no longer available.
type RecoveryBuffer struct {
	mu      sync.RWMutex
	maxSize int
	events  []Event
}

func NewRecoveryBuffer(maxSize int) *RecoveryBuffer {
	if maxSize < 1 { maxSize = 256 }
	return &RecoveryBuffer{maxSize: maxSize, events: make([]Event, 0, maxSize)}
}

func (b *RecoveryBuffer) Append(event Event) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if event.Timestamp.IsZero() { event.Timestamp = time.Now().UTC() }
	b.events = append(b.events, event)
	if excess := len(b.events) - b.maxSize; excess > 0 {
		b.events = append([]Event(nil), b.events[excess:]...)
	}
}

// Since returns events strictly after the supplied sequence. ok=false means
// the gap cannot be completely recovered from this bounded buffer.
func (b *RecoveryBuffer) Since(sequence uint64) (events []Event, ok bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if len(b.events) == 0 { return nil, sequence == 0 }
	first := b.events[0].ID
	last := b.events[len(b.events)-1].ID
	if sequence >= last { return nil, true }
	if sequence+1 < first { return nil, false }
	start := 0
	for start < len(b.events) && b.events[start].ID <= sequence { start++ }
	out := make([]Event, len(b.events)-start)
	copy(out, b.events[start:])
	return out, true
}

func (b *RecoveryBuffer) Len() int {
	b.mu.RLock(); defer b.mu.RUnlock()
	return len(b.events)
}
