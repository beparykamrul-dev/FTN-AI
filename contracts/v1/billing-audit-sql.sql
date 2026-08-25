-- FTN Billing Audit / KYC evidence schema (PostgreSQL)
-- Isolated schema: does not overwrite existing billing tables.
CREATE SCHEMA IF NOT EXISTS ftn_audit;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS ftn_audit.audit_event (
    event_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type text NOT NULL,
    entity_type text NOT NULL,
    entity_id text NOT NULL,
    actor_type text NOT NULL,
    actor_id text,
    role_at_event text,
    occurred_at timestamptz NOT NULL DEFAULT now(),
    source_system text NOT NULL,
    session_id text,
    device_id text,
    ip_address inet,
    payload jsonb NOT NULL DEFAULT '{}'::jsonb,
    evidence_hash text,
    previous_event_hash text,
    created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS audit_event_entity_idx ON ftn_audit.audit_event(entity_type, entity_id, occurred_at DESC);
CREATE INDEX IF NOT EXISTS audit_event_actor_idx ON ftn_audit.audit_event(actor_id, occurred_at DESC);

CREATE TABLE IF NOT EXISTS ftn_audit.collection_observation (
    observation_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id text NOT NULL,
    invoice_id text,
    payment_id text,
    collector_id text,
    collector_role text,
    employment_status_at_event text,
    amount numeric(18,2) NOT NULL CHECK (amount >= 0),
    currency char(3) NOT NULL,
    method text NOT NULL,
    reference text,
    observed_at timestamptz NOT NULL,
    source text NOT NULL,
    session_id text,
    device_id text,
    evidence jsonb NOT NULL DEFAULT '{}'::jsonb,
    matched_ledger boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS collection_customer_idx ON ftn_audit.collection_observation(customer_id, observed_at DESC);
CREATE INDEX IF NOT EXISTS collection_collector_idx ON ftn_audit.collection_observation(collector_id, observed_at DESC);
CREATE INDEX IF NOT EXISTS collection_unmatched_idx ON ftn_audit.collection_observation(matched_ledger) WHERE matched_ledger = false;

CREATE TABLE IF NOT EXISTS ftn_audit.billing_gap (
    gap_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id text NOT NULL,
    billing_period_start date NOT NULL,
    billing_period_end date NOT NULL,
    expected_amount numeric(18,2) NOT NULL CHECK (expected_amount >= 0),
    invoiced_amount numeric(18,2) NOT NULL DEFAULT 0 CHECK (invoiced_amount >= 0),
    paid_amount numeric(18,2) NOT NULL DEFAULT 0 CHECK (paid_amount >= 0),
    entry_status text NOT NULL,
    cause text NOT NULL DEFAULT 'unknown_requires_review',
    evidence jsonb NOT NULL DEFAULT '{}'::jsonb,
    confidence numeric(5,4) CHECK (confidence >= 0 AND confidence <= 1),
    first_detected_at timestamptz NOT NULL DEFAULT now(),
    last_reviewed_at timestamptz,
    reviewed_by text,
    resolved_at timestamptz
);
CREATE UNIQUE INDEX IF NOT EXISTS billing_gap_period_uq ON ftn_audit.billing_gap(customer_id, billing_period_start, billing_period_end);

CREATE TABLE IF NOT EXISTS ftn_audit.kyc_observation (
    kyc_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id text NOT NULL,
    service_id text,
    verification_level text NOT NULL,
    status text NOT NULL,
    provider text,
    evidence_ref text,
    evidence_hash text,
    observed_at timestamptz NOT NULL DEFAULT now(),
    expires_at timestamptz,
    payload jsonb NOT NULL DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS kyc_customer_idx ON ftn_audit.kyc_observation(customer_id, observed_at DESC);

CREATE TABLE IF NOT EXISTS ftn_audit.ai_review (
    review_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    case_type text NOT NULL,
    entity_id text NOT NULL,
    risk_score numeric(5,4) CHECK (risk_score >= 0 AND risk_score <= 1),
    confidence numeric(5,4) CHECK (confidence >= 0 AND confidence <= 1),
    finding text NOT NULL,
    evidence_refs jsonb NOT NULL DEFAULT '[]'::jsonb,
    recommended_action text,
    human_decision text,
    decided_by text,
    decided_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now()
);

-- Safe integration views: replace source-table names in the integration layer,
-- not in this isolated audit schema.
CREATE OR REPLACE VIEW ftn_audit.collection_gap_queue AS
SELECT
    c.observation_id,
    c.customer_id,
    c.invoice_id,
    c.payment_id,
    c.collector_id,
    c.collector_role,
    c.employment_status_at_event,
    c.amount,
    c.currency,
    c.method,
    c.reference,
    c.observed_at,
    CASE
        WHEN c.matched_ledger THEN 'reconciled'
        WHEN c.payment_id IS NOT NULL THEN 'payment_observed_entry_unmatched'
        ELSE 'collection_observed_entry_unmatched'
    END AS status
FROM ftn_audit.collection_observation c
WHERE NOT c.matched_ledger;

CREATE OR REPLACE VIEW ftn_audit.employee_collection_summary AS
SELECT
    collector_id,
    collector_role,
    employment_status_at_event,
    currency,
    count(*) AS observation_count,
    sum(amount) AS observed_amount,
    count(*) FILTER (WHERE matched_ledger) AS reconciled_count,
    sum(amount) FILTER (WHERE matched_ledger) AS reconciled_amount,
    count(*) FILTER (WHERE NOT matched_ledger) AS gap_count,
    sum(amount) FILTER (WHERE NOT matched_ledger) AS gap_amount
FROM ftn_audit.collection_observation
GROUP BY collector_id, collector_role, employment_status_at_event, currency;

CREATE OR REPLACE VIEW ftn_audit.billing_gap_summary AS
SELECT
    cause,
    entry_status,
    count(*) AS case_count,
    sum(expected_amount) AS expected_amount,
    sum(invoiced_amount) AS invoiced_amount,
    sum(paid_amount) AS paid_amount,
    sum(expected_amount - paid_amount) AS outstanding_amount,
    avg(confidence) AS average_confidence
FROM ftn_audit.billing_gap
GROUP BY cause, entry_status;
