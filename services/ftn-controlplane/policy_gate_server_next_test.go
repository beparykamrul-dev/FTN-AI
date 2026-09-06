package controlplane

import "testing"

func TestPolicyGateRequiresKnownServer(t *testing.T) {
 p:=PolicyGate{RequireKnownServer:true}; r:=p.Validate(AnalysisResult{ServerID:"s",Healthy:true},false,true); if r.Approved {t.Fatal("unknown server accepted")}
}
