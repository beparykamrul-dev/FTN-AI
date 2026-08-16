package observability

// Metric is a normalized numeric time-series sample.
type Metric struct {
	Name      string
	Value     float64
	TimestampUnix int64
	Unit      string
}

func (m Metric) Valid() bool {
	return m.Name != "" && m.TimestampUnix > 0
}
