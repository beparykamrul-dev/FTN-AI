package ftndns

import (
	"context"
	"errors"
	"strings"
	"sync"
)

type QueryDeduper struct { mu sync.Mutex; calls map[string]*dedupCall }
type dedupCall struct { done chan struct{}; result ResolveResult; err error }

func NewQueryDeduper() *QueryDeduper { return &QueryDeduper{calls: make(map[string]*dedupCall)} }

func (d *QueryDeduper) Do(ctx context.Context, key string, fn ResolveFunc) (ResolveResult, error) {
	if d == nil { return ResolveResult{}, errors.New("query deduper is required") }
	if ctx == nil { return ResolveResult{}, context.Canceled }
	if fn == nil { return ResolveResult{}, errors.New("resolver function is required") }
	key = strings.TrimSpace(key)
	if key == "" { return ResolveResult{}, errors.New("query key is required") }
	d.mu.Lock()
	if d.calls == nil { d.calls = make(map[string]*dedupCall) }
	if call, ok := d.calls[key]; ok { d.mu.Unlock(); select { case <-call.done: return call.result, call.err; case <-ctx.Done(): return ResolveResult{}, ctx.Err() } }
	call := &dedupCall{done: make(chan struct{})}; d.calls[key] = call; d.mu.Unlock()
	call.result, call.err = fn(ctx, ResolveRequest{Key: key})
	close(call.done)
	d.mu.Lock(); delete(d.calls, key); d.mu.Unlock()
	return call.result, call.err
}
