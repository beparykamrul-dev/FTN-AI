package controlplane

import "math"

type ResourceHealth struct{CPUPercent float64;MemoryPercent float64;DiskPercent float64}
func HealthyResources(r ResourceHealth)bool{for _,v:=range []float64{r.CPUPercent,r.MemoryPercent,r.DiskPercent}{if math.IsNaN(v)||math.IsInf(v,0)||v<0||v>100{return false}};return r.CPUPercent<90&&r.MemoryPercent<90&&r.DiskPercent<95}
