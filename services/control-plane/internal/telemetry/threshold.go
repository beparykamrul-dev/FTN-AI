package telemetry

type Threshold struct {
	Metric string  `json:"metric"`
	Warning float64 `json:"warning"`
	Critical float64 `json:"critical"`
	Unit string `json:"unit,omitempty"`
}
