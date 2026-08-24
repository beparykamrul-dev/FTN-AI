package networkresearch

import "fmt"

var supportedTools = map[ToolKind]bool{
	ToolPing:true, ToolTraceroute:true, ToolDNS:true, ToolHTTP:true, ToolTLS:true,
	ToolBGP:true, ToolFlow:true, ToolRoute:true, ToolMTU:true,
}

func ValidateRequest(r ResearchRequest) error {
	if r.ID == "" || r.Target == "" { return fmt.Errorf("research id and target are required") }
	if !r.Authorized { return fmt.Errorf("network research requires authorization") }
	if len(r.Tools) == 0 { return fmt.Errorf("at least one research tool is required") }
	seen := map[ToolKind]bool{}
	for _, t := range r.Tools {
		if !supportedTools[t] { return fmt.Errorf("unsupported research tool: %s", t) }
		if seen[t] { return fmt.Errorf("duplicate research tool: %s", t) }
		seen[t] = true
	}
	return nil
}

func SupportedTools() []ToolKind {
	return []ToolKind{ToolPing,ToolTraceroute,ToolDNS,ToolHTTP,ToolTLS,ToolBGP,ToolFlow,ToolRoute,ToolMTU}
}
