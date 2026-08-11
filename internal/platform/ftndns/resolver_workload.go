package ftndns

import (
	"context"
	"fmt"
	"sync"
)

type ResolveRequest struct { Key string }
type ResolveResult struct { Value []string; CacheHit bool }
type ResolveFunc func(context.Context, ResolveRequest) (ResolveResult, error)

type ResolverWorkerPool struct {
	workers int
	jobs chan resolveJob
	wg sync.WaitGroup
	once sync.Once
}

type resolveJob struct { ctx context.Context; req ResolveRequest; fn ResolveFunc; result chan resolveResponse }
type resolveResponse struct { result ResolveResult; err error }

func NewResolverWorkerPool(workers, queueSize int) *ResolverWorkerPool {
	if workers < 1 { workers = 1 }; if queueSize < workers { queueSize = workers }
	return &ResolverWorkerPool{workers: workers, jobs: make(chan resolveJob, queueSize)}
}

func (p *ResolverWorkerPool) Start() {
	p.once.Do(func() { for i:=0; i<p.workers; i++ { p.wg.Add(1); go p.worker() } })
}

func (p *ResolverWorkerPool) worker() { defer p.wg.Done(); for job := range p.jobs { r,e := job.fn(job.ctx,job.req); job.result <- resolveResponse{r,e} } }

func (p *ResolverWorkerPool) Submit(ctx context.Context, req ResolveRequest, fn ResolveFunc) (ResolveResult,error) {
	if ctx == nil { return ResolveResult{},fmt.Errorf("context is required") }; if fn == nil { return ResolveResult{},fmt.Errorf("resolver function is required") }
	p.Start(); result:=make(chan resolveResponse,1); job:=resolveJob{ctx:ctx,req:req,fn:fn,result:result}
	select { case p.jobs <- job: case <-ctx.Done(): return ResolveResult{},ctx.Err() }
	select { case response:=<-result: return response.result,response.err; case <-ctx.Done(): return ResolveResult{},ctx.Err() }
}

func (p *ResolverWorkerPool) Stop() { close(p.jobs); p.wg.Wait() }
