package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/netip"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// FTNSiLKStore is the durable database boundary for normalized flow telemetry.
// It stores flow metadata/counters only; raw packet payloads are never persisted.
type FTNSiLKStore struct {
	db       *pgxpool.Pool
	tenantID string
}

func NewFTNSiLKStore(db *pgxpool.Pool, tenantID string) *FTNSiLKStore {
	return &FTNSiLKStore{db: db, tenantID: strings.TrimSpace(tenantID)}
}

func (s *FTNSiLKStore) Insert(ctx context.Context, observedAt time.Time, r FlowRecord, serviceID, subscriberID, mainServerID string) error {
	if s == nil || s.db == nil { return errors.New("database_required") }
	if strings.TrimSpace(s.tenantID) == "" { return errors.New("tenant_required") }
	if strings.TrimSpace(r.ExporterID) == "" { return errors.New("exporter_required") }
	if r.Version != 5 && r.Version != 9 && r.Version != 10 { return errors.New("unsupported_flow_version") }
	if r.SamplingRate == 0 { r.SamplingRate = 1 }
	if r.Bytes > math.MaxInt64 || r.Packets > math.MaxInt64 { return errors.New("flow_counter_overflow") }

	var src, dst *string
	if p, err := netip.ParseAddr(strings.TrimSpace(r.SourceIP)); err == nil { v := p.String(); src = &v }
	if p, err := netip.ParseAddr(strings.TrimSpace(r.DestinationIP)); err == nil { v := p.String(); dst = &v }
	if strings.TrimSpace(r.SourceIP) != "" && src == nil { return errors.New("invalid_source_ip") }
	if strings.TrimSpace(r.DestinationIP) != "" && dst == nil { return errors.New("invalid_destination_ip") }

	if observedAt.IsZero() { observedAt = time.Now().UTC() }
	fingerprint := fmt.Sprintf("%016x", FlowRecordFingerprint(r))
	_, err := s.db.Exec(ctx, `
INSERT INTO ftn_silk_flow_records
(observed_at, tenant_id, exporter_id, version, source_ip, destination_ip,
 source_port, destination_port, protocol, bytes, packets, sampling_rate,
 fingerprint, service_id, subscriber_id, main_server_id)
VALUES ($1,$2,$3,$4,$5::inet,$6::inet,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
ON CONFLICT (tenant_id, exporter_id, fingerprint, observed_at) DO NOTHING`,
		observedAt, s.tenantID, r.ExporterID, r.Version, src, dst,
		r.SourcePort, r.DestinationPort, r.Protocol, int64(r.Bytes), int64(r.Packets),
		r.SamplingRate, fingerprint, nullableString(serviceID), nullableString(subscriberID), nullableString(mainServerID))
	return err
}

func nullableString(v string) *string {
	v = strings.TrimSpace(v)
	if v == "" { return nil }
	return &v
}
