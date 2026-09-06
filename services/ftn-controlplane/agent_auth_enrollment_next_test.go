package controlplane

import "testing"

func TestAuthorizeAgentRequiresEnrollment(t *testing.T) {
 e:=AgentIdentity{ServerID:"a",Fingerprint:"f",Enrolled:true}; p:=e; p.Enrolled=false; if AuthorizeAgent(e,p) {t.Fatal("unenrolled agent authorized")}
}
