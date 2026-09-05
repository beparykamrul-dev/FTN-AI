package proxy
import("strings";"time")
type FailoverPolicy struct{MaxAttempts int;RetryDelay time.Duration;RequireHealthy bool;RequireSecure bool}
func DefaultFailoverPolicy()FailoverPolicy{return FailoverPolicy{MaxAttempts:2,RetryDelay:25*time.Millisecond,RequireHealthy:true,RequireSecure:true}}
func SelectFailover(policy FailoverPolicy,ranked []Upstream,failedID string,attempt int)(Upstream,bool){max:=policy.MaxAttempts;if max<=0{max=2};if attempt<0||attempt>=max{return Upstream{},false};failedID=strings.TrimSpace(failedID);seen:=make(map[string]struct{},len(ranked));for _,u:=range ranked{id:=strings.TrimSpace(u.ID);endpoint:=strings.TrimSpace(u.Endpoint);if id==""||id==failedID||endpoint==""{continue};if _,ok:=seen[id];ok{continue};seen[id]=struct{}{};if policy.RequireHealthy&&!u.Healthy{continue};if policy.RequireSecure&&!u.Secure{continue};u.ID=id;u.Endpoint=endpoint;return u,true};return Upstream{},false}
