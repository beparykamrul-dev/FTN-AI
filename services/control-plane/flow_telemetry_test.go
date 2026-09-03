package main

import (
	"encoding/binary"
	"testing"
)

func TestDecodeNetFlowV5(t *testing.T) {
	p := make([]byte, 24+48)
	binary.BigEndian.PutUint16(p[0:2], 5)
	binary.BigEndian.PutUint16(p[2:4], 1)
	copy(p[24:28], []byte{192, 0, 2, 10})
	copy(p[28:32], []byte{198, 51, 100, 20})
	binary.BigEndian.PutUint32(p[40:44], 100)
	binary.BigEndian.PutUint32(p[44:48], 1000)
	binary.BigEndian.PutUint16(p[56:58], 1234)
	binary.BigEndian.PutUint16(p[58:60], 443)
	p[62] = 6
	got, err := DecodeNetFlowV5(p, "198.51.100.1")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].SourceIP != "192.0.2.10" || got[0].DestinationPort != 443 || got[0].Packets != 100 || got[0].Bytes != 1000 {
		t.Fatalf("unexpected record: %+v", got)
	}
}

func TestFlowTemplateCacheScopedByExporter(t *testing.T) {
	c := NewFlowTemplateCache()
	key := FlowExporterKey{Address: "192.0.2.1", Protocol: "ipfix"}
	if err := c.Put(key, FlowTemplate{ID: 256, Fields: []FlowTemplateField{{IE: 1, Length: 8}}}); err != nil {
		t.Fatal(err)
	}
	if _, ok := c.Get(key, 256); !ok {
		t.Fatal("template not found")
	}
	if _, ok := c.Get(FlowExporterKey{Address: "192.0.2.2", Protocol: "ipfix"}, 256); ok {
		t.Fatal("template leaked across exporters")
	}
}

func TestDecodeIPFIXDataSet(t *testing.T) {
	c := NewFlowTemplateCache()
	key := FlowExporterKey{Address: "192.0.2.1", Protocol: "ipfix"}
	t := FlowTemplate{ID: 256, Fields: []FlowTemplateField{
		{IE: 8, Length: 4}, {IE: 12, Length: 4}, {IE: 7, Length: 2}, {IE: 11, Length: 2},
		{IE: 4, Length: 1}, {IE: 1, Length: 8}, {IE: 2, Length: 8},
	}}
	if err := c.Put(key, t); err != nil {
		t.Fatal(err)
	}
	p := make([]byte, 37)
	copy(p[0:4], []byte{192, 0, 2, 10})
	copy(p[4:8], []byte{198, 51, 100, 20})
	binary.BigEndian.PutUint16(p[8:10], 1234)
	binary.BigEndian.PutUint16(p[10:12], 443)
	p[12] = 6
	binary.BigEndian.PutUint64(p[13:21], 1000)
	binary.BigEndian.PutUint64(p[21:29], 100)
	got, err := DecodeIPFIXDataSet(p, key, 256, c, "192.0.2.1")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Protocol != 6 || got[0].Packets != 100 || got[0].Bytes != 1000 {
		t.Fatalf("unexpected record: %+v", got)
	}
}

func TestNormalizeSampledCounters(t *testing.T) {
	r := NormalizeSampledCounters(FlowRecord{Bytes: 100, Packets: 10, SamplingRate: 100})
	if r.Bytes != 10000 || r.Packets != 1000 {
		t.Fatalf("unexpected normalized counters: %+v", r)
	}
}

func TestFlowPacketLimits(t *testing.T) {
	if _, err := DecodeNetFlowV5(make([]byte, maxFlowPacketSize+1), "x"); err != errFlowPacketTooLarge {
		t.Fatalf("expected size limit, got %v", err)
	}
	if _, err := DecodeNetFlowV5([]byte{0, 5}, "x"); err != errFlowPacketTruncated {
		t.Fatalf("expected truncation, got %v", err)
	}
}
