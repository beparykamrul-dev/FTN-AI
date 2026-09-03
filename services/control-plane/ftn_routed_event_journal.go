package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// FTNRoutedPostgresEventSink persists normalized routing events through the
// durable event-journal primitive. It stores no raw protocol payloads and does
// not execute routing changes.
type FTNRoutedPostgresEventSink struct {
	db       *pgxpool.Pool
	tenantID string
}

func NewFTNRoutedPostgresEventSink(db *pgxpool.Pool, tenantID string) *FTNRoutedPostgresEventSink {
	return &FTNRoutedPostgresEventSink{db: db, tenantID: strings.TrimSpace(tenantID)}
}

func (s *FTNRoutedPostgresEventSink) Publish(ctx context.Context, event FTNBGPEvent) error {
	if s == nil || s.db == nil {
		return errors.New("database_required")
	}
	if strings.TrimSpace(s.tenantID) == "" {
		return errors.New("tenant_required")
	}
	if strings.TrimSpace(event.Type) == "" {
		return errors.New("event_type_required")
	}

	payload, err := json.Marshal(struct {
		Type      string `json:"type"`
		Peer      string `json:"peer,omitempty"`
		State     string `json:"state,omitempty"`
		Prefix    string `json:"prefix,omitempty"`
		Timestamp string `json:"timestamp,omitempty"`
	}{
		Type: event.Type, Peer: event.Peer, State: event.State,
		Prefix: event.Prefix, Timestamp: event.Timestamp,
	})
	if err != nil {
		return err
	}

	correlationID := routeEventCorrelationID(event)
	aggregateID := strings.TrimSpace(event.Peer)
	if aggregateID == "" {
		aggregateID = strings.TrimSpace(event.Prefix)
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = appendEventTx(
		tx,
		ctx,
		s.tenantID,
		strings.TrimSpace(event.Type),
		correlationID,
		correlationID,
		aggregateID,
		json.RawMessage(payload),
	)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func routeEventCorrelationID(event FTNBGPEvent) string {
	payload := strings.Join([]string{
		strings.TrimSpace(event.Type),
		strings.TrimSpace(event.Peer),
		strings.TrimSpace(event.State),
		strings.TrimSpace(event.Prefix),
		strings.TrimSpace(event.Timestamp),
	}, "|")
	sum := sha256.Sum256([]byte(payload))
	return "ftn-routed-" + hex.EncodeToString(sum[:])
}
