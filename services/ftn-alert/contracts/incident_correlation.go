package alert

import "context"

type SignalKind string

const (
	SignalOTDR SignalKind = "otdr"
	SignalOptical SignalKind = "optical"
	SignalDevice SignalKind = "device"
	SignalTopology SignalKind = "topology"
	SignalPPPoE SignalKind = "pppoe"
	SignalCustomerImpact SignalKind = "customer-impact"
	SignalLog SignalKind = "log"
)

type IncidentSignal struct {
	Kind SignalKind
	SourceID string
	ObservedAt string
	Severity string
	Confidence float64
	Fingerprint string
}

type CorrelatedIncident struct {
	ID string
	RootCause string
	Signals []IncidentSignal
	AffectedDevices []string
	AffectedPaths []string
	AffectedCustomers []string
	Confidence float64
	Status string
}

// IncidentCorrelator merges related telemetry into one auditable incident.
type IncidentCorrelator interface {
	Correlate(context.Context, []IncidentSignal) (CorrelatedIncident, error)
	Get(context.Context, string) (CorrelatedIncident, error)
	Resolve(context.Context, string, string) error
}
