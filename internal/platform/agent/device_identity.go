package agent

import "time"

type DeviceKind string

const (
	DeviceServer DeviceKind = "server"
	DeviceRouter DeviceKind = "router"
	DeviceOLT DeviceKind = "olt"
	DeviceONU DeviceKind = "onu"
	DevicePC DeviceKind = "pc"
	DeviceAndroid DeviceKind = "android"
	DeviceTV DeviceKind = "tv"
	DeviceVirtual DeviceKind = "virtual"
)

// DeviceIdentity is the canonical FTN identity shared by all managed nodes.
// Hardware identifiers are metadata only; authentication is handled by the
// FTN agent gateway and credentials are never stored here.
type DeviceIdentity struct {
	ID string `json:"id"`
	Kind DeviceKind `json:"kind"`
	Name string `json:"name"`
	Hostname string `json:"hostname,omitempty"`
	IP string `json:"ip,omitempty"`
	MAC string `json:"mac,omitempty"`
	Serial string `json:"serial,omitempty"`
	OS string `json:"os,omitempty"`
	AgentID string `json:"agent_id,omitempty"`
	Online bool `json:"online"`
	LastSeen time.Time `json:"last_seen,omitempty"`
}

func (d DeviceIdentity) HasStableIdentity() bool { return d.ID != "" || d.AgentID != "" || d.Serial != "" || d.MAC != "" }
