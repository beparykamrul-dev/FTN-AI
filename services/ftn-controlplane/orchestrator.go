package controlplane

import("math";"strings")

type ServerState struct{ID string;Version string;Healthy bool;CPUPercent float64;RAMPercent float64;DiskPercent float64;Services []string}
type DesiredState struct{ServerID string;Services []string}
type AnalysisResult struct{ServerID string;Healthy bool;Desired DesiredState;Reason string}
func Analyze(s ServerState)AnalysisResult{id:=strings.TrimSpace(s.ID);healthy:=s.Healthy&&validPercent(s.CPUPercent)&&validPercent(s.RAMPercent)&&validPercent(s.DiskPercent)&&HealthyResources(ResourceHealth{CPUPercent:s.CPUPercent,MemoryPercent:s.RAMPercent,DiskPercent:s.DiskPercent});return AnalysisResult{ServerID:id,Healthy:healthy,Desired:DesiredState{ServerID:id,Services:NormalizeServices(s.Services)},Reason:"observed state normalized; deployment requires the control-plane policy gate"}}
func validPercent(v float64)bool{return !math.IsNaN(v)&&!math.IsInf(v,0)&&v>=0&&v<=100}
