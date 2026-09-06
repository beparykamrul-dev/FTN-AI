package controlplane

import "testing"

func TestAuthorizeAgentRejectsMismatchedServer(t *testing.T) {
 e:=AgentIdentity{ServerID:"s1",Fingerprint:"fp",Enrolled:true}; p:=e; p.ServerID="s2"
 if AuthorizeAgent(e,p) {t.Fatal("mismatched server authorized")}
}
