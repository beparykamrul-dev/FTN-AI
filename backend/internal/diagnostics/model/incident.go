package model

import "strings"

// Incident is the normalized operational incident identity.
type Incident struct {
	ID string `json:"id"`
	Severity string `json:"severity"`
	NodeID string `json:"node_id,omitempty"`
	MAC string `json:"mac,omitempty"`
	IP string `json:"ip,omitempty"`
	Interface string `json:"interface,omitempty"`
	Service string `json:"service,omitempty"`
	PathID string `json:"path_id,omitempty"`
	EvidenceIDs []string `json:"evidence_ids,omitempty"`
}

func (i Incident) Valid() bool { return i.Normalize().ID != "" && i.Normalize().Severity != "" }

func (i Incident) Normalize() Incident {
	i.ID = strings.TrimSpace(i.ID)
	i.Severity = strings.TrimSpace(i.Severity)
	i.NodeID = strings.TrimSpace(i.NodeID)
	i.MAC = strings.TrimSpace(i.MAC)
	i.IP = strings.TrimSpace(i.IP)
	i.Interface = strings.TrimSpace(i.Interface)
	i.Service = strings.TrimSpace(i.Service)
	i.PathID = strings.TrimSpace(i.PathID)
	seen := make(map[string]struct{}, len(i.EvidenceIDs))
	ids := make([]string, 0, len(i.EvidenceIDs))
	for _, id := range i.EvidenceIDs {
		id = strings.TrimSpace(id)
		if id == "" { continue }
		if _, ok := seen[id]; ok { continue }
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	i.EvidenceIDs = ids
	return i
}
