package ws

import "time"

// TransportPolicy keeps network-transport behavior inside FTN. It does not
// require a hosted realtime provider, external session service, or external
// message broker.
type TransportPolicy struct {
	RequireTLS             bool
	HandshakeTimeout       time.Duration
	ReadTimeout            time.Duration
	WriteTimeout           time.Duration
	PingInterval           time.Duration
	PongTimeout            time.Duration
	MaxConnectionsPerNode  int
	MaxMessageBytes        int64
	MaxSubscriptions       int
	ReconnectBaseDelay     time.Duration
	ReconnectMaxDelay      time.Duration
}

func DefaultTransportPolicy() TransportPolicy {
	return TransportPolicy{
		RequireTLS:            true,
		HandshakeTimeout:      10 * time.Second,
		ReadTimeout:           60 * time.Second,
		WriteTimeout:          10 * time.Second,
		PingInterval:          20 * time.Second,
		PongTimeout:           10 * time.Second,
		MaxConnectionsPerNode: 10000,
		MaxMessageBytes:       1 << 20,
		MaxSubscriptions:      128,
		ReconnectBaseDelay:    time.Second,
		ReconnectMaxDelay:     30 * time.Second,
	}
}

func (p TransportPolicy) Valid() bool {
	return p.HandshakeTimeout > 0 &&
		p.ReadTimeout > 0 &&
		p.WriteTimeout > 0 &&
		p.PingInterval > 0 &&
		p.PongTimeout > 0 &&
		p.MaxConnectionsPerNode > 0 &&
		p.MaxMessageBytes > 0 &&
		p.MaxSubscriptions > 0 &&
		p.ReconnectBaseDelay > 0 &&
		p.ReconnectMaxDelay >= p.ReconnectBaseDelay
}
