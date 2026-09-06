package controlplane

import "testing"

func TestPolicyGateRejectsUnhealthyServer(t *testing.T) {
 a:=Analyze(ServerState{ID:"s1",Healthy:false}); p:=PolicyGate{RequireHealthy:true}.Validate(a,true,true)
 if p.Approved || p.Reason!="server health policy rejected deployment" {t.Fatalf("unexpected plan: %+v",p)}
}
