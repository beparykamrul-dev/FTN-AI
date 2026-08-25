-- FTN Billing Audit Evidence + Authorized Review Layer
-- Database remains the source of truth. This migration adds evidence-chain
-- integrity, role-at-time snapshots, and evidence-backed review exports.

CREATE TABLE IF NOT EXISTS billing_audit_evidence (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  audit_run_id uuid REFERENCES billing_audit_runs(id) ON DELETE CASCADE,
  finding_id uuid REFERENCES billing_audit_findings(id) ON DELETE CASCADE,
  source_type text NOT NULL CHECK (source_type IN (
    'invoice','payment','gateway','employee_session','device_session',
    'api_event','audit_event','kyc','system_event','public_record'
  )),
  source_id text NOT NULL,
  captured_at timestamptz NOT NULL DEFAULT now(),
  evidence jsonb NOT NULL DEFAULT '{}'::jsonb,
  evidence_sha256 text,
  provenance jsonb NOT NULL DEFAULT '{}'::jsonb,
  legal_hold boolean NOT NULL DEFAULT false,
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_billing_audit_evidence_finding ON billing_audit_evidence(finding_id, captured_at DESC);
CREATE INDEX IF NOT EXISTS idx_billing_audit_evidence_source ON billing_audit_evidence(source_type, source_id);
CREATE INDEX IF NOT EXISTS idx_billing_audit_evidence_hold ON billing_audit_evidence(legal_hold, created_at DESC);

CREATE TABLE IF NOT EXISTS billing_actor_role_history (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  actor_id uuid NOT NULL,
  role_name text NOT NULL,
  status text NOT NULL CHECK (status IN ('active','inactive','suspended','revoked','reauthorized')),
  valid_from timestamptz NOT NULL,
  valid_to timestamptz,
  source_audit_event_id uuid REFERENCES security_audit_events(id),
  created_at timestamptz NOT NULL DEFAULT now(),
  CHECK (valid_to IS NULL OR valid_to >= valid_from)
);
CREATE INDEX IF NOT EXISTS idx_actor_role_history_actor ON billing_actor_role_history(actor_id, valid_from DESC);

CREATE OR REPLACE VIEW billing_employee_role_at_collection AS
SELECT
  p.id AS payment_id,
  p.customer_id,
  p.collected_by_employee_id AS employee_id,
  p.amount_minor,
  p.currency,
  p.received_at,
  COALESCE(
    h.role_name,
    'role_not_recorded'
  ) AS role_at_collection,
  COALESCE(h.status, 'status_not_recorded') AS status_at_collection
FROM billing_payments p
LEFT JOIN LATERAL (
  SELECT role_name, status
  FROM billing_actor_role_history h
  WHERE h.actor_id = p.collected_by_employee_id
    AND h.valid_from <= p.received_at
    AND (h.valid_to IS NULL OR h.valid_to >= p.received_at)
  ORDER BY h.valid_from DESC
  LIMIT 1
) h ON true;

CREATE OR REPLACE VIEW billing_evidence_chain AS
SELECT
  e.id AS evidence_id,
  e.audit_run_id,
  e.finding_id,
  e.source_type,
  e.source_id,
  e.captured_at,
  e.evidence_sha256,
  e.provenance,
  e.legal_hold,
  e.created_at
FROM billing_audit_evidence e;

-- Authorized legal/compliance export surface. It exposes only evidence-backed
-- findings and does not infer criminal affiliation.
CREATE OR REPLACE VIEW billing_authorized_case_export AS
SELECT
  f.id AS finding_id,
  f.audit_run_id,
  f.customer_id,
  f.employee_id,
  f.invoice_id,
  f.payment_id,
  f.finding_type,
  f.expected_amount_minor,
  f.observed_amount_minor,
  f.confidence,
  f.severity,
  f.status,
  f.explanation,
  f.evidence_refs,
  f.ai_analysis,
  f.created_at,
  f.reviewed_at,
  f.reviewed_by,
  COALESCE(
    jsonb_agg(
      jsonb_build_object(
        'evidence_id', e.id,
        'source_type', e.source_type,
        'source_id', e.source_id,
        'captured_at', e.captured_at,
        'evidence_sha256', e.evidence_sha256,
        'provenance', e.provenance,
        'legal_hold', e.legal_hold
      ) ORDER BY e.captured_at
    ) FILTER (WHERE e.id IS NOT NULL),
    '[]'::jsonb
  ) AS evidence_chain
FROM billing_audit_findings f
LEFT JOIN billing_audit_evidence e ON e.finding_id = f.id
GROUP BY f.id;

-- Database-backed dashboard summary for management.
CREATE OR REPLACE VIEW billing_audit_management_summary AS
SELECT
  COUNT(*) FILTER (WHERE audit_status = 'outstanding') AS customers_with_outstanding,
  COALESCE(SUM(outstanding_minor),0) AS total_outstanding_minor,
  COUNT(*) FILTER (WHERE audit_status = 'payment_verification_pending') AS payment_verification_pending_customers
FROM billing_customer_balance_audit;

CREATE OR REPLACE VIEW billing_employee_collection_gap AS
SELECT
  employee_id,
  payment_count,
  recorded_collection_minor,
  verified_collection_minor,
  unentered_collection_minor,
  missing_entry_count,
  collector_entry_mismatch_count,
  CASE
    WHEN recorded_collection_minor = 0 THEN 0
    ELSE ROUND((unentered_collection_minor::numeric / recorded_collection_minor::numeric) * 100, 2)
  END AS unentered_percent
FROM billing_employee_collection_audit;

-- AI consumes these views and writes only to billing_ai_reviews / billing_audit_findings.
-- It must never mutate invoices, payments, KYC records, or the authoritative ledger.
