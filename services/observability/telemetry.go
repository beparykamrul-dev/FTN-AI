package observability

// TelemetryEvent is a privacy-conscious normalized event for FTN observability.
type TelemetryEvent struct {
	Service   string
	Operation string
	Status    string
	LatencyMS uint64
	TraceID   string
	ErrorCode string
}

func (e TelemetryEvent) Valid() bool {
	return e.Service != "" && e.Operation != "" && e.Status != ""
}
