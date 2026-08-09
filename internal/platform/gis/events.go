package gis

import (
	"encoding/json"
	"sync"
)

type Event struct {
	Type string `json:"type"`
	Node *MapNode `json:"node,omitempty"`
	Edge *MapEdge `json:"edge,omitempty"`
}

type Subscriber chan []byte

type Hub struct {
	mu   sync.RWMutex
	subs map[Subscriber]struct{}
}

func NewHub() *Hub { return &Hub{subs: make(map[Subscriber]struct{})} }

func (h *Hub) Subscribe(buffer int) Subscriber {
	if buffer < 1 { buffer = 16 }
	ch := make(Subscriber, buffer)
	h.mu.Lock(); h.subs[ch] = struct{}{}; h.mu.Unlock()
	return ch
}

func (h *Hub) Unsubscribe(ch Subscriber) {
	h.mu.Lock()
	if _, ok := h.subs[ch]; ok { delete(h.subs, ch); close(ch) }
	h.mu.Unlock()
}

func (h *Hub) Publish(e Event) error {
	data, err := json.Marshal(e)
	if err != nil { return err }
	h.mu.RLock(); defer h.mu.RUnlock()
	for ch := range h.subs {
		select { case ch <- data: default: }
	}
	return nil
}
