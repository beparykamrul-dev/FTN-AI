package proxy

import "time"

// Policy is a conservative baseline for the FTN proxy data plane.
type Policy struct {
	ConnectTimeout time.Duration
	IdleTimeout    time.Duration
	MaxBodyBytes   int64
	RequireTLS     bool
	AllowHTTP2     bool
	StripHopByHop  bool
}

func DefaultPolicy() Policy {
	return Policy{
		ConnectTimeout: 5 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxBodyBytes:   32 << 20,
		RequireTLS:     true,
		AllowHTTP2:     true,
		StripHopByHop:  true,
	}
}

// Validate applies fail-closed checks before a request is forwarded.
func (p Policy) Validate(scheme string, bodyBytes int64) bool {
	if p.RequireTLS && scheme != "https" {
		return false
	}
	return p.MaxBodyBytes <= 0 || bodyBytes <= p.MaxBodyBytes
}
