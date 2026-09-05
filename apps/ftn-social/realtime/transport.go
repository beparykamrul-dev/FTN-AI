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

const MaxTransportPayload = 1 << 20

func NewTransport() *Transport {
	return &Transport{peers: make(map[string]Peer)}
}

func (t *Transport) Register(peer Peer) error {
	if t == nil || peer == nil || peer.ID() == "" {
		return fmt.Errorf("invalid peer")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.peers == nil {
		t.peers = make(map[string]Peer)
	}
	if _, ok := t.peers[peer.ID()]; ok {
		return fmt.Errorf("peer already registered: %s", peer.ID())
	}
	t.peers[peer.ID()] = peer
	return nil
}

func (t *Transport) Unregister(id string) error {
	if t == nil {
		return nil
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	peer, ok := t.peers[id]
	if !ok {
		return nil
	}
	delete(t.peers, id)
	if peer == nil {
		return nil
	}
	return peer.Close()
}

func (t *Transport) Send(ctx context.Context, id string, payload []byte) error {
	if t == nil {
		return fmt.Errorf("transport is nil")
	}
	if ctx == nil {
		return fmt.Errorf("context is nil")
	}
	if len(payload) > MaxTransportPayload {
		return fmt.Errorf("payload exceeds %d bytes", MaxTransportPayload)
	}
	t.mu.RLock()
	peer, ok := t.peers[id]
	t.mu.RUnlock()
	if !ok || peer == nil {
		return fmt.Errorf("peer not connected: %s", id)
	}
	return peer.Send(ctx, payload)
}
