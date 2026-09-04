package main

import (
	"context"
	"strings"
	"testing"
)

func TestNFCAPDFlowAdapterDecodeCSV(t *testing.T) {
	input := "srcip,dstip,srcport,dstport,proto,bytes,packets,sampling\n192.0.2.10,198.51.100.20,1234,443,6,9000,12,100\n"
	records, err := (NFCAPDFlowAdapter{}).Decode(context.Background(), strings.NewReader(input), "203.0.113.1", 10)
	if err != nil { t.Fatalf("decode: %v", err) }
	if len(records) != 1 { t.Fatalf("records=%d, want 1", len(records)) }
	r := records[0]
	if r.SourceIP != "192.0.2.10" || r.DestinationIP != "198.51.100.20" || r.SourcePort != 1234 || r.DestinationPort != 443 || r.Protocol != 6 || r.Bytes != 9000 || r.Packets != 12 || r.SamplingRate != 100 { t.Fatalf("unexpected record: %+v", r) }
}

func TestNFCAPDFlowAdapterMapsNFDumpAliases(t *testing.T) {
	input := "ts,sa,da,sp,dp,pr,ibytes,ipackets,samplingInterval\n2026-09-04T12:00:00Z,192.0.2.10,198.51.100.20,1234,443,6,9000,12,100\n"
	records, err := (NFCAPDFlowAdapter{}).Decode(context.Background(), strings.NewReader(input), "203.0.113.1", 10)
	if err != nil { t.Fatalf("decode: %v", err) }
	if len(records) != 1 { t.Fatalf("records=%d, want 1", len(records)) }
	if records[0].Bytes != 9000 || records[0].Packets != 12 || records[0].SamplingRate != 100 { t.Fatalf("unexpected record: %+v", records[0]) }
}

func TestNFCAPDFlowAdapterDefaultsSampling(t *testing.T) {
	input := "192.0.2.1,198.51.100.1,1,2,17,64,1,0\n"
	records, err := (NFCAPDFlowAdapter{}).Decode(context.Background(), strings.NewReader(input), "203.0.113.1", 9)
	if err != nil { t.Fatalf("decode: %v", err) }
	if len(records) != 1 || records[0].SamplingRate != 1 { t.Fatalf("unexpected records: %+v", records) }
}

func TestNFCAPDFlowAdapterRejectsInvalidRecord(t *testing.T) {
	input := "not-an-ip,198.51.100.1,1,2,6,10,1,1\n"
	if _, err := (NFCAPDFlowAdapter{}).Decode(context.Background(), strings.NewReader(input), "203.0.113.1", 10); err == nil { t.Fatal("expected invalid address error") }
}
