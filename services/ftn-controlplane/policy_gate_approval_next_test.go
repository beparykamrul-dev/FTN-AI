package controlplane

import "testing"

func TestPolicyGateRequiresExplicitApproval(t *testing.T) {
 p:=PolicyGate{RequireExplicitApproval:true}; r:=p.Validate(AnalysisResult{ServerID:"s",Healthy:true},true,false); if r.Approved {t.Fatal("unapproved deployment accepted")}
}
