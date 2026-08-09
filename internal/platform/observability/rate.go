package observability

import "time"

type TrafficRate struct {
	Interface string  `json:"interface"`
	BPS       float64 `json:"bps"`
	PPS       float64 `json:"pps"`
}

func Rate(prev, current TrafficSample, elapsed time.Duration) TrafficRate {
	if elapsed <= 0 { return TrafficRate{Interface: current.Interface} }
	seconds := elapsed.Seconds()
	var b, p uint64
	if current.Bytes >= prev.Bytes { b = current.Bytes - prev.Bytes }
	if current.Packets >= prev.Packets { p = current.Packets - prev.Packets }
	return TrafficRate{Interface: current.Interface, BPS: float64(b) / seconds, PPS: float64(p) / seconds}
}
