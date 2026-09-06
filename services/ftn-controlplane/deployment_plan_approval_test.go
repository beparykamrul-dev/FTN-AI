package controlplane

import "testing"

func TestVerifyEnvelopeRequiresApproval(t *testing.T) {
 e:=SealPlan(DeploymentPlan{ServerID:"s1",Services:[]string{"dns"}})
 if VerifyEnvelope(e) {t.Fatal("unapproved envelope verified")}
}
