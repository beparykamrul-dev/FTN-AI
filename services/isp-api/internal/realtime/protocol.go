package realtime

import "time"

type Event struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Channel   string    `json:"channel"`
	Sequence  uint64    `json:"sequence"`
	Payload   any       `json:"payload"`
}

type Subscription struct {
	Channel string `json:"channel"`
	Scope   string `json:"scope"`
}

var AllowedChannels = map[string]struct{}{
	"account": {}, "billing": {}, "notifications": {}, "tickets": {},
	"incidents": {}, "topology": {}, "device_telemetry": {},
	"discovery": {}, "recovery": {},
}

func IsAllowedChannel(channel string) bool {
	_, ok := AllowedChannels[channel]
	return ok
}
