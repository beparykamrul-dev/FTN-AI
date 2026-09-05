package observability

import("fmt";"strings")
type DPIFlow struct{SourceIP string `json:"source_ip"`;DestinationIP string `json:"destination_ip"`;Protocol string `json:"protocol"`;Application string `json:"application,omitempty"`;Category string `json:"category,omitempty"`;Bytes uint64 `json:"bytes"`;Packets uint64 `json:"packets"`}
type DPIAdapter struct{Engine string `json:"engine"`;Version string `json:"version"`}
func NewDPIAdapter(version string)*DPIAdapter{return &DPIAdapter{Engine:"ndpi",Version:strings.TrimSpace(version)}}
func(a DPIAdapter)Validate()error{if strings.TrimSpace(a.Engine)==""{return fmt.Errorf("DPI engine is required")};if strings.TrimSpace(a.Version)==""{return fmt.Errorf("DPI engine version is required")};return nil}
func NormalizeDPI(v DPIFlow)TrafficSample{return TrafficSample{SourceIP:strings.TrimSpace(v.SourceIP),DestIP:strings.TrimSpace(v.DestinationIP),Protocol:strings.TrimSpace(v.Protocol),App:strings.TrimSpace(v.Application),Bytes:v.Bytes,Packets:v.Packets}}
