package observability

import "testing"

func TestDPIAdapterAndFlowNormalization(t *testing.T) {
 if err := NewDPIAdapter("4.0").Validate(); err != nil { t.Fatalf("valid DPI adapter rejected: %v", err) }
 if err := NewDPIAdapter("").Validate(); err == nil { t.Fatal("empty DPI version accepted") }
 s := NormalizeDPI(DPIFlow{SourceIP: " 192.0.2.1 ", DestinationIP: " 198.51.100.1 ", Protocol: " TCP ", Application: " Web "})
 if s.SourceIP != "192.0.2.1" || s.DestIP != "198.51.100.1" || s.Protocol != "tcp" || s.App != "Web" { t.Fatalf("normalized sample=%+v", s) }
}
