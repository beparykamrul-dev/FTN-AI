package ftndns

import (
    "context"
    "fmt"
    "time"
)

type FTNDNSResolver struct {
    Cache *ResolverCache
    Negative *NegativeCache
    Admission *CacheAdmission
    Deduper *QueryDeduper
    Metrics *ResolverMetrics
    Workers *ResolverWorkerPool
}

func NewFTNDNSResolver(cache *ResolverCache, negative *NegativeCache, admission *CacheAdmission, workers *ResolverWorkerPool) *FTNDNSResolver {
    return &FTNDNSResolver{Cache:cache, Negative:negative, Admission:admission, Deduper:NewQueryDeduper(), Metrics:&ResolverMetrics{}, Workers:workers}
}

func (r *FTNDNSResolver) Resolve(ctx context.Context, key string, upstream ResolveFunc) (ResolveResult,error) {
    if ctx==nil { return ResolveResult{},fmt.Errorf("context is required") }
    if key=="" || upstream==nil { return ResolveResult{},fmt.Errorf("key and upstream resolver are required") }
    now:=time.Now()
    if r.Cache!=nil { if value,ok:=r.Cache.Get(key,now); ok { r.Metrics.RecordQuery(true); return ResolveResult{Value:value,CacheHit:true},nil } }
    if r.Negative!=nil && r.Negative.Hit(key,now) { r.Metrics.NegativeHits.Add(1); r.Metrics.RecordQuery(false); return ResolveResult{},nil }
    if r.Admission!=nil { r.Admission.Observe(key) }
    r.Metrics.RecordQuery(false)
    resolve:=func(c context.Context, req ResolveRequest)(ResolveResult,error){
        if r.Workers!=nil { return r.Workers.Submit(c,req,upstream) }
        return upstream(c,req)
    }
    return r.Deduper.Do(ctx,key,resolve)
}
