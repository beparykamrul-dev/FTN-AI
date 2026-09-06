package controlplane

import "testing"

func TestAnalyzeCopiesServices(t *testing.T) {
 s:=ServerState{ID:"s",Healthy:true,Services:[]string{"dns"}}; r:=Analyze(s); s.Services[0]="changed"; if r.Desired.Services[0]!="dns" {t.Fatal("analysis retained caller slice")}
}
