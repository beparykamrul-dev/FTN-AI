package mesh

import (
 "sync"
 "time"
)

type LinkState string
const (
 LinkUp LinkState = "up"
 LinkDown LinkState = "down"
 LinkDegraded LinkState = "degraded"
)

type Link struct {
 ID string `json:"id"`
 From string `json:"from"`
 To string `json:"to"`
 State LinkState `json:"state"`
 LatencyMS float64 `json:"latency_ms"`
 LossPercent float64 `json:"loss_percent"`
 Metric uint32 `json:"metric"`
 UpdatedAt time.Time `json:"updated_at"`
}

type LinkStateStore struct { mu sync.RWMutex; links map[string]Link }
func NewLinkStateStore() *LinkStateStore { return &LinkStateStore{links: make(map[string]Link)} }
func (s *LinkStateStore) Upsert(l Link) { s.mu.Lock(); defer s.mu.Unlock(); if l.UpdatedAt.IsZero(){l.UpdatedAt=time.Now().UTC()}; s.links[l.ID]=l }
func (s *LinkStateStore) Get(id string) (Link,bool) { s.mu.RLock(); defer s.mu.RUnlock(); l,ok:=s.links[id]; return l,ok }
func (s *LinkStateStore) Snapshot() []Link { s.mu.RLock(); defer s.mu.RUnlock(); out:=make([]Link,0,len(s.links)); for _,l:=range s.links {out=append(out,l)}; return out }
