package observability

import (
	"fmt"
	"strings"
)

type DPIFlow struct {
	SourceIP      string `json:"source_ip"`
	DestinationIP string `json:"destination_ip"`
	Protocol      string `json:"protocol"`
	Application   string `json:"application,omitempty"`
	Category      string `json:"category,omitempty"`
	Bytes         uint64 `json:"bytes"`
	Packets       uint64 `json:"packets"`
}

// DPIAdapter is the integration boundary for an nDPI-backed collector.
// Classification is performed by the authorized collector/agent; the control
// plane receives normalized metadata rather than raw packet payloads.
type DPIAdapter struct { Engine string `json:"engine"`; Version string `json:"version"` }

func NewDPIAdapter(version string) *DPIAdapter { return &DPIAdapter{Engine: "ndpi", Version: version} }

func (a DPIAdapter) Validate() error {
	if strings.TrimSpace(a.Version) == "" { return fmt.Errorf("DPI engine version is required") }
	return nil
}

func NormalizeDPI(v DPIFlow) TrafficSample {
	return TrafficSample{SourceIP: v.SourceIP, DestIP: v.DestinationIP, Protocol: v.Protocol, App: v.Application, Bytes: v.Bytes, Packets: v.Packets}
}
