package controlplane

import("strings";"sync";"time")

type RateLimiter struct{mu sync.Mutex;limit int;window time.Duration;seen map[string][]time.Time}
func NewRateLimiter(limit int,window time.Duration)*RateLimiter{if limit<1||window<=0{return nil};return &RateLimiter{limit:limit,window:window,seen:make(map[string][]time.Time)}}
func(r *RateLimiter)Allow(key string,now time.Time)bool{if r==nil{return false};key=strings.TrimSpace(key);if key==""{return false};if now.IsZero(){now=time.Now().UTC()}else{now=now.UTC()};r.mu.Lock();defer r.mu.Unlock();cut:=now.Add(-r.window);old:=r.seen[key];i:=0;for i<len(old)&&old[i].After(cut){i++};old=old[i:];if len(old)>=r.limit{r.seen[key]=old;return false};r.seen[key]=append(old,now);return true}
