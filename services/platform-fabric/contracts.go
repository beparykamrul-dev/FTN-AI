package platformfabric

import "time"

type Identity struct {
	Subject string `json:"subject"`
	Kind string `json:"kind"`
	Issuer string `json:"issuer"`
}

type AuditEvent struct {
	ID string `json:"id"`
	Subject string `json:"subject"`
	Action string `json:"action"`
	Resource string `json:"resource"`
	CorrelationID string `json:"correlation_id"`
	Timestamp time.Time `json:"timestamp"`
}

type ServiceEvent struct {
	Type string `json:"type"`
	ServiceID string `json:"service_id"`
	Subject string `json:"subject,omitempty"`
	Status string `json:"status,omitempty"`
	CorrelationID string `json:"correlation_id"`
	Timestamp time.Time `json:"timestamp"`
}

type MetricSample struct {
	Name string `json:"name"`
	Value float64 `json:"value"`
	Unit string `json:"unit,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}
