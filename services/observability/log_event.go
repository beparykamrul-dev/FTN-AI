package observability

// LogEvent is a normalized structured log entry for the FTN observability pipeline.
type LogEvent struct {
	TimestampUnix int64
	Service string
	Level string
	Message string
	TraceID string
}

func (e LogEvent) Valid() bool {
	return e.TimestampUnix > 0 && e.Service != "" && e.Level != "" && e.Message != ""
}
