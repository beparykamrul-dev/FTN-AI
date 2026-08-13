package realtime

import (
	"context"
	"fmt"
	"sync"
)

// Hub routes a payload to all currently connected peers in a room.
// It intentionally contains no provider-specific WebSocket implementation.
type Hub struct {
	mu    sync.RWMutex
	rooms map[string]map[string]Peer
}

func NewHub() *Hub {
	return &Hub{rooms: make(map[string]map[string]Peer)}
}

func (h *Hub) Join(roomID string, peer Peer) error {
	if roomID == "" || peer == nil || peer.ID() == "" {
		return fmt.Errorf("room and peer are required")
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	room := h.rooms[roomID]
	if room == nil {
		room = make(map[string]Peer)
		h.rooms[roomID] = room
	}
	room[peer.ID()] = peer
	return nil
}

func (h *Hub) Leave(roomID, peerID string) {
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
	h.mu.RLock()
	room := h.rooms[roomID]
	peers := make([]Peer, 0, len(room))
	for id, peer := range room {
		if id != senderID {
			peers = append(peers, peer)
		}
	}
	h.mu.RUnlock()

	result := make(map[string]error, len(peers))
	for _, peer := range peers {
		result[peer.ID()] = peer.Send(ctx, payload)
	}
	return result
}
