package controlplane

import "testing"

func TestPolicyGateRequiresApproval(t *testing.T) {
 a:=Analyze(ServerState{ID:"s1",Healthy:true,Services:[]string{"dns"}})
 p:=PolicyGate{RequireExplicitApproval:true}.Validate(a,true,false)
 if p.Approved || p.Reason!="explicit approval required" {t.Fatalf("unexpected plan: %+v",p)}
}
