package transport

import (
	"net"
	"sync"
	"time"
)

type rateEntry struct { window time.Time; count int }
type IPRateLimiter struct { mu sync.Mutex; entries map[string]rateEntry; MaxAttempts int; Window time.Duration }
func NewIPRateLimiter(maxAttempts int, window time.Duration) *IPRateLimiter { if maxAttempts <= 0 { maxAttempts = 20 }; if window <= 0 { window = time.Minute }; return &IPRateLimiter{entries:make(map[string]rateEntry),MaxAttempts:maxAttempts,Window:window} }
func (r *IPRateLimiter) Allow(ip net.IP, now time.Time) bool { if r == nil || ip == nil { return false }; key:=ip.String(); if key=="<nil>"||key=="" { return false }; if now.IsZero(){now=time.Now().UTC()}else{now=now.UTC()}; r.mu.Lock(); defer r.mu.Unlock(); if r.entries==nil{r.entries=make(map[string]rateEntry)}; e:=r.entries[key]; if e.window.IsZero()||now.Sub(e.window)>=r.Window { r.entries[key]=rateEntry{window:now,count:1}; return true }; if now.Before(e.window){return false}; if e.count>=r.MaxAttempts{return false}; e.count++; r.entries[key]=e; return true }
func (r *IPRateLimiter) Cleanup(now time.Time) { if r==nil{return}; if now.IsZero(){now=time.Now().UTC()}else{now=now.UTC()}; r.mu.Lock(); defer r.mu.Unlock(); for ip,e:=range r.entries { if !now.Before(e.window)&&now.Sub(e.window)>=r.Window {delete(r.entries,ip)} } }
