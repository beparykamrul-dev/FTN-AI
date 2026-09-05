package realtime

import (
	"context"
	"fmt"
	"sync"
)

const MaxBroadcastPayload = 1 << 20

type Hub struct {
	mu    sync.RWMutex
	rooms map[string]map[string]Peer
}

func NewHub() *Hub {
	return &Hub{rooms: make(map[string]map[string]Peer)}
}

func (h *Hub) Join(roomID string, peer Peer) error {
	if h == nil || roomID == "" || peer == nil || peer.ID() == "" {
		return fmt.Errorf("room and peer are required")
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.rooms == nil {
		h.rooms = make(map[string]map[string]Peer)
	}
	room := h.rooms[roomID]
	if room == nil {
		room = make(map[string]Peer)
		h.rooms[roomID] = room
	}
	room[peer.ID()] = peer
	return nil
}

func (h *Hub) Leave(roomID, peerID string) {
	if h == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	room := h.rooms[roomID]
	if room == nil {
		return
	}
	delete(room, peerID)
	if len(room) == 0 {
		delete(h.rooms, roomID)
	}
}

func (h *Hub) Broadcast(ctx context.Context, roomID, senderID string, payload []byte) map[string]error {
	result := make(map[string]error)
	if h == nil {
		result["_"] = fmt.Errorf("hub is nil")
		return result
	}
	if ctx == nil {
		result["_"] = fmt.Errorf("context is nil")
		return result
	}
	if len(payload) > MaxBroadcastPayload {
		result["_"] = fmt.Errorf("payload exceeds %d bytes", MaxBroadcastPayload)
		return result
	}

	h.mu.RLock()
	room := h.rooms[roomID]
	peers := make([]Peer, 0, len(room))
	for id, p := range room {
		if id != senderID && p != nil {
			peers = append(peers, p)
		}
	}
	h.mu.RUnlock()
	for _, p := range peers {
		result[p.ID()] = p.Send(ctx, payload)
	}
	return result
}
