package observability

// TraceSpan links one FTN operation to its parent trace without storing payload contents.
type TraceSpan struct {
	TraceID string
	SpanID  string
	ParentID string
	Service string
	Operation string
	DurationMS uint64
	Status string
}

func (s TraceSpan) Valid() bool {
	return s.TraceID != "" && s.SpanID != "" && s.Service != "" && s.Operation != ""
}
