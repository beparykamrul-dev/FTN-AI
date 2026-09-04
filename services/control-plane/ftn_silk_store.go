package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"net/netip"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const maxSiLKBatchSize = 1000

const silkFlowInsertSQL = `
INSERT INTO ftn_silk_flow_records
(observed_at, first_seen, last_seen, tenant_id, exporter_id, version, source_ip, destination_ip,
 source_port, destination_port, protocol, bytes, packets, sampling_rate,
 fingerprint, service_id, subscriber_id, main_server_id)
VALUES ($1,$2,$3,$4,$5,$6,$7::inet,$8::inet,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)
ON CONFLICT (tenant_id, exporter_id, fingerprint, observed_at) DO NOTHING`

type FTNSiLKStore struct { db *pgxpool.Pool; tenantID string; clickhouse *FTNClickHouseStore }

func NewFTNSiLKStore(db *pgxpool.Pool, tenantID string) *FTNSiLKStore {
	clickhouse, err := NewFTNClickHouseStoreFromEnv()
	if err != nil { log.Printf("FTN ClickHouse analytics disabled: %v", err) }
	return &FTNSiLKStore{db: db, tenantID: strings.TrimSpace(tenantID), clickhouse: clickhouse}
}

func (s *FTNSiLKStore) validateRecord(r FlowRecord) (FlowRecord, *string, *string, error) {
	if s == nil || s.db == nil { return FlowRecord{}, nil, nil, errors.New("database_required") }
	if strings.TrimSpace(s.tenantID) == "" { return FlowRecord{}, nil, nil, errors.New("tenant_required") }
	if strings.TrimSpace(r.ExporterID) == "" { return FlowRecord{}, nil, nil, errors.New("exporter_required") }
	if r.Version != 5 && r.Version != 9 && r.Version != 10 { return FlowRecord{}, nil, nil, errors.New("unsupported_flow_version") }
	if r.SamplingRate == 0 { r.SamplingRate = 1 }
	if r.Bytes > math.MaxInt64 || r.Packets > math.MaxInt64 { return FlowRecord{}, nil, nil, errors.New("flow_counter_overflow") }
	var src, dst *string
	if p, err := netip.ParseAddr(strings.TrimSpace(r.SourceIP)); err == nil { v := p.String(); src = &v }
	if p, err := netip.ParseAddr(strings.TrimSpace(r.DestinationIP)); err == nil { v := p.String(); dst = &v }
	if strings.TrimSpace(r.SourceIP) != "" && src == nil { return FlowRecord{}, nil, nil, errors.New("invalid_source_ip") }
	if strings.TrimSpace(r.DestinationIP) != "" && dst == nil { return FlowRecord{}, nil, nil, errors.New("invalid_destination_ip") }
	return r, src, dst, nil
}

func (s *FTNSiLKStore) args(observedAt time.Time, firstSeen, lastSeen *time.Time, r FlowRecord, src, dst *string, serviceID, subscriberID, mainServerID string) []any {
	if observedAt.IsZero() { observedAt = time.Now().UTC() }
	return []any{observedAt.UTC(), firstSeen, lastSeen, s.tenantID, r.ExporterID, r.Version, src, dst, r.SourcePort, r.DestinationPort, r.Protocol, int64(r.Bytes), int64(r.Packets), r.SamplingRate, fmt.Sprintf("%016x", FlowRecordFingerprint(r)), nullableString(serviceID), nullableString(subscriberID), nullableString(mainServerID)}
}

func (s *FTNSiLKStore) InsertObservation(ctx context.Context, o FlowObservation, serviceID, subscriberID, mainServerID string) error {
	if err := o.Validate(); err != nil { return err }
	r, src, dst, err := s.validateRecord(o.FlowRecord)
	if err != nil { return err }
	var first, last *time.Time
	if !o.FirstSeen.IsZero() { v := o.FirstSeen.UTC(); first = &v }
	if !o.LastSeen.IsZero() { v := o.LastSeen.UTC(); last = &v }
	_, err = s.db.Exec(ctx, silkFlowInsertSQL, s.args(o.ObservedAt, first, last, r, src, dst, serviceID, subscriberID, mainServerID)...)
	return err
}

func (s *FTNSiLKStore) Insert(ctx context.Context, observedAt time.Time, r FlowRecord, serviceID, subscriberID, mainServerID string) error {
	o, err := normalizeFlowObservation(r, observedAt, time.Time{}, time.Time{})
	if err != nil { return err }
	return s.InsertObservation(ctx, o, serviceID, subscriberID, mainServerID)
}

func (s *FTNSiLKStore) InsertBatchObservations(ctx context.Context, observations []FlowObservation, serviceID, subscriberID, mainServerID string) error {
	if len(observations) == 0 { return nil }
	if len(observations) > maxSiLKBatchSize { return fmt.Errorf("flow_batch_too_large: %d > %d", len(observations), maxSiLKBatchSize) }
	batch := &pgx.Batch{}
	for _, o := range observations {
		if err := o.Validate(); err != nil { return err }
		r, src, dst, err := s.validateRecord(o.FlowRecord)
		if err != nil { return err }
		var first, last *time.Time
		if !o.FirstSeen.IsZero() { v := o.FirstSeen.UTC(); first = &v }
		if !o.LastSeen.IsZero() { v := o.LastSeen.UTC(); last = &v }
		batch.Queue(silkFlowInsertSQL, s.args(o.ObservedAt, first, last, r, src, dst, serviceID, subscriberID, mainServerID)...)
	}
	results := s.db.SendBatch(ctx, batch)
	defer results.Close()
	for range observations { if _, err := results.Exec(); err != nil { return err } }
	if s.clickhouse != nil {
		if err := s.clickhouse.InsertBatchObservations(ctx, s.tenantID, observations, serviceID, subscriberID, mainServerID); err != nil { log.Printf("FTN ClickHouse flow persistence failed: %v", err) }
	}
	return nil
}

func (s *FTNSiLKStore) InsertBatch(ctx context.Context, observedAt time.Time, records []FlowRecord, serviceID, subscriberID, mainServerID string) error {
	if len(records) == 0 { return nil }
	observations := make([]FlowObservation, 0, len(records))
	for _, r := range records {
		o, err := normalizeFlowObservation(r, observedAt, time.Time{}, time.Time{})
		if err != nil { return err }
		observations = append(observations, o)
	}
	return s.InsertBatchObservations(ctx, observations, serviceID, subscriberID, mainServerID)
}

func nullableString(v string) *string { v = strings.TrimSpace(v); if v == "" { return nil }; return &v }
