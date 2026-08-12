package ws

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

const Protocol = "ftn-ws/1"

var ErrClosed = errors.New("websocket fabric closed")

type Event struct {
	ID        uint64         `json:"id"`
	Protocol  string         `json:"protocol"`
	Type      string         `json:"type"`
	Topic     string         `json:"topic"`
	Source    string         `json:"source"`
	NodeID    string         `json:"node_id,omitempty"`
	Timestamp time.Time      `json:"timestamp"`
	Data      map[string]any `json:"data,omitempty"`
}

type Subscription struct {
	ID    string
	Topic string
}

type Client interface {
	ID() string
	Send(context.Context, Event) error
}

type Fabric struct {
	mu          sync.RWMutex
	clients     map[string]Client
	subs        map[string]map[string]Subscription
	sequence    atomic.Uint64
	closed      atomic.Bool
	bufferSize  int
	deliveryTTL time.Duration
}

func NewFabric(bufferSize int, deliveryTTL time.Duration) *Fabric {
	if bufferSize < 1 { bufferSize = 64 }
	if deliveryTTL <= 0 { deliveryTTL = 5 * time.Second }
	return &Fabric{
		clients: make(map[string]Client),
		subs: make(map[string]map[string]Subscription),
		bufferSize: bufferSize,
		deliveryTTL: deliveryTTL,
	}
}

func (f *Fabric) Register(c Client) error {
	if c == nil || c.ID() == "" { return errors.New("invalid websocket client") }
	if f.closed.Load() { return ErrClosed }
	f.mu.Lock()
	defer f.mu.Unlock()
	f.clients[c.ID()] = c
	if _, ok := f.subs[c.ID()]; !ok { f.subs[c.ID()] = make(map[string]Subscription) }
	return nil
}

func (f *Fabric) Unregister(id string) {
	f.mu.Lock()
	delete(f.clients, id)
	delete(f.subs, id)
	f.mu.Unlock()
}

func (f *Fabric) Subscribe(id, subscriptionID, topic string) error {
	if subscriptionID == "" || topic == "" { return errors.New("subscription id and topic are required") }
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.clients[id]; !ok { return errors.New("client not registered") }
	if _, ok := f.subs[id]; !ok { f.subs[id] = make(map[string]Subscription) }
	f.subs[id][subscriptionID] = Subscription{ID: subscriptionID, Topic: topic}
	return nil
}

func (f *Fabric) Unsubscribe(id, subscriptionID string) {
	f.mu.Lock()
	if subscriptions, ok := f.subs[id]; ok { delete(subscriptions, subscriptionID) }
	f.mu.Unlock()
}

func (f *Fabric) Publish(ctx context.Context, event Event) error {
	if f.closed.Load() { return ErrClosed }
	if event.Protocol == "" { event.Protocol = Protocol }
	event.ID = f.sequence.Add(1)
	if event.Timestamp.IsZero() { event.Timestamp = time.Now().UTC() }

	f.mu.RLock()
	clients := make(map[string]Client, len(f.clients))
	for id, c := range f.clients {
		for _, sub := range f.subs[id] {
			if topicMatches(sub.Topic, event.Topic) { clients[id] = c; break }
		}
	}
	f.mu.RUnlock()

	for _, c := range clients {
		deliverCtx, cancel := context.WithTimeout(ctx, f.deliveryTTL)
		err := c.Send(deliverCtx, event)
		cancel()
		if err != nil { f.Unregister(c.ID()) }
	}
	return nil
}

func (f *Fabric) Close() {
	if !f.closed.CompareAndSwap(false, true) { return }
	f.mu.Lock()
	f.clients = make(map[string]Client)
	f.subs = make(map[string]map[string]Subscription)
	f.mu.Unlock()
}

func (f *Fabric) ClientCount() int {
	f.mu.RLock(); defer f.mu.RUnlock()
	return len(f.clients)
}

func topicMatches(subscription, topic string) bool {
	if subscription == "*" || subscription == topic { return true }
	if len(subscription) > 2 && subscription[len(subscription)-2:] == "/*" {
		prefix := subscription[:len(subscription)-1]
		return len(topic) >= len(prefix) && topic[:len(prefix)] == prefix
	}
	return false
}
