package dns

import("context";"fmt";"strings";"sync";"time")
type ProbePolicy struct{FailureThreshold uint32 `json:"failure_threshold"`;SuccessThreshold uint32 `json:"success_threshold"`;MaxLatency time.Duration `json:"max_latency"`}
type ProbeState struct{Failures uint32 `json:"failures"`;Successes uint32 `json:"successes"`;Healthy bool `json:"healthy"`}
type ProbePolicyController struct{mu sync.Mutex;states map[string]ProbeState}
func NewProbePolicyController()*ProbePolicyController{return &ProbePolicyController{states:make(map[string]ProbeState)}}
func(p ProbePolicy)Validate()error{if p.FailureThreshold==0||p.SuccessThreshold==0{return fmt.Errorf("probe thresholds must be positive")};if p.MaxLatency<=0{return fmt.Errorf("max latency must be positive")};return nil}
func(c *ProbePolicyController)Observe(ctx context.Context,key string,result DNSProbeResult,policy ProbePolicy)(bool,error){if c==nil{return false,fmt.Errorf("probe policy controller is required")};if err:=policy.Validate();err!=nil{return false,err};if ctx==nil{return false,context.Canceled};if err:=ctx.Err();err!=nil{return false,err};key=strings.TrimSpace(key);if key==""{return false,fmt.Errorf("probe policy key is required")};c.mu.Lock();defer c.mu.Unlock();state:=c.states[key];good:=result.Reachable&&result.Latency<=policy.MaxLatency;if good{state.Successes++;state.Failures=0;if state.Successes>=policy.SuccessThreshold{state.Healthy=true}}else{state.Failures++;state.Successes=0;if state.Failures>=policy.FailureThreshold{state.Healthy=false}};c.states[key]=state;return state.Healthy,nil}
