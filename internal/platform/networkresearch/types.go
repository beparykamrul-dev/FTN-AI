package networkresearch

import("fmt";"strings";"time")
type ToolKind string
const(ToolPing ToolKind="ping";ToolTraceroute ToolKind="traceroute";ToolDNS ToolKind="dns";ToolHTTP ToolKind="http";ToolTLS ToolKind="tls";ToolBGP ToolKind="bgp";ToolFlow ToolKind="flow";ToolRoute ToolKind="route";ToolMTU ToolKind="mtu")
type ResearchRequest struct{ID string `json:"id"`;Target string `json:"target"`;Tools []ToolKind `json:"tools"`;Authorized bool `json:"authorized"`;CreatedAt time.Time `json:"created_at"`}
type ResearchResult struct{RequestID string `json:"request_id"`;Tool ToolKind `json:"tool"`;Status string `json:"status"`;Summary string `json:"summary,omitempty"`;Evidence []string `json:"evidence,omitempty"`;CheckedAt time.Time `json:"checked_at"`}
func(r ResearchRequest)Validate()error{if strings.TrimSpace(r.ID)==""||len(strings.TrimSpace(r.ID))>256{return fmt.Errorf("invalid research request id")};if strings.TrimSpace(r.Target)==""||len(strings.TrimSpace(r.Target))>1024{return fmt.Errorf("invalid research target")};if len(r.Tools)==0||len(r.Tools)>32{return fmt.Errorf("invalid research tool set")};seen:=map[ToolKind]struct{}{};for _,tool:=range r.Tools{if _,ok:=seen[tool];ok{return fmt.Errorf("duplicate research tool")};seen[tool]=struct{}{}};if !r.Authorized{return fmt.Errorf("research request is not authorized")};return nil}
func(r ResearchResult)Valid()bool{return strings.TrimSpace(r.RequestID)!=""&&len(strings.TrimSpace(r.RequestID))<=256&&strings.TrimSpace(string(r.Tool))!=""&&len(r.Summary)<=8192&&len(r.Evidence)<=1024}
