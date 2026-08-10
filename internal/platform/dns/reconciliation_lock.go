package dns

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type ReconciliationLock interface {
	Acquire(ctx context.Context, key string, lease time.Duration) (bool, error)
	Release(ctx context.Context, key string) error
}

type memoryLock struct {
	mu sync.Mutex
	leases map[string]time.Time
}

func NewMemoryReconciliationLock() ReconciliationLock {
	return &memoryLock{leases: make(map[string]time.Time)}
}

func (l *memoryLock) Acquire(ctx context.Context, key string, lease time.Duration) (bool, error) {
	if key == "" { return false, fmt.Errorf("lock key is required") }
	if lease <= 0 { return false, fmt.Errorf("lock lease must be positive") }
	select { case <-ctx.Done(): return false, ctx.Err(); default: }
	l.mu.Lock(); defer l.mu.Unlock()
	now := time.Now().UTC()
	if expires, ok := l.leases[key]; ok && expires.After(now) { return false, nil }
	l.leases[key] = now.Add(lease)
	return true, nil
}

func (l *memoryLock) Release(ctx context.Context, key string) error {
	if key == "" { return fmt.Errorf("lock key is required") }
	select { case <-ctx.Done(): return ctx.Err(); default: }
	l.mu.Lock(); defer l.mu.Unlock()
	delete(l.leases, key)
	return nil
}
