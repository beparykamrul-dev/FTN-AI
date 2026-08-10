package fiber

import "time"

type LinkMetrics struct {
	LinkID string `json:"link_id"`
	LatencyMS float64 `json:"latency_ms"`
	PacketLossPct float64 `json:"packet_loss_pct"`
	RxMbps float64 `json:"rx_mbps"`
	TxMbps float64 `json:"tx_mbps"`
	OpticalPowerDBM *float64 `json:"optical_power_dbm,omitempty"`
	AvailabilityPct float64 `json:"availability_pct"`
	ObservedAt time.Time `json:"observed_at"`
}

type CustomerMetrics struct {
	CustomerID string `json:"customer_id"`
	ServiceStatus string `json:"service_status"`
	LatencyMS float64 `json:"latency_ms"`
	PacketLossPct float64 `json:"packet_loss_pct"`
	UsageMbps float64 `json:"usage_mbps"`
	ObservedAt time.Time `json:"observed_at"`
}

type MetricAlert struct {
	TargetID string `json:"target_id"`
	Metric string `json:"metric"`
	Value float64 `json:"value"`
	Threshold float64 `json:"threshold"`
	Severity string `json:"severity"`
	ObservedAt time.Time `json:"observed_at"`
}

func EvaluateLinkMetrics(m LinkMetrics) []MetricAlert {
	alerts := make([]MetricAlert, 0)
	if m.PacketLossPct >= 5 { alerts = append(alerts, MetricAlert{TargetID:m.LinkID, Metric:"packet_loss_pct", Value:m.PacketLossPct, Threshold:5, Severity:"high", ObservedAt:m.ObservedAt}) }
	if m.LatencyMS >= 200 { alerts = append(alerts, MetricAlert{TargetID:m.LinkID, Metric:"latency_ms", Value:m.LatencyMS, Threshold:200, Severity:"medium", ObservedAt:m.ObservedAt}) }
	if m.AvailabilityPct > 0 && m.AvailabilityPct < 99 { alerts = append(alerts, MetricAlert{TargetID:m.LinkID, Metric:"availability_pct", Value:m.AvailabilityPct, Threshold:99, Severity:"medium", ObservedAt:m.ObservedAt}) }
	return alerts
}
