package ftndns

import (
	"context"
	"sync"
)

type QueryDeduper struct { mu sync.Mutex; calls map[string]*dedupCall }
type dedupCall struct { done chan struct{}; result ResolveResult; err error }

func NewQueryDeduper() *QueryDeduper { return &QueryDeduper{calls: make(map[string]*dedupCall)} }

func (d *QueryDeduper) Do(ctx context.Context, key string, fn ResolveFunc) (ResolveResult, error) {
	if ctx == nil { return ResolveResult{}, context.Canceled }
	d.mu.Lock()
	if call, ok := d.calls[key]; ok { d.mu.Unlock(); select { case <-call.done: return call.result, call.err; case <-ctx.Done(): return ResolveResult{}, ctx.Err() } }
	call := &dedupCall{done: make(chan struct{})}; d.calls[key]=call; d.mu.Unlock()
	call.result, call.err = fn(ctx, ResolveRequest{Key:key})
	close(call.done)
	d.mu.Lock(); delete(d.calls,key); d.mu.Unlock()
	return call.result, call.err
}
