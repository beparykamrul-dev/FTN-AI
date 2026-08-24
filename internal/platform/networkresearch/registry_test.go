package networkresearch

import "testing"

func TestValidateAuthorizedResearch(t *testing.T) {
	r := ResearchRequest{ID:"r1", Target:"ftn-dns", Authorized:true, Tools:[]ToolKind{ToolDNS,ToolTLS,ToolRoute}}
	if err := ValidateRequest(r); err != nil { t.Fatal(err) }
}

func TestRejectUnauthorizedResearch(t *testing.T) {
	r := ResearchRequest{ID:"r1", Target:"example", Authorized:false, Tools:[]ToolKind{ToolPing}}
	if err := ValidateRequest(r); err == nil { t.Fatal("expected authorization error") }
}
