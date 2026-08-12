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
		select {
		case l.sem <- struct{}{}:
			go func() {
				defer func() { <-l.sem; _ = conn.Close() }()
				if l.cfg.ConnectionTimeout > 0 {
					_ = conn.SetDeadline(time.Now().Add(l.cfg.ConnectionTimeout))
					defer conn.SetDeadline(time.Time{})
				}
				_ = ServeSession(ctx, conn, l.cfg.Session, auth, delivery)
			}()
		default:
			_ = conn.Close()
		}
	}
}
