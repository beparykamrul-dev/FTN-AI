package realtime

import (
	"context"
	"fmt"
	"sync"
)

type Peer interface {
	ID() string
	Send(context.Context, []byte) error
	Close() error
}

type Transport struct {
	mu    sync.RWMutex
	peers map[string]Peer
}

func NewTransport() *Transport {
	return &Transport{peers: make(map[string]Peer)}
}

func (t *Transport) Register(peer Peer) error {
	if peer == nil || peer.ID() == "" {
		return fmt.Errorf("invalid peer")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, exists := t.peers[peer.ID()]; exists {
		return fmt.Errorf("peer already registered: %s", peer.ID())
	}
	t.peers[peer.ID()] = peer
	return nil
}

func (t *Transport) Unregister(id string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	peer, ok := t.peers[id]
	if !ok {
		return nil
	}
	delete(t.peers, id)
	return peer.Close()
}

func (t *Transport) Send(ctx context.Context, id string, payload []byte) error {
	t.mu.RLock()
	peer, ok := t.peers[id]
	t.mu.RUnlock()
	if !ok {
		return fmt.Errorf("peer not connected: %s", id)
	}
	return peer.Send(ctx, payload)
}
