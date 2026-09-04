package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

// NFCAPDFlowAdapter accepts metadata exported by nfdump/nfcapd in CSV form.
// It is intentionally an import boundary: no packet payloads are read or stored.
type NFCAPDFlowAdapter struct{}

type nfcapdColumns struct {
	src, dst, sport, dport, proto, bytes, packets, sampling int
}

func (NFCAPDFlowAdapter) Name() string { return "nfcapd-csv" }

func normalizeNFCAPDColumn(s string) string {
	s := strings.ToLower(strings.TrimSpace(s))
	s = strings.Trim(s, "\"'")
	return strings.ReplaceAll(s, "_", "")
}

func findNFCAPDColumn(header []string, aliases ...string) int {
	want := make(map[string]struct{}, len(aliases))
	for _, alias := range aliases {
		want[normalizeNFCAPDColumn(alias)] = struct{}{}
	}
	for i, name := range header {
		if _, ok := want[normalizeNFCAPDColumn(name)]; ok {
			return i
		}
	}
	return -1
}

func parseNFCAPDUint(rec []string, idx int, field string, bits int) (uint64, error) {
	if idx < 0 || idx >= len(rec) {
		return 0, fmt.Errorf("%s_required", field)
	}
	value := strings.TrimSpace(rec[idx])
	if value == "" {
		return 0, fmt.Errorf("%s_required", field)
	}
	v, err := strconv.ParseUint(value, 10, bits)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", field, err)
	}
	return v, nil
}

func (a NFCAPDFlowAdapter) Decode(ctx context.Context, r io.Reader, exporter string, version uint16) ([]FlowRecord, error) {
	if r == nil {
		return nil, errors.New("flow_reader_required")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	exporter = strings.TrimSpace(exporter)
	if exporter == "" || net.ParseIP(exporter) == nil {
		return nil, errors.New("invalid_exporter")
	}
	if version != 5 && version != 9 && version != 10 {
		return nil, errors.New("unsupported_flow_version")
	}

	cr := csv.NewReader(bufio.NewReader(r))
	cr.TrimLeadingSpace = true
	cr.FieldsPerRecord = -1
	out := make([]FlowRecord, 0, 128)
	var columns *nfcapdColumns

	for {
		select {
		case <-ctx.Done():
			return out, ctx.Err()
		default:
		}
		rec, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read_nfcapd_csv: %w", err)
		}
		if len(rec) == 0 {
			continue
		}
		for i := range rec {
			rec[i] = strings.TrimSpace(rec[i])
		}

		// Accept both the stable eight-column FTN format and common nfdump
		// header names. Once a header is recognized, data columns are mapped
		// by name instead of relying on exporter-specific column ordering.
		if columns == nil && len(rec) >= 8 && (strings.EqualFold(rec[0], "ts") || strings.EqualFold(rec[0], "first") || strings.EqualFold(rec[0], "srcip") || strings.EqualFold(rec[0], "sa")) {
			c := &nfcapdColumns{
				src:      findNFCAPDColumn(rec, "srcip", "sa", "sourceip", "sourceipv4address", "sourceipv6address"),
				dst:      findNFCAPDColumn(rec, "dstip", "da", "destinationip", "destinationipv4address", "destinationipv6address"),
				sport:    findNFCAPDColumn(rec, "srcport", "sp", "spt", "sourcetransportport"),
				dport:    findNFCAPDColumn(rec, "dstport", "dp", "dpt", "destinationtransportport"),
				proto:    findNFCAPDColumn(rec, "proto", "pr", "protocol", "protocolidentifier"),
				sampling: findNFCAPDColumn(rec, "sampling", "samplinginterval"),
			}
			if c.src >= 0 && c.dst >= 0 && c.sport >= 0 && c.dport >= 0 && c.proto >= 0 {
				c.bytes = findNFCAPDColumn(rec, "bytes", "ibytes", "octettotalcount")
				c.packets = findNFCAPDColumn(rec, "packets", "ipackets", "packettotalcount")
				if c.bytes >= 0 && c.packets >= 0 {
					columns = c
					continue
				}
			}
		}
		if len(rec) >= 1 && (strings.EqualFold(rec[0], "ts") || strings.EqualFold(rec[0], "first") || strings.EqualFold(rec[0], "srcip") || strings.EqualFold(rec[0], "sa")) {
			continue
		}

		srcIdx, dstIdx, sportIdx, dportIdx, protoIdx := 0, 1, 2, 3, 4
		bytesIdx, packetsIdx, samplingIdx := 5, 6, 7
		if columns != nil {
			srcIdx, dstIdx, sportIdx, dportIdx, protoIdx = columns.src, columns.dst, columns.sport, columns.dport, columns.proto
			bytesIdx, packetsIdx, samplingIdx = columns.bytes, columns.packets, columns.sampling
		}
		if len(rec) <= bytesIdx || len(rec) <= packetsIdx {
			return nil, errors.New("nfcapd_record_too_short")
		}
		src, dst := rec[srcIdx], rec[dstIdx]
		if net.ParseIP(src) == nil || net.ParseIP(dst) == nil {
			return nil, errors.New("flow_source_destination_required")
		}
		sport, err := parseNFCAPDUint(rec, sportIdx, "source_port", 16)
		if err != nil {
			return nil, err
		}
		dport, err := parseNFCAPDUint(rec, dportIdx, "destination_port", 16)
		if err != nil {
			return nil, err
		}
		proto, err := parseNFCAPDUint(rec, protoIdx, "protocol", 8)
		if err != nil {
			return nil, err
		}
		bytes, err := parseNFCAPDUint(rec, bytesIdx, "bytes", 64)
		if err != nil {
			return nil, err
		}
		packets, err := parseNFCAPDUint(rec, packetsIdx, "packets", 64)
		if err != nil {
			return nil, err
		}
		rate := uint64(1)
		if samplingIdx >= 0 && samplingIdx < len(rec) && rec[samplingIdx] != "" {
			rate, err = parseNFCAPDUint(rec, samplingIdx, "sampling_rate", 32)
			if err != nil {
				return nil, err
			}
			if rate == 0 {
				rate = 1
			}
		}
		out = append(out, FlowRecord{
			ExporterID: exporter, Version: version, SourceIP: src, DestinationIP: dst,
			SourcePort: uint16(sport), DestinationPort: uint16(dport), Protocol: uint8(proto),
			Bytes: bytes, Packets: packets, SamplingRate: uint32(rate),
		})
		if len(out) > maxSiLKBatchSize {
			return nil, errors.New("flow_batch_limit")
		}
	}
	return out, nil
}
