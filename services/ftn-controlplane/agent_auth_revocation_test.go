package controlplane

import "testing"

func TestAuthorizeAgentRejectsRevokedIdentity(t *testing.T) {
 e:=AgentIdentity{ServerID:"s1",Fingerprint:"fp",Enrolled:true}; p:=e; p.Revoked=true
 if AuthorizeAgent(e,p) {t.Fatal("revoked agent authorized")}
}
