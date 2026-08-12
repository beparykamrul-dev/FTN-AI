package transport

import (
	"context"
	"errors"
	"net"
	"time"
)

// AdmissionMiddleware is the connection-level integration point for the
// lightweight IP admission and authentication throttling controls.
type AdmissionMiddleware struct {
	Guard *AdmissionGuard
}

func (m *AdmissionMiddleware) Accept(ctx context.Context, conn net.Conn) error {
	if m == nil || m.Guard == nil || conn == nil {
		return errors.New("mail admission unavailable")
	}
	remote, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		return errors.New("invalid remote address")
	}
	ip := net.ParseIP(remote)
	if ip == nil || !m.Guard.Allow(ip, time.Now()) {
		return errors.New("connection refused by mail admission policy")
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
