package controlplane

import "testing"

func TestPolicyGateRejectsUnknownServer(t *testing.T) {
 a:=Analyze(ServerState{ID:"s1",Healthy:true}); p:=PolicyGate{RequireKnownServer:true}.Validate(a,false,true)
 if p.Approved || p.Reason!="server is not enrolled" {t.Fatalf("unexpected plan: %+v",p)}
}
