CREATE TABLE IF NOT EXISTS billing_audit_runs (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  run_type text NOT NULL CHECK (run_type IN ('scheduled','manual','incident')),
  period_start date,
  period_end date,
  status text NOT NULL DEFAULT 'running' CHECK (status IN ('running','completed','failed')),
  started_at timestamptz NOT NULL DEFAULT now(),
  completed_at timestamptz,
  summary jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_by uuid
);

CREATE TABLE IF NOT EXISTS billing_audit_findings (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  audit_run_id uuid REFERENCES billing_audit_runs(id) ON DELETE CASCADE,
  customer_id uuid,
  employee_id uuid,
  invoice_id uuid REFERENCES billing_invoices(id),
  payment_id uuid REFERENCES billing_payments(id),
  finding_type text NOT NULL CHECK (finding_type IN (
    'outstanding_balance','missing_payment_entry','unmatched_payment',
    'employee_entry_missing','collector_entry_mismatch','duplicate_payment',
    'underpayment','overpayment','late_entry','kyc_mismatch','unusual_collection'
  )),
  expected_amount_minor bigint,
  observed_amount_minor bigint,
  confidence numeric(6,5) CHECK (confidence >= 0 AND confidence <= 1),
  severity text NOT NULL DEFAULT 'info' CHECK (severity IN ('info','low','medium','high','critical')),
  status text NOT NULL DEFAULT 'open' CHECK (status IN ('open','reviewing','confirmed','resolved','dismissed')),
  explanation text NOT NULL,
  evidence_refs jsonb NOT NULL DEFAULT '[]'::jsonb,
  ai_analysis jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamptz NOT NULL DEFAULT now(),
  reviewed_at timestamptz,
  reviewed_by uuid
);

