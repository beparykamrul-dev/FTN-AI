-- FTN Billing Reconciliation + AI Review Layer
-- Database is the source of truth. AI is advisory and never mutates the ledger.

CREATE TABLE IF NOT EXISTS billing_ai_reviews (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  review_type text NOT NULL CHECK (review_type IN ('kyc','missing_entry','balance','collector','anomaly','relationship')),
  customer_id uuid,
  employee_id uuid,
  payment_id uuid REFERENCES billing_payments(id),
  invoice_id uuid REFERENCES billing_invoices(id),
  result text NOT NULL CHECK (result IN ('clear','needs_review','high_priority')),
  confidence numeric(6,5) CHECK (confidence >= 0 AND confidence <= 1),
  explanation text NOT NULL,
  evidence_refs jsonb NOT NULL DEFAULT '[]'::jsonb,
  model_ref text,
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_billing_ai_reviews_customer ON billing_ai_reviews(customer_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_billing_ai_reviews_employee ON billing_ai_reviews(employee_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_billing_ai_reviews_result ON billing_ai_reviews(result, created_at DESC);

CREATE OR REPLACE VIEW billing_customer_balance AS
SELECT
  i.customer_id,
  i.currency,
  SUM(i.amount_minor) FILTER (WHERE i.status IN ('issued','partially_paid','overdue')) AS invoiced_minor,
  COALESCE(SUM(p.amount_minor) FILTER (WHERE p.verification_status = 'verified'), 0) AS verified_paid_minor,
  GREATEST(
    SUM(i.amount_minor) FILTER (WHERE i.status IN ('issued','partially_paid','overdue'))
      - COALESCE(SUM(p.amount_minor) FILTER (WHERE p.verification_status = 'verified'), 0),
    0
  ) AS outstanding_minor
FROM billing_invoices i
LEFT JOIN billing_payments p ON p.customer_id = i.customer_id AND p.invoice_id = i.id
GROUP BY i.customer_id, i.currency;

CREATE OR REPLACE VIEW billing_employee_collection_summary AS
SELECT
  p.collected_by_employee_id AS employee_id,
  p.currency,
  COUNT(*) AS payment_count,
  SUM(p.amount_minor) AS collected_minor,
  SUM(p.amount_minor) FILTER (WHERE p.verification_status = 'verified') AS verified_collected_minor,
  SUM(p.amount_minor) FILTER (WHERE p.entered_at IS NULL) AS unentered_minor,
  COUNT(*) FILTER (WHERE p.entered_at IS NULL) AS missing_entry_count,
  COUNT(*) FILTER (WHERE p.entered_by IS DISTINCT FROM p.collected_by_employee_id AND p.entered_at IS NOT NULL) AS collector_entry_mismatch_count
FROM billing_payments p
GROUP BY p.collected_by_employee_id, p.currency;

CREATE OR REPLACE VIEW billing_missing_entry_review AS
SELECT
  p.id AS payment_id,
  p.customer_id,
  p.invoice_id,
  p.collected_by_employee_id,
  p.entered_by,
  p.amount_minor,
  p.currency,
  p.received_at,
  p.entered_at,
  CASE
    WHEN p.entered_at IS NOT NULL THEN 'entered'
    WHEN p.evidence ? 'customer_payment' THEN 'employee_not_entered'
    WHEN p.evidence ? 'employee_collection_task' THEN 'employee_not_entered'
    WHEN p.evidence ? 'queue_failure' OR p.evidence ? 'system_failure' THEN 'system_failure'
    WHEN p.evidence ? 'duplicate' OR p.evidence ? 'validation_error' THEN 'duplicate_or_invalid'
    ELSE 'unknown_requires_review'
  END AS probable_cause,
  p.evidence AS evidence
FROM billing_payments p
WHERE p.entered_at IS NULL;

CREATE OR REPLACE VIEW billing_employee_entry_gap_summary AS
SELECT
  collected_by_employee_id AS employee_id,
  currency,
  COUNT(*) FILTER (WHERE entered_at IS NULL) AS missing_entries,
  COALESCE(SUM(amount_minor) FILTER (WHERE entered_at IS NULL), 0) AS missing_entry_minor,
  COUNT(*) AS total_collections,
  COALESCE(SUM(amount_minor), 0) AS total_collection_minor
FROM billing_payments
GROUP BY collected_by_employee_id, currency;

CREATE OR REPLACE VIEW billing_reconciliation_dashboard AS
SELECT
  r.id AS case_id,
  r.customer_id,
  r.invoice_id,
  r.payment_id,
  r.case_type,
  r.reason,
  r.expected_amount_minor,
  r.observed_amount_minor,
  r.confidence,
  r.status,
  r.evidence,
  r.created_at,
  r.resolved_at
FROM billing_reconciliation_cases r;

-- Authorized legal/compliance review must use evidence-backed records only.
CREATE OR REPLACE VIEW billing_authorized_evidence AS
SELECT
  e.id,
  e.created_at,
  e.actor_type,
  e.actor_id,
  e.action,
  e.object_type,
  e.object_id,
  e.source_ip,
  e.request_id,
  e.before_state,
  e.after_state,
  e.evidence_refs,
  e.risk_score
FROM security_audit_events e
WHERE e.evidence_refs <> '[]'::jsonb;
