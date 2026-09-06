package controlplane

import "time"

type Readiness struct{Started bool;Dependencies int;HealthyDependencies int;LastCheck time.Time}
func(r Readiness)Ready()bool{return r.Started&&r.Dependencies>0&&r.HealthyDependencies==r.Dependencies&&!r.LastCheck.IsZero()}
