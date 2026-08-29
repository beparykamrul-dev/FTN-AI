package realtime

import (
	"sync"
)

type Subscriber struct {
	ID    string
	Scope string
	Send  func(Event) error
}

type Hub struct {
	mu   sync.RWMutex
	subs map[string]map[string]Subscriber
}

func NewHub() *Hub {
	return &Hub{subs: make(map[string]map[string]Subscriber)}
}

func (h *Hub) Subscribe(s Subscriber, channel string) bool {
	if s.ID == "" || s.Send == nil || !IsAllowedChannel(channel) {
		return false
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.subs[channel] == nil {
		h.subs[channel] = make(map[string]Subscriber)
	}
	h.subs[channel][s.ID] = s
	return true
}

func (h *Hub) Unsubscribe(id, channel string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if clients := h.subs[channel]; clients != nil {
		delete(clients, id)
		if len(clients) == 0 {
			delete(h.subs, channel)
		}
	}
}

func (h *Hub) Publish(channel string, event Event) {
	if !IsAllowedChannel(channel) {
		return
	}
	h.mu.RLock()
	clients := make([]Subscriber, 0, len(h.subs[channel]))
	for _, s := range h.subs[channel] {
		clients = append(clients, s)
	}
	h.mu.RUnlock()
	for _, s := range clients {
		_ = s.Send(event)
	}
}
