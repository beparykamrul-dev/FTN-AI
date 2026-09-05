package runtime

import("context";"sort";"strings";"time")
type HealthState string
const(HealthHealthy HealthState="healthy";HealthFailed HealthState="failed";HealthUnknown HealthState="unknown")
type HealthResult struct{Service string;State HealthState;CheckedAt time.Time;Error string}
func CheckHealth(ctx context.Context,manager *ServiceManager)[]HealthResult{if ctx==nil||manager==nil{return []HealthResult{}};if err:=ctx.Err();err!=nil{return []HealthResult{{State:HealthUnknown,CheckedAt:time.Now().UTC(),Error:err.Error()}}};health:=manager.Health(ctx);results:=make([]HealthResult,0,len(health));now:=time.Now().UTC();for name,err:=range health{name=strings.TrimSpace(name);if name==""{continue};result:=HealthResult{Service:name,State:HealthHealthy,CheckedAt:now};if err!=nil{result.State=HealthFailed;result.Error=strings.TrimSpace(err.Error())};results=append(results,result)};sort.Slice(results,func(i,j int)bool{return results[i].Service<results[j].Service});return results}
