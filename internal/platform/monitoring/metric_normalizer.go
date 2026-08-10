package monitoring

import "time"

type MetricKind string

const (
	MetricCounter MetricKind = "counter"
	MetricGauge MetricKind = "gauge"
)

type NormalizedMetric struct {
	DeviceID string `json:"device_id"`
	Metric string `json:"metric"`
	Kind MetricKind `json:"kind"`
	Value float64 `json:"value"`
	Unit string `json:"unit,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
	ObservedAt time.Time `json:"observed_at"`
}

func NormalizeSNMPSample(sample SNMPSample, deviceID, metric string, kind MetricKind, value float64, unit string, labels map[string]string) NormalizedMetric {
	observed := sample.ObservedAt
	if observed.IsZero() { observed = time.Now().UTC() }
	return NormalizedMetric{DeviceID:deviceID, Metric:metric, Kind:kind, Value:value, Unit:unit, Labels:labels, ObservedAt:observed}
}
