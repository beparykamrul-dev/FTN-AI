package transport

import (
	"context"
	"errors"
	"net"
	"time"
)

type ListenerConfig struct {
	Address           string
	MaxConnections    int
	ConnectionTimeout time.Duration
	Session           SessionConfig
	Admission         *AdmissionMiddleware
}

type Listener struct {
	cfg ListenerConfig
	sem chan struct{}
}

func NewListener(cfg ListenerConfig) (*Listener, error) {
	if cfg.Address == "" || cfg.MaxConnections <= 0 || cfg.Session.MaxMessageSize <= 0 {
		return nil, errors.New("invalid SMTP listener configuration")
	}
	return &Listener{cfg: cfg, sem: make(chan struct{}, cfg.MaxConnections)}, nil
}

func (l *Listener) Serve(ctx context.Context, ln net.Listener, auth Authenticator, delivery Delivery) error {
	if l == nil || ln == nil { return errors.New("invalid listener") }
	go func() { <-ctx.Done(); _ = ln.Close() }()
	for {
		conn, err := ln.Accept()
		if err != nil {
			select { case <-ctx.Done(): return nil; default: return err }
		}
		if l.cfg.Admission != nil && l.cfg.Admission.Accept(ctx, conn) != nil {
			_ = conn.Close()
			continue
		}
		select {
		case l.sem <- struct{}{}:
			go func(c net.Conn) {
				defer func() { <-l.sem; _ = c.Close() }()
				if l.cfg.ConnectionTimeout > 0 {
					_ = c.SetDeadline(time.Now().Add(l.cfg.ConnectionTimeout))
					defer c.SetDeadline(time.Time{})
				}
				_ = ServeSession(ctx, c, l.cfg.Session, auth, delivery)
			}(conn)
		default:
			_ = conn.Close()
		}
	}
}
