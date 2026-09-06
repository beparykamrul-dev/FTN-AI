package controlplane

import "testing"

func TestAnalyzeCopiesServices(t *testing.T) {
 services:=[]string{"dns","mesh"}; r:=Analyze(ServerState{ID:"s1",Services:services}); services[0]="changed"
 if r.Desired.Services[0]!="dns" {t.Fatalf("analysis retained caller slice: %v",r.Desired.Services)}
}
