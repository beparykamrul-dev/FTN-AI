package networkresearch

import "time"

type ToolKind string

const (
	ToolPing ToolKind = "ping"
	ToolTraceroute ToolKind = "traceroute"
	ToolDNS ToolKind = "dns"
	ToolHTTP ToolKind = "http"
	ToolTLS ToolKind = "tls"
	ToolBGP ToolKind = "bgp"
	ToolFlow ToolKind = "flow"
	ToolRoute ToolKind = "route"
	ToolMTU ToolKind = "mtu"
)

type ResearchRequest struct {
	ID string `json:"id"`
	Target string `json:"target"`
	Tools []ToolKind `json:"tools"`
	Authorized bool `json:"authorized"`
	CreatedAt time.Time `json:"created_at"`
}

type ResearchResult struct {
	RequestID string `json:"request_id"`
	Tool ToolKind `json:"tool"`
	Status string `json:"status"`
	Summary string `json:"summary,omitempty"`
	Evidence []string `json:"evidence,omitempty"`
	CheckedAt time.Time `json:"checked_at"`
}
