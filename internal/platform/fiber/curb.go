package fiber

import "time"

// CURB (Customer/Utility/Route Boundary) records a field/network boundary
// relevant to GIS and recovery planning without directly changing the network.
type CURBRecord struct {
	ID string `json:"id"`
	Type string `json:"type"`
	NodeID string `json:"node_id,omitempty"`
	LinkID string `json:"link_id,omitempty"`
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Status string `json:"status"`
	ObservedAt time.Time `json:"observed_at"`
}

func BuildCURBRecord(id, typ, nodeID, linkID string, lat, lon float64, status string, observedAt time.Time) CURBRecord {
	if observedAt.IsZero() { observedAt = time.Now().UTC() }
	return CURBRecord{ID:id, Type:typ, NodeID:nodeID, LinkID:linkID, Latitude:lat, Longitude:lon, Status:status, ObservedAt:observedAt}
}