CREATE INDEX IF NOT EXISTS idx_billing_audit_findings_run ON billing_audit_findings(audit_run_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_billing_audit_findings_customer ON billing_audit_findings(customer_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_billing_audit_findings_employee ON billing_audit_findings(employee_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_billing_audit_findings_status ON billing_audit_findings(status, severity, created_at DESC);

CREATE OR REPLACE VIEW billing_customer_balance_audit AS
WITH invoice_totals AS (
  SELECT customer_id,
         SUM(amount_minor) FILTER (WHERE status IN ('issued','partially_paid','overdue')) AS billed_minor
  FROM billing_invoices
  GROUP BY customer_id
), payment_totals AS (
  SELECT customer_id,
         SUM(amount_minor) FILTER (WHERE verification_status = 'verified') AS verified_paid_minor,
         SUM(amount_minor) AS recorded_paid_minor
  FROM billing_payments
  GROUP BY customer_id
)
SELECT
  COALESCE(i.customer_id, p.customer_id) AS customer_id,
  COALESCE(i.billed_minor,0) AS billed_minor,
  COALESCE(p.verified_paid_minor,0) AS verified_paid_minor,
  COALESCE(p.recorded_paid_minor,0) AS recorded_paid_minor,
  GREATEST(COALESCE(i.billed_minor,0) - COALESCE(p.verified_paid_minor,0),0) AS outstanding_minor,
  CASE
    WHEN COALESCE(i.billed_minor,0) <= COALESCE(p.verified_paid_minor,0) THEN 'clear'
    WHEN COALESCE(p.recorded_paid_minor,0) > COALESCE(p.verified_paid_minor,0) THEN 'payment_verification_pending'
    ELSE 'outstanding'
  END AS audit_status
FROM invoice_totals i
FULL OUTER JOIN payment_totals p ON p.customer_id = i.customer_id;

CREATE OR REPLACE VIEW billing_employee_collection_audit AS
SELECT
  p.collected_by_employee_id AS employee_id,
  COUNT(*) AS payment_count,
  COALESCE(SUM(p.amount_minor),0) AS recorded_collection_minor,
  COALESCE(SUM(p.amount_minor) FILTER (WHERE p.verification_status='verified'),0) AS verified_collection_minor,
  COALESCE(SUM(p.amount_minor) FILTER (WHERE p.entered_at IS NULL),0) AS unentered_collection_minor,
  COUNT(*) FILTER (WHERE p.entered_at IS NULL) AS missing_entry_count,
  COUNT(*) FILTER (WHERE p.entered_by IS DISTINCT FROM p.collected_by_employee_id) AS collector_entry_mismatch_count
FROM billing_payments p
WHERE p.collected_by_employee_id IS NOT NULL
GROUP BY p.collected_by_employee_id;

CREATE OR REPLACE VIEW billing_missing_entry_queue AS
SELECT
  p.id AS payment_id,
  p.customer_id,
  p.invoice_id,
  p.collected_by_employee_id AS employee_id,
  p.amount_minor,
  p.currency,
  p.payment_method,
  p.external_ref,
  p.received_at,
  p.entered_at,
  CASE
    WHEN p.entered_at IS NULL AND p.entered_by IS NULL THEN 'payment_recorded_but_employee_entry_missing'
    WHEN p.entered_at IS NULL THEN 'employee_entry_missing'
    WHEN p.entered_by IS DISTINCT FROM p.collected_by_employee_id THEN 'collector_and_entry_actor_differ'
    ELSE 'ok'
  END AS reason
FROM billing_payments p
WHERE p.entered_at IS NULL
   OR p.entered_by IS DISTINCT FROM p.collected_by_employee_id;

CREATE OR REPLACE VIEW billing_customer_kyc_audit AS
SELECT DISTINCT ON (k.customer_id)
  k.customer_id,
  k.verification_status,
  k.verification_method,
  k.score,
  k.reasons,
  k.evidence_refs,
  k.verified_at,
  k.updated_at
FROM customer_kyc_verifications k
ORDER BY k.customer_id, k.created_at DESC;

-- Evidence-first relationship view. This does not label anyone as part of a criminal
-- group; it only exposes auditable relationships that can be reviewed by authorized staff.
CREATE OR REPLACE VIEW billing_actor_relationship_audit AS
SELECT
  p.customer_id,
  p.collected_by_employee_id AS employee_id,
  COUNT(*) AS payment_count,
  SUM(p.amount_minor) AS amount_minor,
  MIN(p.received_at) AS first_seen_at,
  MAX(p.received_at) AS last_seen_at,
  COUNT(*) FILTER (WHERE p.entered_by IS NULL) AS missing_entry_count,
  COUNT(*) FILTER (WHERE p.entered_by IS DISTINCT FROM p.collected_by_employee_id) AS entry_mismatch_count
FROM billing_payments p
WHERE p.collected_by_employee_id IS NOT NULL
GROUP BY p.customer_id, p.collected_by_employee_id;

-- Deterministic audit findings; AI may enrich explanations but must not create unsupported accusations.
CREATE OR REPLACE FUNCTION generate_billing_audit_findings(p_run_id uuid)
RETURNS void
LANGUAGE plpgsql
AS $$
BEGIN
  INSERT INTO billing_audit_findings (
    audit_run_id, customer_id, invoice_id, payment_id, employee_id,
    finding_type, expected_amount_minor, observed_amount_minor,
    confidence, severity, explanation, evidence_refs
  )
  SELECT
    p_run_id, customer_id, NULL, NULL, NULL,
    'outstanding_balance', billed_minor, verified_paid_minor,
    1.0,
    CASE WHEN outstanding_minor > 0 THEN 'medium' ELSE 'info' END,
    'Database reconciliation shows an outstanding verified balance.',
    jsonb_build_array(jsonb_build_object('view','billing_customer_balance_audit','customer_id',customer_id))
  FROM billing_customer_balance_audit
  WHERE outstanding_minor > 0;

  INSERT INTO billing_audit_findings (
    audit_run_id, customer_id, invoice_id, payment_id, employee_id,
    finding_type, observed_amount_minor, confidence, severity, explanation, evidence_refs
  )
  SELECT
    p_run_id, customer_id, invoice_id, payment_id, employee_id,
    'employee_entry_missing', amount_minor, 1.0, 'high',
    reason,
    jsonb_build_array(jsonb_build_object('view','billing_missing_entry_queue','payment_id',payment_id))
  FROM billing_missing_entry_queue
  WHERE payment_id IS NOT NULL;
END;
$$;
