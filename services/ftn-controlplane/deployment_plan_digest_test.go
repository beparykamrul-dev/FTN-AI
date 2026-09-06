package controlplane

import "testing"

func TestVerifyEnvelopeDetectsPlanMutation(t *testing.T) {
 e:=SealPlan(DeploymentPlan{ServerID:"s1",Services:[]string{"dns"},Approved:true})
 e.Plan.Services[0]="mesh"
 if VerifyEnvelope(e) {t.Fatal("mutated envelope verified")}
}
