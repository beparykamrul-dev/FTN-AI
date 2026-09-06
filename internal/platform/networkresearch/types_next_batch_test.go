package networkresearch

import "testing"

func TestResearchRequestAuthorizationAndDuplicates(t *testing.T) {
 r := ResearchRequest{ID: "r1", Target: "example.com", Tools: []ToolKind{ToolDNS, ToolHTTP}, Authorized: true}
 if err := r.Validate(); err != nil { t.Fatalf("valid research request rejected: %v", err) }
 r.Tools = []ToolKind{ToolDNS, ToolDNS}; if err := r.Validate(); err == nil { t.Fatal("duplicate tool accepted") }
 r.Tools = []ToolKind{ToolDNS}; r.Authorized = false; if err := r.Validate(); err == nil { t.Fatal("unauthorized research accepted") }
 if !(ResearchResult{RequestID: "r1", Tool: ToolDNS, Summary: "ok"}).Valid() { t.Fatal("valid research result rejected") }
}
