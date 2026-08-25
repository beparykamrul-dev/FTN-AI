-- FTN Billing Audit hardening layer.
-- Safe to apply after contracts/v1/billing-audit-sql.sql.
-- Does not alter or guess the existing billing/ledger source schema.

CREATE SCHEMA IF NOT EXISTS ftn_audit;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Append-only audit events: corrections are new events, never UPDATE/DELETE.
CREATE OR REPLACE FUNCTION ftn_audit.reject_audit_mutation()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
  RAISE EXCEPTION 'ftn_audit.audit_event is append-only; create a compensating audit event instead';
END;
$$;

DROP TRIGGER IF EXISTS audit_event_immutable ON ftn_audit.audit_event;
CREATE TRIGGER audit_event_immutable
BEFORE UPDATE OR DELETE ON ftn_audit.audit_event
FOR EACH ROW EXECUTE FUNCTION ftn_audit.reject_audit_mutation();

-- Immutable collection observations preserve the original collection evidence.
CREATE OR REPLACE FUNCTION ftn_audit.reject_collection_mutation()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
  RAISE EXCEPTION 'ftn_audit.collection_observation is immutable; reconcile through audit events';
END;
$$;

DROP TRIGGER IF EXISTS collection_observation_immutable ON ftn_audit.collection_observation;
CREATE TRIGGER collection_observation_immutable
BEFORE UPDATE OR DELETE ON ftn_audit.collection_observation
FOR EACH ROW EXECUTE FUNCTION ftn_audit.reject_collection_mutation();

-- Evidence-aware collection observation entry point for the integration layer.
CREATE OR REPLACE FUNCTION ftn_audit.record_collection_observation(
  p_customer_id text,
  p_invoice_id text,
  p_payment_id text,
  p_collector_id text,
  p_collector_role text,
  p_employment_status_at_event text,
  p_amount numeric,
  p_currency char(3),
  p_method text,
  p_reference text,
  p_observed_at timestamptz,
  p_source text,
  p_session_id text,
  p_device_id text,
  p_evidence jsonb
)
RETURNS uuid
LANGUAGE plpgsql
AS $$
DECLARE
  v_id uuid;
BEGIN
  IF p_customer_id IS NULL OR btrim(p_customer_id) = '' THEN
    RAISE EXCEPTION 'customer_id is required';
  END IF;
  IF p_amount IS NULL OR p_amount < 0 THEN
    RAISE EXCEPTION 'amount must be non-negative';
  END IF;
  IF p_currency IS NULL OR length(btrim(p_currency::text)) <> 3 THEN
    RAISE EXCEPTION '3-letter currency is required';
  END IF;
  IF p_source IS NULL OR btrim(p_source) = '' THEN
    RAISE EXCEPTION 'source is required';
  END IF;
  IF p_observed_at IS NULL THEN
    RAISE EXCEPTION 'observed_at is required';
  END IF;

  INSERT INTO ftn_audit.collection_observation(
    customer_id, invoice_id, payment_id, collector_id, collector_role,
    employment_status_at_event, amount, currency, method, reference,
    observed_at, source, session_id, device_id, evidence
  ) VALUES (
    p_customer_id, p_invoice_id, p_payment_id, p_collector_id, p_collector_role,
    p_employment_status_at_event, p_amount, upper(p_currency), p_method, p_reference,
    p_observed_at, p_source, p_session_id, p_device_id, coalesce(p_evidence, '{}'::jsonb)
  ) RETURNING observation_id INTO v_id;

  RETURN v_id;
END;
$$;

-- Deterministic billing-gap classification. AI may explain/prioritize, but this
-- function only accepts evidence-backed causes from the integration layer.
CREATE OR REPLACE FUNCTION ftn_audit.classify_billing_gap(
  p_has_verified_payment boolean,
  p_has_collection_evidence boolean,
  p_has_matching_ledger_entry boolean,
  p_has_system_failure boolean,
  p_has_duplicate_or_validation_evidence boolean
)
RETURNS text
LANGUAGE sql
IMMUTABLE
AS $$
  SELECT CASE
    WHEN p_has_duplicate_or_validation_evidence THEN 'duplicate_or_invalid'
    WHEN p_has_system_failure THEN 'system_failure'
    WHEN p_has_collection_evidence AND NOT p_has_matching_ledger_entry THEN 'employee_not_entered'
    WHEN NOT p_has_verified_payment THEN 'customer_not_paid'
    ELSE 'unknown_requires_review'
  END;
