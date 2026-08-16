package observability

// DataTemperature classifies telemetry by access/value characteristics.
type DataTemperature string

const (
	TemperatureHot DataTemperature = "hot"
	TemperatureWarm DataTemperature = "warm"
	TemperatureCold DataTemperature = "cold"
	TemperatureArchive DataTemperature = "archive"
)

// TemperaturePolicy provides bounded thresholds for tier decisions.
type TemperaturePolicy struct {
	HotAgeSeconds int64
	WarmAgeSeconds int64
	ColdAgeSeconds int64
}

func (p TemperaturePolicy) Valid() bool {
	return p.HotAgeSeconds >= 0 && p.WarmAgeSeconds >= p.HotAgeSeconds && p.ColdAgeSeconds >= p.WarmAgeSeconds
}
