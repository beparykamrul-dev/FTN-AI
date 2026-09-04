package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// DecodeObservations preserves exporter timestamps when an nfdump CSV export
// contains first/last (or ts/te) columns. It uses the existing bounded flow
// decoder for the normalized counters, then attaches the source timestamps.
func (a NFCAPDFlowAdapter) DecodeObservations(ctx context.Context, r io.Reader, exporter string, version uint16) ([]FlowObservation, error) {
	if r == nil { return nil, errors.New("flow_reader_required") }
	if ctx == nil { ctx = context.Background() }
	data, err := io.ReadAll(io.LimitReader(r, 16<<20))
	if err != nil { return nil, fmt.Errorf("read_nfcapd_csv: %w", err) }
	if len(data) == 0 { return nil, nil }

	cr := csv.NewReader(bufio.NewReader(strings.NewReader(string(data))))
	cr.TrimLeadingSpace = true; cr.FieldsPerRecord = -1
	records, err := readNFCAPDRecords(ctx, cr)
	if err != nil { return nil, err }
	if len(records) == 0 { return nil, nil }

	firstIdx, lastIdx := -1, -1
	if len(records.header) > 0 {
		firstIdx = findNFCAPDColumn(records.header, "first", "ts", "starttime")
		lastIdx = findNFCAPDColumn(records.header, "last", "te", "endtime", "end")
	}
	if firstIdx < 0 && lastIdx < 0 {
		flows, err := a.Decode(ctx, strings.NewReader(string(data)), exporter, version)
		if err != nil { return nil, err }
		out := make([]FlowObservation, 0, len(flows))
		now := time.Now().UTC()
		for _, f := range flows { o, e := normalizeFlowObservation(f, now, time.Time{}, time.Time{}); if e != nil { return nil, e }; out = append(out, o) }
		return out, nil
	}

	// Re-decode normalized records through the canonical adapter, while using
	// the same data for timestamp extraction. This keeps counter validation in
	// one place and prevents the observation path from diverging.
	flows, err := a.Decode(ctx, strings.NewReader(string(data)), exporter, version)
	if err != nil { return nil, err }
	if len(flows) != len(records.rows) { return nil, errors.New("nfcapd_record_alignment_failed") }
	out := make([]FlowObservation, 0, len(flows))
	for i, f := range flows {
		var first, last time.Time
		if firstIdx >= 0 { first, err = parseNFCAPDTime(records.rows[i][firstIdx]); if err != nil { return nil, fmt.Errorf("first_seen: %w", err) } }
		if lastIdx >= 0 { last, err = parseNFCAPDTime(records.rows[i][lastIdx]); if err != nil { return nil, fmt.Errorf("last_seen: %w", err) } }
		observed := last; if observed.IsZero() { observed = first }; if observed.IsZero() { observed = time.Now().UTC() }
		o, err := normalizeFlowObservation(f, observed, first, last); if err != nil { return nil, err }; out = append(out, o)
	}
	return out, nil
}

type nfcapdRows struct { header []string; rows [][]string }
func readNFCAPDRecords(ctx context.Context, cr *csv.Reader) (nfcapdRows, error) {
	var out nfcapdRows
	for {
		select { case <-ctx.Done(): return out, ctx.Err(); default: }
		rec, err := cr.Read(); if err == io.EOF { return out, nil }; if err != nil { return out, fmt.Errorf("read_nfcapd_csv: %w", err) }
		for i := range rec { rec[i] = strings.TrimSpace(rec[i]) }
		if len(rec) == 0 { continue }
		if len(out.header) == 0 && (strings.EqualFold(rec[0], "first") || strings.EqualFold(rec[0], "ts") || strings.EqualFold(rec[0], "srcip") || strings.EqualFold(rec[0], "sa")) { out.header = append([]string(nil), rec...); continue }
		if len(out.header) > 0 { if strings.EqualFold(rec[0], "first") || strings.EqualFold(rec[0], "ts") || strings.EqualFold(rec[0], "srcip") || strings.EqualFold(rec[0], "sa") { continue }; out.rows = append(out.rows, rec) }
	}
}

func parseNFCAPDTime(raw string) (time.Time, error) {
	raw = strings.TrimSpace(strings.Trim(raw, "\"'")); if raw == "" { return time.Time{}, nil }
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02 15:04:05.999999999Z07:00", "2006-01-02 15:04:05"} {
		if t, err := time.Parse(layout, raw); err == nil { return t.UTC(), nil }
	}
	if v, err := strconv.ParseInt(raw, 10, 64); err == nil {
		if len(strings.TrimLeft(raw, "-")) >= 13 { return time.UnixMilli(v).UTC(), nil }
		return time.Unix(v, 0).UTC(), nil
	}
	return time.Time{}, fmt.Errorf("invalid_timestamp")
}
