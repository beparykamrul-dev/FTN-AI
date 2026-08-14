package fiber

import "context"

type FiberMetric struct {
	Name string
	Value float64
	Unit string
	Labels map[string]string
	Timestamp string
}

// MetricsRepository is intentionally local/provider-neutral.
type MetricsRepository interface {
	Observe(context.Context, FiberMetric) error
	Latest(context.Context, string, string) (FiberMetric, error)
}

type FiberHealth struct {
	PathID string
	OpticalLossDb float64
	DistanceMeters float64
	ActiveUsers int
	ActiveONUs int
	OpenCuts int
	HealthScore float64
	Status string
}
