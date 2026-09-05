package proxy

import (
	"strings"
	"time"
)

type Policy struct {
	ConnectTimeout time.Duration
	IdleTimeout    time.Duration
	MaxBodyBytes   int64
	RequireTLS     bool
	AllowHTTP2     bool
	StripHopByHop  bool
}

func DefaultPolicy() Policy {
	return Policy{ConnectTimeout: 5 * time.Second, IdleTimeout: 60 * time.Second, MaxBodyBytes: 32 << 20, RequireTLS: true, AllowHTTP2: true, StripHopByHop: true}
}

func (p Policy) Valid() bool {
	return p.ConnectTimeout > 0 && p.IdleTimeout >= p.ConnectTimeout && p.MaxBodyBytes >= 0
}

func (p Policy) Validate(scheme string, bodyBytes int64) bool {
	scheme = strings.ToLower(strings.TrimSpace(scheme))
	if !p.Valid() || bodyBytes < 0 {
		return false
	}
	if p.RequireTLS && scheme != "https" {
		return false
	}
	return p.MaxBodyBytes <= 0 || bodyBytes <= p.MaxBodyBytes
}
