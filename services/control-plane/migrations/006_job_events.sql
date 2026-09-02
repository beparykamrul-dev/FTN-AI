-- FTN job lifecycle event integrity. Additive/idempotent.
CREATE INDEX IF NOT EXISTS event_journal_aggregate_sequence_idx
  ON event_journal(tenant_id, aggregate_id, sequence);

-- Prevent duplicate lifecycle events for the same aggregate/action/attempt marker.
CREATE UNIQUE INDEX IF NOT EXISTS event_journal_job_lifecycle_dedupe_idx
  ON event_journal(tenant_id, aggregate_id, event_type, causation_id)
  WHERE aggregate_id <> '' AND causation_id <> '';
