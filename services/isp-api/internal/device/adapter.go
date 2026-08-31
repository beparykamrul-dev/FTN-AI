package device

import "context"

type Kind string

const (
	KindMikroTik Kind = "mikrotik"
	KindOLT      Kind = "olt"
	KindONU      Kind = "onu"
	KindSwitch   Kind = "switch"
	KindServer   Kind = "server"
	KindFTNNode  Kind = "ftn-node"
)

type Device struct {
	ID       string `json:"id"`
	Kind     Kind   `json:"kind"`
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
}

type Telemetry struct {
	DeviceID string         `json:"device_id"`
	Status   string         `json:"status"`
	Metrics  map[string]any `json:"metrics,omitempty"`
}

type Adapter interface {
	Kind() Kind
	GetStatus(context.Context, Device) (Telemetry, error)
	GetInventory(context.Context, Device) (map[string]any, error)
	GetTelemetry(context.Context, Device) (Telemetry, error)
	Close(context.Context) error
}
