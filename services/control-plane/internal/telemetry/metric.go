package telemetry

type Metric struct {
	DeviceID    string  `json:"device_id"`
	InterfaceID string  `json:"interface_id,omitempty"`
	Name        string  `json:"name"`
	Value       float64 `json:"value"`
	Unit        string  `json:"unit,omitempty"`
}
