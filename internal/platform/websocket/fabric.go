package websocket

import (
	"context"
	"errors"
	"sync"
	"time"
)

const ProtocolVersion = "ftn-ws/1"

var ErrClosed = errors.New("ftn websocket fabric is closed")

// Event is the canonical realtime envelope shared by Control, Monitoring and
// future FTN services. Payload remains opaque to the transport layer.
type Event struct {
	Version   string      `json:"version"`
	EventID   string      `json:"event_id"`
	Type      string      `json:"type"`
	Source    string      `json:"source"`
	NodeID    string      `json:"node_id,omitempty"`
	Sequence  uint64      `json:"sequence"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data,omitempty"`
}

// Subscriber receives ordered events for its subscriptions.
type Subscriber struct {
	ID     string
	Events chan Event
	Topics map[string]struct{}
}

// Fabric is the in-process event core. A transport adapter can expose it over
// WebSocket, while the core remains independent of HTTP/WebSocket libraries.
type Fabric struct {
	mu          sync.RWMutex
	subscribers map[string]*Subscriber
	sequence    uint64
	closed      bool
}

func NewFabric() *Fabric {
	return &Fabric{subscribers: make(map[string]*Subscriber)}
}

func (f *Fabric) Subscribe(id string, topics ...string) (*Subscriber, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.closed {
		return nil, ErrClosed
	}
	if id == "" {
		return nil, errors.New("subscriber id is required")
	}
	topicSet := make(map[string]struct{}, len(topics))
	for _, topic := range topics {
		if topic != "" {
			topicSet[topic] = struct{}{}
		}
	}
	s := &Subscriber{ID: id, Events: make(chan Event, 128), Topics: topicSet}
	f.subscribers[id] = s
	return s, nil
}

func (f *Fabric) Unsubscribe(id string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if s, ok := f.subscribers[id]; ok {
		delete(f.subscribers, id)
		close(s.Events)
	}
}

func (f *Fabric) Publish(ctx context.Context, event Event) error {
	f.mu.Lock()
	if f.closed {
		f.mu.Unlock()
		return ErrClosed
	}
	f.sequence++
	event.Version = ProtocolVersion
	event.Sequence = f.sequence
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}
	for _, s := range f.subscribers {
		if !matches(s.Topics, event.Type) {
			continue
		}
		select {
		case s.Events <- event:
		default:
			// Slow consumers are not allowed to block the fabric. The transport
			// layer can reconnect/replay from the sequence store later.
		}
	}
	f.mu.Unlock()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (f *Fabric) Close() {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.closed {
		return
	}
	f.closed = true
	for id, s := range f.subscribers {
		delete(f.subscribers, id)
		close(s.Events)
	}
}

func matches(topics map[string]struct{}, eventType string) bool {
	if len(topics) == 0 {
		return true
	}
	if _, ok := topics[eventType]; ok {
		return true
	}
	for topic := range topics {
		if len(topic) < len(eventType) && eventType[:len(topic)] == topic && eventType[len(topic)] == '.' {
			return true
		}
	}
	return false
}
