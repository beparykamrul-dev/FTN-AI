package observability

// QoSClass defines relative priority for FTN traffic.
type QoSClass uint8

const (
	QoSBulk QoSClass = iota
	QoSStorage
	QoSControl
	QoSLive
)

// QoSReservation protects a minimum bandwidth budget for higher-priority traffic.
type QoSReservation struct {
	Class QoSClass
	ReservedMbps float64
	MaxMbps float64
}

func (r QoSReservation) Valid() bool {
	return r.MaxMbps > 0 && r.ReservedMbps >= 0 && r.ReservedMbps <= r.MaxMbps
}
