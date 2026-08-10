package monitoring

import "time"

type SNMPTarget struct {
	ID string `json:"id"`
	Address string `json:"address"`
	Port uint16 `json:"port"`
	Community string `json:"community,omitempty"`
	Version string `json:"version"`
	Timeout time.Duration `json:"timeout"`
	Retries int `json:"retries"`
}

type SNMPSample struct {
	TargetID string `json:"target_id"`
	OID string `json:"oid"`
	Value any `json:"value"`
	ObservedAt time.Time `json:"observed_at"`
}

type SNMPPoller interface {
	Poll(target SNMPTarget, oids []string) ([]SNMPSample, error)
}

func PollTarget(poller SNMPPoller, target SNMPTarget, oids []string) ([]SNMPSample, error) {
	if poller == nil { return nil, nil }
	return poller.Poll(target, oids)
}
