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

func (NFCAPDFlowAdapter) Name() string { return "nfcapd-csv" }

func (a NFCAPDFlowAdapter) Decode(ctx context.Context, r io.Reader, exporter string, version uint16) ([]FlowRecord, error) {
	if r == nil { return nil, errors.New("flow_reader_required") }
	if ctx == nil { ctx = context.Background() }
	exporter = strings.TrimSpace(exporter)
	if exporter == "" || net.ParseIP(exporter) == nil { return nil, errors.New("invalid_exporter") }
	if version != 5 && version != 9 && version != 10 { return nil, errors.New("unsupported_flow_version") }

	cr := csv.NewReader(bufio.NewReader(r))
	cr.TrimLeadingSpace = true
	cr.FieldsPerRecord = -1
	out := make([]FlowRecord, 0, 128)
	for {
		select { case <-ctx.Done(): return out, ctx.Err(); default: }
		rec, err := cr.Read()
		if err == io.EOF { break }
		if err != nil { return nil, fmt.Errorf("read_nfcapd_csv: %w", err) }
		if len(rec) == 0 { continue }
		for i := range rec { rec[i] = strings.TrimSpace(rec[i]) }
		if strings.EqualFold(rec[0], "ts") || strings.EqualFold(rec[0], "first") || strings.EqualFold(rec[0], "srcip") { continue }
		if len(rec) < 8 { return nil, errors.New("nfcapd_record_too_short") }

		src, dst := rec[0], rec[1]
		if net.ParseIP(src) == nil || net.ParseIP(dst) == nil { return nil, errors.New("flow_source_destination_required") }
		sport, err := strconv.ParseUint(rec[2], 10, 16); if err != nil { return nil, fmt.Errorf("source_port: %w", err) }
		dport, err := strconv.ParseUint(rec[3], 10, 16); if err != nil { return nil, fmt.Errorf("destination_port: %w", err) }
		proto, err := strconv.ParseUint(rec[4], 10, 8); if err != nil { return nil, fmt.Errorf("protocol: %w", err) }
		bytes, err := strconv.ParseUint(rec[5], 10, 64); if err != nil { return nil, fmt.Errorf("bytes: %w", err) }
		packets, err := strconv.ParseUint(rec[6], 10, 64); if err != nil { return nil, fmt.Errorf("packets: %w", err) }
		rate := uint64(1)
		if rec[7] != "" { rate, err = strconv.ParseUint(rec[7], 10, 32); if err != nil { return nil, fmt.Errorf("sampling_rate: %w", err) }; if rate == 0 { rate = 1 } }
		out = append(out, FlowRecord{ExporterID: exporter, Version: version, SourceIP: src, DestinationIP: dst, SourcePort: uint16(sport), DestinationPort: uint16(dport), Protocol: uint8(proto), Bytes: bytes, Packets: packets, SamplingRate: uint32(rate)})
		if len(out) > maxSiLKBatchSize { return nil, errors.New("flow_batch_limit") }
	}
	return out, nil
}
