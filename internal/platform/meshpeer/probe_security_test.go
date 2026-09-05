package meshpeer

import "testing"
func TestProbeRejectsOversizedEndpoint(t *testing.T){r:=Probe(nil,string(make([]byte,2049)),"",nil,0);if r.Error==""{t.Fatal("expected endpoint rejection")}}
func TestProbeRequiresContext(t *testing.T){r:=Probe(nil,"127.0.0.1:1","",nil,0);if r.Error!="context is required"{t.Fatalf("error=%q",r.Error)}}
