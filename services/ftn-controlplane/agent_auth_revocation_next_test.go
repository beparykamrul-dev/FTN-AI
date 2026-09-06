package controlplane

import "testing"

func TestAuthorizeAgentRejectsRevokedIdentity(t *testing.T) {
 e:=AgentIdentity{ServerID:"a",Fingerprint:"f",Enrolled:true}; p:=e; e.Revoked=true; if AuthorizeAgent(e,p) {t.Fatal("revoked identity authorized")}; p=e; p.Revoked=false; e.Revoked=false; p.Revoked=true; if AuthorizeAgent(e,p) {t.Fatal("revoked presented identity authorized")}
}
