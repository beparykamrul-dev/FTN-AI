package data

import "context"

// CockroachBackend reserves the adapter boundary for CockroachDB. The actual
// SQL/driver implementation belongs behind this contract.
type CockroachBackend struct {
	PingFunc func(context.Context) error
}

func (c CockroachBackend) Name() string { return "cockroachdb" }

func (c CockroachBackend) Ping(ctx context.Context) error {
	if c.PingFunc == nil {
		return nil
	}
	return c.PingFunc(ctx)
}
