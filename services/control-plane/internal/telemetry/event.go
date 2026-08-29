package telemetry

import "time"

type Event struct {
	ID        string    `json:"id"`
	DeviceID  string    `json:"device_id"`
	Source    string    `json:"source"`
	Metric    string    `json:"metric"`
	Value     float64   `json:"value"`
	Unit      string    `json:"unit,omitempty"`
	ObservedAt time.Time `json:"observed_at"`
}
