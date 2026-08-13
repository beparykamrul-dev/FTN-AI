package model

// Incident is the normalized operational incident identity.
type Incident struct {
	ID          string   `json:"id"`
	Severity    string   `json:"severity"`
	NodeID      string   `json:"node_id,omitempty"`
	MAC         string   `json:"mac,omitempty"`
	IP          string   `json:"ip,omitempty"`
	Interface   string   `json:"interface,omitempty"`
	Service     string   `json:"service,omitempty"`
	PathID      string   `json:"path_id,omitempty"`
	EvidenceIDs []string `json:"evidence_ids,omitempty"`
}
