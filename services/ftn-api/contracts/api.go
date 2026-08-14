package api

import "context"

type Service interface {
	Name() string
	Version() string
	Health(context.Context) error
}

type Metric struct {
	Name string
	Value float64
	Unit string
	Labels map[string]string
}

// MetricsSink is provider-neutral; production deployments can persist locally
// and expose metrics through the FTN API without exporting data to third parties.
type MetricsSink interface {
	Observe(context.Context, Metric) error
}
