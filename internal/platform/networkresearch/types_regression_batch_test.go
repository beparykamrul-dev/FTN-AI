package networkresearch

import "testing"

func TestResearchRequestRequiresAuthorization(t *testing.T) {
	r := ResearchRequest{ID: "r1", Target: "router-1", Tools: []ToolKind{ToolPing}}
	if r.Validate() == nil { t.Fatal("unauthorized research request must be rejected") }
}

func TestResearchRequestRejectsDuplicateTools(t *testing.T) {
	r := ResearchRequest{ID: "r1", Target: "router-1", Tools: []ToolKind{ToolPing, ToolPing}, Authorized: true}
	if r.Validate() == nil { t.Fatal("duplicate research tools must be rejected") }
}