$$;

-- Useful for the control panel: employee collection totals and unreconciled exposure.
CREATE OR REPLACE VIEW ftn_audit.employee_accountability AS
SELECT
  collector_id AS employee_id,
  collector_role,
  employment_status_at_event,
  currency,
  count(*) AS observation_count,
  sum(amount) AS observed_amount,
  count(*) FILTER (WHERE matched_ledger) AS reconciled_count,
  coalesce(sum(amount) FILTER (WHERE matched_ledger), 0) AS reconciled_amount,
  count(*) FILTER (WHERE NOT matched_ledger) AS unreconciled_count,
  coalesce(sum(amount) FILTER (WHERE NOT matched_ledger), 0) AS unreconciled_amount,
  min(observed_at) AS first_observation_at,
  max(observed_at) AS last_observation_at
FROM ftn_audit.collection_observation
GROUP BY collector_id, collector_role, employment_status_at_event, currency;

-- Customer-level exposure view for the audit dashboard.
CREATE OR REPLACE VIEW ftn_audit.customer_exposure AS
SELECT
  customer_id,
  billing_period_start,
  billing_period_end,
  expected_amount,
  invoiced_amount,
  paid_amount,
  greatest(expected_amount - paid_amount, 0) AS outstanding_amount,
  entry_status,
  cause,
  confidence,
  evidence,
  first_detected_at,
  last_reviewed_at,
  resolved_at
FROM ftn_audit.billing_gap;

-- Investigation cases are evidence containers, not criminal-affiliation labels.
CREATE TABLE IF NOT EXISTS ftn_audit.investigation_case (
  case_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  case_type text NOT NULL,
  subject_type text NOT NULL,
  subject_id text NOT NULL,
  status text NOT NULL DEFAULT 'open',
  severity text NOT NULL DEFAULT 'review',
  deterministic_reason text NOT NULL,
  evidence_refs jsonb NOT NULL DEFAULT '[]'::jsonb,
  source_event_ids jsonb NOT NULL DEFAULT '[]'::jsonb,
  created_by text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  closed_by text,
  closed_at timestamptz
);
CREATE INDEX IF NOT EXISTS investigation_subject_idx
  ON ftn_audit.investigation_case(subject_type, subject_id, created_at DESC);
CREATE INDEX IF NOT EXISTS investigation_status_idx
  ON ftn_audit.investigation_case(status, severity, created_at DESC);

-- Explicit evidence links make later authorized/legal review traceable.
CREATE TABLE IF NOT EXISTS ftn_audit.investigation_evidence (
  evidence_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  case_id uuid NOT NULL REFERENCES ftn_audit.investigation_case(case_id),
  evidence_type text NOT NULL,
  source_system text NOT NULL,
  source_id text NOT NULL,
  collected_at timestamptz NOT NULL DEFAULT now(),
  collected_by text,
  content_hash text,
  metadata jsonb NOT NULL DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS investigation_evidence_case_idx
  ON ftn_audit.investigation_evidence(case_id, collected_at);

-- No AI/legal conclusion is stored here. Human decisions remain explicit and auditable.
CREATE OR REPLACE VIEW ftn_audit.authorized_case_evidence AS
SELECT
  c.case_id,
  c.case_type,
  c.subject_type,
  c.subject_id,
  c.status,
  c.severity,
  c.deterministic_reason,
  e.evidence_id,
  e.evidence_type,
  e.source_system,
  e.source_id,
  e.collected_at,
  e.collected_by,
  e.content_hash,
  e.metadata
FROM ftn_audit.investigation_case c
JOIN ftn_audit.investigation_evidence e ON e.case_id = c.case_id;

-- Performance indexes for the main dashboard queries.
CREATE INDEX IF NOT EXISTS billing_gap_customer_status_idx
  ON ftn_audit.billing_gap(customer_id, entry_status, billing_period_end DESC);
CREATE INDEX IF NOT EXISTS billing_gap_cause_idx
  ON ftn_audit.billing_gap(cause, first_detected_at DESC);
CREATE INDEX IF NOT EXISTS audit_event_time_idx
  ON ftn_audit.audit_event(occurred_at DESC);
CREATE INDEX IF NOT EXISTS kyc_status_idx
  ON ftn_audit.kyc_observation(status, observed_at DESC);
