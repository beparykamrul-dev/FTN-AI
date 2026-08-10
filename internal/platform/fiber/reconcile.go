package fiber

import "sort"

type ChangeKind string

const (
	Added ChangeKind = "added"
	Removed ChangeKind = "removed"
	Changed ChangeKind = "changed"
)

type TopologyChange struct {
	Kind ChangeKind `json:"kind"`
	Entity DiscoveredEntity `json:"entity"`
	Previous *DiscoveredEntity `json:"previous,omitempty"`
}

func Reconcile(previous, current TopologySnapshot) []TopologyChange {
	old := make(map[string]DiscoveredEntity, len(previous.Entities))
	for _, e := range previous.Entities { old[e.ExternalID] = e }
	changes := make([]TopologyChange, 0)
	seen := make(map[string]bool, len(current.Entities))
	for _, e := range current.Entities {
		seen[e.ExternalID] = true
		p, ok := old[e.ExternalID]
		if !ok { changes = append(changes, TopologyChange{Kind: Added, Entity:e}); continue }
		if p.ParentID != e.ParentID || p.Kind != e.Kind || p.Status != e.Status || p.Name != e.Name {
			prev := p
			changes = append(changes, TopologyChange{Kind:Changed, Entity:e, Previous:&prev})
		}
	}
	for _, e := range previous.Entities {
		if !seen[e.ExternalID] { changes = append(changes, TopologyChange{Kind:Removed, Entity:e}) }
	}
	sort.Slice(changes, func(i,j int) bool { return changes[i].Entity.ExternalID < changes[j].Entity.ExternalID })
	return changes
}
