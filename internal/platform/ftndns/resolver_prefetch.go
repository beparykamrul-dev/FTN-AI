package ftndns

import (
	"context"
	"strings"
	"sync"
	"time"
)

type PrefetchPolicy struct { RefreshBefore time.Duration }
type Prefetcher struct { cache *ResolverCache; policy PrefetchPolicy; mu sync.Mutex; running map[string]bool }
func NewPrefetcher(cache *ResolverCache, policy PrefetchPolicy) *Prefetcher { if policy.RefreshBefore<=0{policy.RefreshBefore=30*time.Second};return &Prefetcher{cache:cache,policy:policy,running:make(map[string]bool)} }
func (p *Prefetcher) ShouldRefresh(key string,expiresAt,now time.Time) bool { if p==nil||strings.TrimSpace(key)==""||expiresAt.IsZero(){return false};return !expiresAt.After(now.Add(p.policy.RefreshBefore)) }
func (p *Prefetcher) Trigger(ctx context.Context,key string,refresh func(context.Context)([]string,time.Time,error)){if p==nil||p.cache==nil||refresh==nil||ctx==nil{return};key=strings.TrimSpace(key);if key==""{return};select{case<-ctx.Done():return;default:};p.mu.Lock();if p.running==nil{p.running=make(map[string]bool)};if p.running[key]{p.mu.Unlock();return};p.running[key]=true;p.mu.Unlock();go func(){defer func(){p.mu.Lock();delete(p.running,key);p.mu.Unlock()}();values,expires,err:=refresh(ctx);if err==nil{p.cache.Set(key,values,expires)}}()}
