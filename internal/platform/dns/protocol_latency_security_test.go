package dns

import "testing"
func TestProbeProtocolLatencyRejectsProtocol(t *testing.T){r:=ProbeProtocolLatency(nil,"icmp","127.0.0.1:53",0);if r.Error==""{t.Fatal("expected protocol rejection")}}
func TestProbeProtocolLatencyRejectsHugeAddress(t *testing.T){r:=ProbeProtocolLatency(nil,"tcp",string(make([]byte,2049)),0);if r.Error==""{t.Fatal("expected address rejection")}}
