package transport

import (
	"net"
	"sync"
	"time"
)

type IPAdmission struct { Limiter *IPRateLimiter; mu sync.Mutex; blocked map[string]time.Time; BlockFor time.Duration }
func NewIPAdmission(limiter *IPRateLimiter, blockFor time.Duration) *IPAdmission { if blockFor<=0{blockFor=10*time.Minute};return &IPAdmission{Limiter:limiter,blocked:make(map[string]time.Time),BlockFor:blockFor} }
func (a *IPAdmission) Allow(ip net.IP, now time.Time) bool { if a==nil||a.Limiter==nil||ip==nil{return false};key:=ip.String();if key=="<nil>"||key==""{return false};if now.IsZero(){now=time.Now().UTC()}else{now=now.UTC()};a.mu.Lock();until,exists:=a.blocked[key];if exists&&now.Before(until){a.mu.Unlock();return false};if exists{delete(a.blocked,key)};a.mu.Unlock();return a.Limiter.Allow(ip,now) }
func (a *IPAdmission) Block(ip net.IP, now time.Time) { if a==nil||ip==nil{return};key:=ip.String();if key=="<nil>"||key==""{return};if now.IsZero(){now=time.Now().UTC()}else{now=now.UTC()};a.mu.Lock();if a.blocked==nil{a.blocked=make(map[string]time.Time)};a.blocked[key]=now.Add(a.BlockFor);a.mu.Unlock() }
func (a *IPAdmission) Cleanup(now time.Time) { if a==nil{return};if now.IsZero(){now=time.Now().UTC()}else{now=now.UTC()};a.mu.Lock();for ip,until:=range a.blocked{if !now.Before(until){delete(a.blocked,ip)}};a.mu.Unlock();if a.Limiter!=nil{a.Limiter.Cleanup(now)} }
