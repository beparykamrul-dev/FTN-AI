package controlplane

import "testing"

func TestPolicyGateRejectsUnhealthyServer(t *testing.T) {
 p:=PolicyGate{RequireHealthy:true}; r:=p.Validate(AnalysisResult{ServerID:"s",Healthy:false},true,true); if r.Approved || r.Reason=="" {t.Fatalf("unexpected plan: %+v",r)}
}
