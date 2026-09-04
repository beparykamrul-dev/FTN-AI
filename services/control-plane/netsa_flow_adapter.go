package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

// NetSAFlowAdapter accepts normalized JSON records commonly emitted by
// YAF/super_mediator pipelines. It deliberately consumes metadata only; raw
// packets and opaque payloads are never retained or forwarded.
type NetSAFlowAdapter struct{}

type netSAFlowRecord struct {
	ExporterID string `json:"exporter_id"`
	SourceIP string `json:"sourceIPv4Address"`
	DestinationIP string `json:"destinationIPv4Address"`
	SourceIPv6 string `json:"sourceIPv6Address"`
	DestinationIPv6 string `json:"destinationIPv6Address"`
	SourcePort uint16 `json:"sourceTransportPort"`
	DestinationPort uint16 `json:"destinationTransportPort"`
	Protocol uint8 `json:"protocolIdentifier"`
	Bytes uint64 `json:"octetTotalCount"`
	Packets uint64 `json:"packetTotalCount"`
	SamplingRate uint32 `json:"samplingInterval"`
}

func (NetSAFlowAdapter) Name() string { return "netsa-jsonl" }

func (a NetSAFlowAdapter) Decode(ctx context.Context, r io.Reader, exporter string, version uint16) ([]FlowRecord, error) {
	if r == nil { return nil, errors.New("flow_reader_required") }
	if ctx == nil { ctx = context.Background() }
	if version != 5 && version != 9 && version != 10 { return nil, errors.New("unsupported_flow_version") }
	exporter = strings.TrimSpace(exporter)
	if exporter == "" { return nil, errors.New("exporter_required") }
	if net.ParseIP(exporter) == nil { return nil, errors.New("invalid_exporter") }
	out := make([]FlowRecord, 0, 128)
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 4096), 1<<20)
	for scanner.Scan() {
		select { case <-ctx.Done(): return out, ctx.Err(); default: }
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 { continue }
		var raw netSAFlowRecord
		if err := json.Unmarshal(line, &raw); err != nil { return nil, fmt.Errorf("invalid_netsa_json: %w", err) }
		source := raw.SourceIP; if source == "" { source = raw.SourceIPv6 }
		destination := raw.DestinationIP; if destination == "" { destination = raw.DestinationIPv6 }
		if net.ParseIP(source) == nil || net.ParseIP(destination) == nil { return nil, errors.New("flow_source_destination_required") }
		rate := raw.SamplingRate; if rate == 0 { rate = 1 }
		id := strings.TrimSpace(raw.ExporterID); if id == "" { id = exporter }
		out = append(out, FlowRecord{ExporterID: id, Version: version, SourceIP: source, DestinationIP: destination, SourcePort: raw.SourcePort, DestinationPort: raw.DestinationPort, Protocol: raw.Protocol, Bytes: raw.Bytes, Packets: raw.Packets, SamplingRate: rate})
		if len(out) > maxSiLKBatchSize { return nil, errors.New("flow_batch_limit") }
	}
	if err := scanner.Err(); err != nil { return nil, fmt.Errorf("read_netsa_json: %w", err) }
	return out, nil
}

// ParseUnsigned is kept small and strict for text-oriented NetSA adapters.
func ParseUnsigned(s string, bits int) (uint64, error) {
	s = strings.TrimSpace(s)
	if s == "" { return 0, errors.New("number_required") }
	return strconv.ParseUint(s, 10, bits)
}
