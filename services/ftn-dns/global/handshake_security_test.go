package global

import("context";"testing")
func TestProbeRejectsOversizedAddress(t *testing.T){r:=Probe(context.Background(),ResolverProbe{Address:string(make([]byte,2049)),Network:"tcp"},0);if r.Error==""{t.Fatal("expected address rejection")}}
func TestProbeRejectsUnsupportedNetwork(t *testing.T){r:=Probe(context.Background(),ResolverProbe{Address:"127.0.0.1:53",Network:"icmp"},0);if r.Error==""{t.Fatal("expected network rejection")}}
