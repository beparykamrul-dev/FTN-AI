package networkresearch

import "testing"

func TestResearchRequestValidate(t *testing.T){r:=ResearchRequest{ID:"r1",Target:"router-1",Tools:[]ToolKind{ToolPing,ToolDNS},Authorized:true};if err:=r.Validate();err!=nil{t.Fatalf("valid request rejected: %v",err)};r.Tools=append(r.Tools,ToolPing);if err:=r.Validate();err==nil{t.Fatal("duplicate tool accepted")}}
func TestResearchResultValidBounds(t *testing.T){r:=ResearchResult{RequestID:"r1",Tool:ToolPing,Status:"ok"};if !r.Valid(){t.Fatal("valid result rejected")}}
