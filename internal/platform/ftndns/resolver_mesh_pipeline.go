package ftndns

import (
    "context"
    "fmt"
    "time"
)

type MeshResolverConfig struct {
    Upstreams []Upstream
    CacheNodes []CacheNode
    LocalSnapshotHash string
    Prefetch *Prefetcher
    RefreshGuard *RefreshGuard
}

type MeshResolver struct {
    Resolver *FTNDNSResolver
    Selector UpstreamSelector
    CachePeers CacheMeshCoordinator
    Config MeshResolverConfig
}

func NewMeshResolver(resolver *FTNDNSResolver, config MeshResolverConfig) *MeshResolver {
    return &MeshResolver{Resolver:resolver, Config:config}
}

func (m *MeshResolver) Resolve(ctx context.Context, key string, resolveUpstream func(context.Context, ResolveRequest, Upstream) (ResolveResult,error)) (ResolveResult,error) {
    if m == nil || m.Resolver == nil { return ResolveResult{},fmt.Errorf("resolver is required") }
    selected, ok := m.Selector.Select(m.Config.Upstreams)
    if !ok { return ResolveResult{},fmt.Errorf("no healthy upstream available") }
    result, err := m.Resolver.Resolve(ctx,key,func(c context.Context, req ResolveRequest)(ResolveResult,error){
        return resolveUpstream(c,req,selected)
    })
    if err != nil { return ResolveResult{},err }
    if m.Config.Prefetch != nil && m.Config.RefreshGuard != nil && m.Config.Prefetch.ShouldRefresh(key,time.Now().Add(time.Second),time.Now()) {
        if m.Config.RefreshGuard.TryAcquire(key,time.Now(),time.Second) { m.Config.Prefetch.Trigger(ctx,key,func(c context.Context)([]string,time.Time,error){ r,e:=resolveUpstream(c,ResolveRequest{Key:key},selected); return r.Value,time.Now().Add(time.Minute),e }); go m.Config.RefreshGuard.Release(key) }
    }
    return result,nil
}
