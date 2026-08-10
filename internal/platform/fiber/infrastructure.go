package fiber

import "time"

type CustomerProfile struct {
	ID string `json:"id"`
	Name string `json:"name,omitempty"`
	ServiceID string `json:"service_id,omitempty"`
	Package string `json:"package,omitempty"`
	PPPoEUser string `json:"pppoe_user,omitempty"`
	ONU string `json:"onu,omitempty"`
	Router string `json:"router,omitempty"`
	IP string `json:"ip,omitempty"`
	Status string `json:"status,omitempty"`
}

type FiberIncident struct {
	ID string `json:"id"`
	LinkID string `json:"link_id"`
	Status string `json:"status"`
	DetectedAt time.Time `json:"detected_at"`
	Recovery RecoveryCandidate `json:"recovery"`
}

type GISRecord struct {
	Node FiberNode `json:"node"`
	Links []FiberLink `json:"links,omitempty"`
	Customers []CustomerProfile `json:"customers,omitempty"`
	Incidents []FiberIncident `json:"incidents,omitempty"`
}

// BuildGISRecord composes the API-facing graph record. Persistence belongs to
// the database adapter; this model intentionally contains no DB credentials.
func BuildGISRecord(node FiberNode, links []FiberLink, customers []CustomerProfile, incidents []FiberIncident) GISRecord {
	return GISRecord{Node: node, Links: links, Customers: customers, Incidents: incidents}
}
