package main

import (
	"encoding/binary"
	"errors"
	"net"
)

var (
	errFlowV9HeaderTruncated   = errors.New("netflow v9 header truncated")
	errFlowIPFIXHeaderTruncated = errors.New("ipfix header truncated")
	errFlowSetInvalid          = errors.New("invalid flow set")
	errFlowTemplateInvalid     = errors.New("invalid flow template")
	errFlowVariableField       = errors.New("unsupported variable-length flow field")
)

type FlowPacketMetadata struct {
	Exporter   string
	Protocol   string
	Version    uint16
	Sequence   uint32
	ObservationDomain uint32
}

type FlowDispatchResult struct {
	Metadata FlowPacketMetadata
	Records  []FlowRecord
	Templates int
}

// DispatchFlowPacket fully separates protocol framing from persistence. It accepts
// NetFlow v5/v9 and IPFIX (v10), updates the shared template cache, and never retains
// raw packet bytes after returning.
func DispatchFlowPacket(packet []byte, exporter string, cache *FlowTemplateCache) (FlowDispatchResult, error) {
	if cache == nil {
		return FlowDispatchResult{}, errors.New("template cache required")
	}
	if len(packet) < 2 || len(packet) > maxFlowPacketSize {
		return FlowDispatchResult{}, errFlowPacketTruncated
	}
	version := binary.BigEndian.Uint16(packet[:2])
	switch version {
	case 5:
		r, err := DecodeNetFlowV5(packet, exporter)
		return FlowDispatchResult{Metadata: FlowPacketMetadata{Exporter: exporter, Protocol: "netflow5", Version: 5}, Records: r}, err
	case 9:
		return dispatchNetFlowV9(packet, exporter, cache)
	case 10:
		return dispatchIPFIX(packet, exporter, cache)
	default:
		return FlowDispatchResult{}, errFlowUnsupportedVersion
	}
}

func dispatchNetFlowV9(packet []byte, exporter string, cache *FlowTemplateCache) (FlowDispatchResult, error) {
	if len(packet) < 20 { return FlowDispatchResult{}, errFlowV9HeaderTruncated }
	count := int(binary.BigEndian.Uint16(packet[2:4]))
	if count > maxRecordsPerPacket { return FlowDispatchResult{}, errFlowInvalidCount }
	sequence := binary.BigEndian.Uint32(packet[12:16])
	sourceID := binary.BigEndian.Uint32(packet[16:20])
	key := FlowExporterKey{Address: exporter, Protocol: "netflow9"}
	result := FlowDispatchResult{Metadata: FlowPacketMetadata{Exporter: exporter, Protocol: "netflow9", Version: 9, Sequence: sequence, ObservationDomain: sourceID}}
	off := 20
	sets := 0
	for off < len(packet) {
		if off+4 > len(packet) { return FlowDispatchResult{}, errFlowSetInvalid }
		setID := binary.BigEndian.Uint16(packet[off:off+2])
		setLen := int(binary.BigEndian.Uint16(packet[off+2:off+4]))
		if setLen < 4 || off+setLen > len(packet) { return FlowDispatchResult{}, errFlowSetInvalid }
		body := packet[off+4 : off+setLen]
		switch {
		case setID == 0:
			n, err := parseV9TemplateSet(body, key, cache)
			if err != nil { return FlowDispatchResult{}, err }
			result.Templates += n
		case setID >= 256:
			t, ok := cache.Get(key, setID)
			if !ok { return FlowDispatchResult{}, errFlowTemplateMissing }
			r, err := decodeFlowRecords(body, t, exporter, 9)
			if err != nil { return FlowDispatchResult{}, err }
			result.Records = append(result.Records, r...)
			if len(result.Records) > maxRecordsPerPacket { return FlowDispatchResult{}, errFlowInvalidCount }
		}
		sets++
		off += setLen
	}
	_ = sets
	return result, nil
}

func parseV9TemplateSet(body []byte, key FlowExporterKey, cache *FlowTemplateCache) (int, error) {
	off := 0
	count := 0
	for off < len(body) {
		if off+4 > len(body) { return 0, errFlowTemplateInvalid }
		id := binary.BigEndian.Uint16(body[off:off+2])
		fields := int(binary.BigEndian.Uint16(body[off+2:off+4]))
		off += 4
		if id < 256 || fields <= 0 || fields > 64 { return 0, errFlowTemplateInvalid }
		need := fields * 4
		if off+need > len(body) { return 0, errFlowTemplateInvalid }
		f := make([]FlowTemplateField, 0, fields)
		for i := 0; i < fields; i++ {
			ie := binary.BigEndian.Uint16(body[off:off+2])
			ln := binary.BigEndian.Uint16(body[off+2:off+4])
			if ln == 0 { return 0, errFlowTemplateInvalid }
			f = append(f, FlowTemplateField{IE: ie, Length: ln})
			off += 4
		}
		if err := cache.Put(key, FlowTemplate{ID: id, Fields: f}); err != nil { return 0, err }
		count++
	}
	return count, nil
}

