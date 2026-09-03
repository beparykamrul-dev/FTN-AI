package mesh

import (
	"errors"
	"sync"
)

type Topology struct {
	mu sync.RWMutex
	nodes map[string]Node
	links []Link
}

func NewTopology() *Topology { return &Topology{nodes: make(map[string]Node)} }

func (t *Topology) UpsertNode(n Node) error {
	if n.ID == "" { return errors.New("mesh node id is required") }
	t.mu.Lock(); defer t.mu.Unlock()
	t.nodes[n.ID] = n
	return nil
}

func (t *Topology) SetLink(l Link) error {
	if l.From == "" || l.To == "" || l.From == l.To { return errors.New("invalid mesh link") }
	t.mu.Lock(); defer t.mu.Unlock()
	for i, x := range t.links {
		if (x.From == l.From && x.To == l.To) || (x.From == l.To && x.To == l.From) {
			t.links[i] = l
			return nil
		}
	}
	t.links = append(t.links, l)
	return nil
}

func (t *Topology) Snapshot() ([]Node, []Link) {
	t.mu.RLock(); defer t.mu.RUnlock()
	nodes := make([]Node, 0, len(t.nodes))
	for _, n := range t.nodes { nodes = append(nodes, n) }
	links := append([]Link(nil), t.links...)
	return nodes, links
}
