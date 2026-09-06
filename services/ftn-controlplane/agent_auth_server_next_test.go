package controlplane

import "testing"

func TestAuthorizeAgentRequiresSameServer(t *testing.T) {
 e:=AgentIdentity{ServerID:"a",Fingerprint:"f",Enrolled:true}; p:=e; p.ServerID="b"; if AuthorizeAgent(e,p) {t.Fatal("cross-server identity authorized")}
}
