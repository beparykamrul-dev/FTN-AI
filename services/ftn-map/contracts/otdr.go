package fiber

import "context"

type OTDRTrace struct {
	ID string
	FiberPathID string
	WavelengthNm int
	PulseWidthNs int
	RangeMeters float64
	TotalLossDb float64
	ORLDb float64
	EventCount int
	MeasuredAt string
	Status string
}

type OTDREvent struct {
	DistanceMeters float64
	LossDb float64
	ReflectanceDb float64
	EventType string
	Severity string
}

// OTDRRepository keeps traces/events local to FTN storage.
type OTDRRepository interface {
	SaveTrace(context.Context, OTDRTrace, []OTDREvent) error
	LatestTrace(context.Context, string) (OTDRTrace, []OTDREvent, error)
}

// OTDRAnalyzer converts traces into topology/fault evidence. It is advisory;
// network recovery remains policy-gated.
type OTDRAnalyzer interface {
	Analyze(context.Context, OTDRTrace, []OTDREvent) (string, error)
	CorrelateCut(context.Context, FiberPath, OTDRTrace, []OTDREvent) (*CutEvent, error)
}
