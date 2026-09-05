package transport

import (
	"net"
	"sync"
	"time"
)

type AuthThrottle struct { mu sync.Mutex; failures map[string]rateEntry; MaxFailures int; Window time.Duration; BlockFor time.Duration; admission *IPAdmission }
func NewAuthThrottle(admission *IPAdmission, maxFailures int, window, blockFor time.Duration) *AuthThrottle { if maxFailures<=0{maxFailures=5};if window<=0{window=5*time.Minute};if blockFor<=0{blockFor=15*time.Minute};return &AuthThrottle{failures:make(map[string]rateEntry),MaxFailures:maxFailures,Window:window,BlockFor:blockFor,admission:admission} }
func (t *AuthThrottle) Failed(ip net.IP, now time.Time) bool { if t==nil||ip==nil{return true};key:=ip.String();if key=="<nil>"||key==""{return true};if now.IsZero(){now=time.Now().UTC()}else{now=now.UTC()};t.mu.Lock();if t.failures==nil{t.failures=make(map[string]rateEntry)};e:=t.failures[key];if e.window.IsZero()||now.Sub(e.window)>=t.Window{e=rateEntry{window:now}};if now.Before(e.window){t.mu.Unlock();return true};e.count++;t.failures[key]=e;blocked:=e.count>=t.MaxFailures;t.mu.Unlock();if blocked&&t.admission!=nil{t.admission.Block(ip,now)};return blocked }
func (t *AuthThrottle) Success(ip net.IP) { if t==nil||ip==nil{return};key:=ip.String();t.mu.Lock();delete(t.failures,key);t.mu.Unlock() }
func (t *AuthThrottle) Cleanup(now time.Time) { if t==nil{return};if now.IsZero(){now=time.Now().UTC()}else{now=now.UTC()};t.mu.Lock();for ip,e:=range t.failures{if !now.Before(e.window)&&now.Sub(e.window)>=t.Window{delete(t.failures,ip)}};t.mu.Unlock() }
