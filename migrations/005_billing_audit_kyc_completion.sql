-- FTN Billing Audit + KYC Completion Layer
-- Database is authoritative. AI is advisory only.
-- This migration adds customer/employee reconciliation surfaces and
-- evidence-backed relationship signals without inferring criminal conduct.

CREATE TABLE IF NOT EXISTS billing_kyc_ai_reviews (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  customer_id uuid NOT NULL,
  verification_status text NOT NULL CHECK (verification_status IN ('clear','needs_review','manual_review')),
  confidence numeric(6,5) CHECK (confidence >= 0 AND confidence <= 1),
  reasons jsonb NOT NULL DEFAULT '[]'::jsonb,
  evidence_refs jsonb NOT NULL DEFAULT '[]'::jsonb,
  model_ref text,
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_billing_kyc_ai_customer ON billing_kyc_ai_reviews(customer_id, created_at DESC);

CREATE OR REPLACE VIEW billing_customer_accountability AS
SELECT
  b.customer_id,
  b.currency,
  b.billed_minor AS invoiced_minor,
  b.verified_paid_minor,
  GREATEST(b.outstanding_minor, 0) AS outstanding_minor,
  b.audit_status AS balance_status,
  k.verification_status AS latest_kyc_status,
  k.score AS latest_kyc_score,
  a.verification_status AS latest_ai_kyc_status,
  a.confidence AS latest_ai_kyc_confidence,
  a.reasons AS latest_ai_kyc_reasons,
  a.evidence_refs AS latest_ai_kyc_evidence_refs
FROM billing_customer_balance_audit b
LEFT JOIN LATERAL (
  SELECT verification_status, score
  FROM customer_kyc_verifications k
  WHERE k.customer_id = b.customer_id
  ORDER BY k.created_at DESC
  LIMIT 1
) k ON true
LEFT JOIN LATERAL (
  SELECT verification_status, confidence, reasons, evidence_refs
  FROM billing_kyc_ai_reviews a
  WHERE a.customer_id = b.customer_id
  ORDER BY a.created_at DESC
  LIMIT 1
) a ON true;

CREATE OR REPLACE VIEW billing_payment_entry_accountability AS
SELECT
  p.id AS payment_id,
  p.customer_id,
  p.invoice_id,
  p.collected_by_employee_id AS collector_employee_id,
  p.entered_by AS entry_employee_id,
  p.amount_minor,
  p.currency,
  p.payment_method,
  p.external_ref,
  p.received_at,
  p.entered_at,
  p.verification_status,
  CASE
    WHEN p.entered_at IS NOT NULL AND p.entered_by IS DISTINCT FROM p.collected_by_employee_id
      THEN 'collector_entry_mismatch'
    WHEN p.entered_at IS NOT NULL THEN 'entered'
    WHEN p.evidence ? 'customer_payment' THEN 'customer_payment_not_entered'
    WHEN p.evidence ? 'employee_collection_task' THEN 'employee_collection_not_entered'
    WHEN p.evidence ? 'queue_failure' OR p.evidence ? 'system_failure' THEN 'system_failure'
    ELSE 'cause_unknown_requires_review'
  END AS entry_gap_reason,
  p.evidence AS payment_evidence
FROM billing_payments p;

CREATE OR REPLACE VIEW billing_employee_money_accountability AS
SELECT
  p.collected_by_employee_id AS employee_id,
  p.currency,
  COUNT(*) AS collection_count,
  COALESCE(SUM(p.amount_minor), 0) AS collected_minor,
  COALESCE(SUM(p.amount_minor) FILTER (WHERE p.verification_status = 'verified'), 0) AS verified_collected_minor,
  COALESCE(SUM(p.amount_minor) FILTER (WHERE p.entered_at IS NULL), 0) AS unentered_minor,
  COUNT(*) FILTER (WHERE p.entered_at IS NULL) AS unentered_count,
  COUNT(*) FILTER (
    WHERE p.entered_at IS NOT NULL
      AND p.entered_by IS DISTINCT FROM p.collected_by_employee_id
  ) AS collector_entry_mismatch_count
FROM billing_payments p
WHERE p.collected_by_employee_id IS NOT NULL
GROUP BY p.collected_by_employee_id, p.currency;

CREATE OR REPLACE VIEW billing_customer_payment_history AS
SELECT
  p.customer_id,
  p.currency,
  COUNT(*) AS payment_count,
  COALESCE(SUM(p.amount_minor), 0) AS payment_minor,
  COUNT(*) FILTER (WHERE p.verification_status = 'verified') AS verified_payment_count,
  COALESCE(SUM(p.amount_minor) FILTER (WHERE p.verification_status = 'verified'), 0) AS verified_payment_minor,
  COUNT(*) FILTER (WHERE p.entered_at IS NULL) AS missing_entry_count,
  COALESCE(SUM(p.amount_minor) FILTER (WHERE p.entered_at IS NULL), 0) AS missing_entry_minor,
  MAX(p.received_at) AS last_payment_at
FROM billing_payments p
GROUP BY p.customer_id, p.currency;

CREATE OR REPLACE VIEW billing_relationship_review_signals AS
SELECT
  e.actor_id AS employee_id,
  e.object_id,
  COUNT(*) AS event_count,
  MIN(e.created_at) AS first_seen_at,
  MAX(e.created_at) AS last_seen_at,
  MAX(e.risk_score) AS max_risk_score,
  jsonb_agg(DISTINCT e.action) AS actions,
  jsonb_agg(DISTINCT e.object_type) AS object_types,
  jsonb_agg(DISTINCT e.evidence_refs) AS evidence_refs
FROM security_audit_events e
WHERE e.actor_type IN ('employee','admin','system','ai')
  AND e.evidence_refs <> '[]'::jsonb
GROUP BY e.actor_id, e.object_id;

CREATE OR REPLACE VIEW billing_audit_case_package AS
SELECT
  f.finding_id,
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
  f.evidence_chain
FROM billing_authorized_case_export f;

-- AI may read these views and write review records; it must not mutate invoices,
-- payments, KYC records, or the authoritative billing ledger.
