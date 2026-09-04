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
(observed_at, tenant_id, exporter_id, version, source_ip, destination_ip,
 source_port, destination_port, protocol, bytes, packets, sampling_rate,
 fingerprint, service_id, subscriber_id, main_server_id)
VALUES ($1,$2,$3,$4,$5::inet,$6::inet,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
ON CONFLICT (tenant_id, exporter_id, fingerprint, observed_at) DO NOTHING`

// FTNSiLKStore is the durable database boundary for normalized flow telemetry.
// It stores flow metadata/counters only; raw packet payloads are never persisted.
type FTNSiLKStore struct {
	db         *pgxpool.Pool
	tenantID   string
	clickhouse *FTNClickHouseStore
}

func NewFTNSiLKStore(db *pgxpool.Pool, tenantID string) *FTNSiLKStore {
	clickhouse, err := NewFTNClickHouseStoreFromEnv()
	if err != nil {
		log.Printf("FTN ClickHouse analytics disabled: %v", err)
	}
	return &FTNSiLKStore{db: db, tenantID: strings.TrimSpace(tenantID), clickhouse: clickhouse}
}

func (s *FTNSiLKStore) validateRecord(r FlowRecord) (FlowRecord, *string, *string, error) {
	if s == nil || s.db == nil {
		return FlowRecord{}, nil, nil, errors.New("database_required")
	}
	if strings.TrimSpace(s.tenantID) == "" {
		return FlowRecord{}, nil, nil, errors.New("tenant_required")
	}
	if strings.TrimSpace(r.ExporterID) == "" {
		return FlowRecord{}, nil, nil, errors.New("exporter_required")
	}
	if r.Version != 5 && r.Version != 9 && r.Version != 10 {
		return FlowRecord{}, nil, nil, errors.New("unsupported_flow_version")
	}
	if r.SamplingRate == 0 {
		r.SamplingRate = 1
	}
	// PostgreSQL BIGINT is signed. FlowRecord counters are uint64, so reject
	// values that cannot be represented rather than allowing a wrapping cast.
	if r.Bytes > math.MaxInt64 || r.Packets > math.MaxInt64 {
		return FlowRecord{}, nil, nil, errors.New("flow_counter_overflow")
	}

	var src, dst *string
	if p, err := netip.ParseAddr(strings.TrimSpace(r.SourceIP)); err == nil {
		v := p.String()
		src = &v
	}
	if p, err := netip.ParseAddr(strings.TrimSpace(r.DestinationIP)); err == nil {
		v := p.String()
		dst = &v
	}
	if strings.TrimSpace(r.SourceIP) != "" && src == nil {
		return FlowRecord{}, nil, nil, errors.New("invalid_source_ip")
	}
	if strings.TrimSpace(r.DestinationIP) != "" && dst == nil {
		return FlowRecord{}, nil, nil, errors.New("invalid_destination_ip")
	}
	return r, src, dst, nil
}

func (s *FTNSiLKStore) args(observedAt time.Time, r FlowRecord, src, dst *string, serviceID, subscriberID, mainServerID string) []any {
	if observedAt.IsZero() {
		observedAt = time.Now().UTC()
	}
	fingerprint := fmt.Sprintf("%016x", FlowRecordFingerprint(r))
	return []any{
		observedAt, s.tenantID, r.ExporterID, r.Version, src, dst,
		r.SourcePort, r.DestinationPort, r.Protocol, int64(r.Bytes), int64(r.Packets),
		r.SamplingRate, fingerprint, nullableString(serviceID), nullableString(subscriberID), nullableString(mainServerID),
	}
}

func (s *FTNSiLKStore) Insert(ctx context.Context, observedAt time.Time, r FlowRecord, serviceID, subscriberID, mainServerID string) error {
	r, src, dst, err := s.validateRecord(r)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(ctx, silkFlowInsertSQL, s.args(observedAt, r, src, dst, serviceID, subscriberID, mainServerID)...)
	return err
}

// InsertBatch persists a bounded batch with pgx's pipeline/batch protocol.
// The batch is intentionally capped to prevent an exporter burst from turning
// into an unbounded database request or transaction. Duplicate flow records
// remain harmless through the existing unique dedupe index.
func (s *FTNSiLKStore) InsertBatch(ctx context.Context, observedAt time.Time, records []FlowRecord, serviceID, subscriberID, mainServerID string) error {
	if len(records) == 0 {
		return nil
	}
	if len(records) > maxSiLKBatchSize {
		return fmt.Errorf("flow_batch_too_large: %d > %d", len(records), maxSiLKBatchSize)
	}
	if _, _, _, err := s.validateRecord(records[0]); err != nil {
		return err
	}

	batch := &pgx.Batch{}
	for _, record := range records {
		r, src, dst, err := s.validateRecord(record)
		if err != nil {
			return err
		}
		batch.Queue(silkFlowInsertSQL, s.args(observedAt, r, src, dst, serviceID, subscriberID, mainServerID)...)
	}
	results := s.db.SendBatch(ctx, batch)
	defer results.Close()
	for range records {
		if _, err := results.Exec(); err != nil {
			return err
		}
	}

	if s.clickhouse != nil {
		if err := s.clickhouse.InsertBatch(ctx, observedAt, s.tenantID, records, serviceID, subscriberID, mainServerID); err != nil {
			// Analytics is deliberately non-blocking for the transactional path.
			log.Printf("FTN ClickHouse flow persistence failed: %v", err)
		}
	}
	return nil
}

func nullableString(v string) *string {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	return &v
}
