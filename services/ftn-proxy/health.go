package proxy

import "time"

type HealthState uint8
const(HealthUnknown HealthState=iota;HealthHealthy;HealthUnhealthy)
type HealthCheck struct{Interval time.Duration;Timeout time.Duration;FailureLimit int;SuccessLimit int}
func DefaultHealthCheck()HealthCheck{return HealthCheck{Interval:10*time.Second,Timeout:3*time.Second,FailureLimit:3,SuccessLimit:2}}
type HealthTracker struct{State HealthState;Failures int;Successes int}
func(h *HealthTracker)Observe(success bool,policy HealthCheck)HealthState{if h==nil{return HealthUnknown};if success{h.Failures=0;if h.Successes< int(^uint(0)>>1){h.Successes++};limit:=policy.SuccessLimit;if limit<=0{limit=2};if h.Successes>=limit{h.State=HealthHealthy};return h.State};h.Successes=0;if h.Failures< int(^uint(0)>>1){h.Failures++};limit:=policy.FailureLimit;if limit<=0{limit=3};if h.Failures>=limit{h.State=HealthUnhealthy};return h.State}
