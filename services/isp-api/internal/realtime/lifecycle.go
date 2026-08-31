package realtime

import (
	"context"
	"time"
)

type Transport interface {
	Read(context.Context) (Event, error)
	Write(context.Context, Event) error
	Ping(context.Context) error
	Close() error
}

type Lifecycle struct {
	Heartbeat time.Duration
}

func (l Lifecycle) Run(ctx context.Context, transport Transport, onEvent func(Event) error) error {
	if transport == nil || onEvent == nil {
		return ErrUnauthorized
	}
	heartbeat := l.Heartbeat
	if heartbeat <= 0 {
		heartbeat = 30 * time.Second
	}
	ticker := time.NewTicker(heartbeat)
	defer ticker.Stop()
	defer transport.Close()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := transport.Ping(ctx); err != nil {
				return err
			}
		default:
			event, err := transport.Read(ctx)
			if err != nil {
				return err
			}
			if err := onEvent(event); err != nil {
				return err
			}
		}
	}
}
