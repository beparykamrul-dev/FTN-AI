package main

import (
	"encoding/binary"
	"testing"
)

func TestDispatchNetFlowV9TemplateAndData(t *testing.T) {
	cache := NewFlowTemplateCache()
	// v9 header (20) + template set (16) + data set (20).
	packet := make([]byte, 20+16+20)
	binary.BigEndian.PutUint16(packet[0:2], 9)
	binary.BigEndian.PutUint16(packet[2:4], 2)
	binary.BigEndian.PutUint32(packet[12:16], 7)
	binary.BigEndian.PutUint32(packet[16:20], 42)
	off := 20
	binary.BigEndian.PutUint16(packet[off:off+2], 0)
	binary.BigEndian.PutUint16(packet[off+2:off+4], 16)
	binary.BigEndian.PutUint16(packet[off+4:off+6], 256)
	binary.BigEndian.PutUint16(packet[off+6:off+8], 2)
	binary.BigEndian.PutUint16(packet[off+8:off+10], 8)
	binary.BigEndian.PutUint16(packet[off+10:off+12], 4)
	binary.BigEndian.PutUint16(packet[off+12:off+14], 12)
	binary.BigEndian.PutUint16(packet[off+14:off+16], 4)
	off += 16
	binary.BigEndian.PutUint16(packet[off:off+2], 256)
	binary.BigEndian.PutUint16(packet[off+2:off+4], 20)
	binary.BigEndian.PutUint32(packet[off+4:off+8], 0x0a000001)
	binary.BigEndian.PutUint32(packet[off+8:off+12], 0x08080808)
	binary.BigEndian.PutUint32(packet[off+12:off+16], 1234)
	binary.BigEndian.PutUint32(packet[off+16:off+20], 0)

	got, err := DispatchFlowPacket(packet, "192.0.2.1", cache)
	if err != nil {
		t.Fatalf("dispatch: %v", err)
	}
	if got.Metadata.Version != 9 || got.Metadata.ObservationDomain != 42 || got.Metadata.Sequence != 7 {
		t.Fatalf("bad metadata: %+v", got.Metadata)
	}
	if got.Templates != 1 || len(got.Records) != 1 {
		t.Fatalf("templates=%d records=%d", got.Templates, len(got.Records))
	}
	if got.Records[0].SourceIP != "10.0.0.1" || got.Records[0].DestinationIP != "8.8.8.8" || got.Records[0].Bytes != 1234 {
		t.Fatalf("bad record: %+v", got.Records[0])
	}
}

func TestDispatchIPFIXEnterpriseTemplate(t *testing.T) {
	cache := NewFlowTemplateCache()
	// IPFIX header (16) + template set (24) + data set (12).
	packet := make([]byte, 16+24+12)
	binary.BigEndian.PutUint16(packet[0:2], 10)
	binary.BigEndian.PutUint16(packet[2:4], uint16(len(packet)))
	binary.BigEndian.PutUint32(packet[8:12], 9)
	binary.BigEndian.PutUint32(packet[12:16], 77)
	off := 16
	binary.BigEndian.PutUint16(packet[off:off+2], 2)
	binary.BigEndian.PutUint16(packet[off+2:off+4], 24)
	binary.BigEndian.PutUint16(packet[off+4:off+6], 300)
	binary.BigEndian.PutUint16(packet[off+6:off+8], 3)
	binary.BigEndian.PutUint16(packet[off+8:off+10], 8)
	binary.BigEndian.PutUint16(packet[off+10:off+12], 4)
	binary.BigEndian.PutUint16(packet[off+12:off+14], 12)
	binary.BigEndian.PutUint16(packet[off+14:off+16], 4)
	binary.BigEndian.PutUint16(packet[off+16:off+18], 0x8001)
	binary.BigEndian.PutUint16(packet[off+18:off+20], 4)
	binary.BigEndian.PutUint32(packet[off+20:off+24], 123)

	off += 24
	binary.BigEndian.PutUint16(packet[off:off+2], 300)
	binary.BigEndian.PutUint16(packet[off+2:off+4], 12)
	binary.BigEndian.PutUint32(packet[off+4:off+8], 0x0a000002)
	binary.BigEndian.PutUint32(packet[off+8:off+12], 0x0a000003)

	got, err := DispatchFlowPacket(packet, "192.0.2.2", cache)
	if err != nil {
		t.Fatalf("dispatch: %v", err)
	}
	if got.Metadata.Version != 10 || got.Metadata.ObservationDomain != 77 || got.Metadata.Sequence != 9 {
		t.Fatalf("bad metadata: %+v", got.Metadata)
	}
	if got.Templates != 1 || len(got.Records) != 1 {
		t.Fatalf("templates=%d records=%d", got.Templates, len(got.Records))
	}
	if got.Records[0].SourceIP != "10.0.0.2" || got.Records[0].DestinationIP != "10.0.0.3" {
		t.Fatalf("bad record: %+v", got.Records[0])
	}
}
