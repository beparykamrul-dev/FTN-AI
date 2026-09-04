package main

import (
	"context"
	"strings"
	"testing"
)

func TestNetSAFlowAdapterDecodeYAFStyleJSONL(t *testing.T) {
	input := `{"sourceIPv4Address":"192.0.2.10","destinationIPv4Address":"198.51.100.20","sourceTransportPort":1234,"destinationTransportPort":443,"protocolIdentifier":6,"octetTotalCount":9000,"packetTotalCount":12,"samplingInterval":100}` + "\n"
	records, err := (NetSAFlowAdapter{}).Decode(context.Background(), strings.NewReader(input), "203.0.113.1", 10)
	if err != nil { t.Fatalf("decode: %v", err) }
	if len(records) != 1 { t.Fatalf("records=%d, want 1", len(records)) }
	r := records[0]
	if r.ExporterID != "203.0.113.1" || r.SourceIP != "192.0.2.10" || r.DestinationIP != "198.51.100.20" || r.Bytes != 9000 || r.Packets != 12 || r.SamplingRate != 100 { t.Fatalf("unexpected normalized record: %+v", r) }
}

func TestNetSAFlowAdapterDecodeIPv6(t *testing.T) {
	input := `{"sourceIPv6Address":"2001:db8::1","destinationIPv6Address":"2001:db8::2","protocolIdentifier":17,"octetTotalCount":64,"packetTotalCount":1}`
	records, err := (NetSAFlowAdapter{}).Decode(context.Background(), strings.NewReader(input), "203.0.113.2", 9)
	if err != nil { t.Fatalf("decode: %v", err) }
	if len(records) != 1 || records[0].SourceIP != "2001:db8::1" || records[0].DestinationIP != "2001:db8::2" { t.Fatalf("unexpected IPv6 record: %+v", records) }
	if records[0].SamplingRate != 1 { t.Fatalf("sampling=%d, want 1", records[0].SamplingRate) }
}

func TestNetSAFlowAdapterRejectsInvalidAddress(t *testing.T) {
	input := `{"sourceIPv4Address":"not-an-ip","destinationIPv4Address":"198.51.100.20"}`
	if _, err := (NetSAFlowAdapter{}).Decode(context.Background(), strings.NewReader(input), "203.0.113.1", 10); err == nil { t.Fatal("expected invalid source address error") }
}

func TestNetSAFlowAdapterRejectsOversizedBatch(t *testing.T) {
	line := `{"sourceIPv4Address":"192.0.2.1","destinationIPv4Address":"198.51.100.1"}`
	input := strings.Repeat(line+"\n", maxSiLKBatchSize+1)
	if _, err := (NetSAFlowAdapter{}).Decode(context.Background(), strings.NewReader(input), "203.0.113.1", 10); err == nil { t.Fatal("expected batch limit error") }
}

func TestParseUnsigned(t *testing.T) {
	v, err := ParseUnsigned("42", 16)
	if err != nil || v != 42 { t.Fatalf("value=%d err=%v", v, err) }
	if _, err := ParseUnsigned("-1", 16); err == nil { t.Fatal("expected negative value rejection") }
}