func dispatchIPFIX(packet []byte, exporter string, cache *FlowTemplateCache) (FlowDispatchResult, error) {
	if len(packet) < 16 { return FlowDispatchResult{}, errFlowIPFIXHeaderTruncated }
	messageLen := int(binary.BigEndian.Uint16(packet[2:4]))
	if messageLen < 16 || messageLen > len(packet) { return FlowDispatchResult{}, errFlowPacketTruncated }
	sequence := binary.BigEndian.Uint32(packet[8:12])
	domain := binary.BigEndian.Uint32(packet[12:16])
	key := FlowExporterKey{Address: exporter, Protocol: "ipfix"}
	result := FlowDispatchResult{Metadata: FlowPacketMetadata{Exporter: exporter, Protocol: "ipfix", Version: 10, Sequence: sequence, ObservationDomain: domain}}
	off := 16
	for off < messageLen {
		if off+4 > messageLen { return FlowDispatchResult{}, errFlowSetInvalid }
		setID := binary.BigEndian.Uint16(packet[off:off+2])
		setLen := int(binary.BigEndian.Uint16(packet[off+2:off+4]))
		if setLen < 4 || off+setLen > messageLen { return FlowDispatchResult{}, errFlowSetInvalid }
		body := packet[off+4 : off+setLen]
		switch {
		case setID == 2:
			n, err := parseIPFIXTemplateSet(body, key, cache)
			if err != nil { return FlowDispatchResult{}, err }
			result.Templates += n
		case setID >= 256:
			t, ok := cache.Get(key, setID)
			if !ok { return FlowDispatchResult{}, errFlowTemplateMissing }
			r, err := decodeFlowRecords(body, t, exporter, 10)
			if err != nil { return FlowDispatchResult{}, err }
			result.Records = append(result.Records, r...)
			if len(result.Records) > maxRecordsPerPacket { return FlowDispatchResult{}, errFlowInvalidCount }
		}
		off += setLen
	}
	return result, nil
}

func parseIPFIXTemplateSet(body []byte, key FlowExporterKey, cache *FlowTemplateCache) (int, error) {
	off := 0
	count := 0
	for off < len(body) {
		if off+4 > len(body) { return 0, errFlowTemplateInvalid }
		id := binary.BigEndian.Uint16(body[off:off+2])
		fields := int(binary.BigEndian.Uint16(body[off+2:off+4]))
		off += 4
		if id < 256 || fields <= 0 || fields > 64 { return 0, errFlowTemplateInvalid }
		f := make([]FlowTemplateField, 0, fields)
		for i := 0; i < fields; i++ {
			if off+4 > len(body) { return 0, errFlowTemplateInvalid }
			rawIE := binary.BigEndian.Uint16(body[off:off+2])
			ln := binary.BigEndian.Uint16(body[off+2:off+4])
			off += 4
			enterprise := uint32(0)
			ie := rawIE & 0x7fff
			if rawIE&0x8000 != 0 {
				if off+4 > len(body) { return 0, errFlowTemplateInvalid }
				enterprise = binary.BigEndian.Uint32(body[off:off+4])
				off += 4
			}
			if ln == 0 { return 0, errFlowTemplateInvalid }
			f = append(f, FlowTemplateField{IE: ie, Length: ln, Enterprise: enterprise})
		}
		if err := cache.Put(key, FlowTemplate{ID: id, Fields: f}); err != nil { return 0, err }
		count++
	}
	return count, nil
}

func decodeFlowRecords(body []byte, t FlowTemplate, exporter string, version uint16) ([]FlowRecord, error) {
	recordLen := 0
	for _, f := range t.Fields {
		if f.Length == 65535 { return nil, errFlowVariableField }
		recordLen += int(f.Length)
	}
	if recordLen <= 0 { return nil, errFlowTemplateInvalid }
	if len(body)%recordLen != 0 { return nil, errFlowPacketTruncated }
	count := len(body) / recordLen
	if count > maxRecordsPerPacket { return nil, errFlowInvalidCount }
	out := make([]FlowRecord, 0, count)
	for n := 0; n < count; n++ {
		r := FlowRecord{ExporterID: exporter, Version: version, SamplingRate: 1}
		off := n * recordLen
		for _, f := range t.Fields {
			v := body[off : off+int(f.Length)]
			off += int(f.Length)
			if f.Enterprise != 0 { continue }
			switch f.IE {
			case 1: r.Bytes = decodeUint(v)
			case 2: r.Packets = decodeUint(v)
			case 4: if len(v) > 0 { r.Protocol = v[len(v)-1] }
			case 7: if len(v) <= 2 { r.SourcePort = uint16(decodeUint(v)) }
			case 8: if len(v) == 4 { r.SourceIP = net.IP(v).String() }
			case 11: if len(v) <= 2 { r.DestinationPort = uint16(decodeUint(v)) }
			case 12: if len(v) == 4 { r.DestinationIP = net.IP(v).String() }
			case 27: if len(v) == 16 { r.SourceIP = net.IP(v).String() }
			case 28: if len(v) == 16 { r.DestinationIP = net.IP(v).String() }
			case 34: r.SamplingRate = uint32(decodeUint(v))
			}
		}
		if r.SamplingRate == 0 { r.SamplingRate = 1 }
		out = append(out, r)
	}
	return out, nil
}

// NormalizeFlowPacketMetadata creates the stable exporter/cache key used by the dispatcher.
func NormalizeFlowPacketMetadata(exporter, protocol string) FlowExporterKey {
	return FlowExporterKey{Address: exporter, Protocol: protocol}
}
