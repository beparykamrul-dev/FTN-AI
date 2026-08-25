CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS customer_kyc_verifications (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  customer_id uuid NOT NULL,
  verification_status text NOT NULL CHECK (verification_status IN ('pending','verified','rejected','manual_review')),
  verification_method text NOT NULL DEFAULT 'rule_engine',
  score numeric(6,5),
  reasons jsonb NOT NULL DEFAULT '[]'::jsonb,
  evidence_refs jsonb NOT NULL DEFAULT '[]'::jsonb,
  verified_at timestamptz,
  verified_by uuid,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_kyc_customer ON customer_kyc_verifications(customer_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_kyc_status ON customer_kyc_verifications(verification_status, created_at DESC);

CREATE TABLE IF NOT EXISTS billing_invoices (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  customer_id uuid NOT NULL,
  period_start date NOT NULL,
  period_end date NOT NULL,
  amount_minor bigint NOT NULL CHECK (amount_minor >= 0),
  currency char(3) NOT NULL,
  status text NOT NULL CHECK (status IN ('draft','issued','partially_paid','paid','overdue','void')),
  source_ref text,
  issued_at timestamptz,
  due_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CHECK (period_end >= period_start)
);
CREATE INDEX IF NOT EXISTS idx_invoice_customer_period ON billing_invoices(customer_id, period_start, period_end);
CREATE INDEX IF NOT EXISTS idx_invoice_status_due ON billing_invoices(status, due_at);

CREATE TABLE IF NOT EXISTS billing_payments (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  invoice_id uuid REFERENCES billing_invoices(id),
  customer_id uuid NOT NULL,
  collected_by_employee_id uuid,
  amount_minor bigint NOT NULL CHECK (amount_minor > 0),
  currency char(3) NOT NULL,
  payment_method text NOT NULL,
  external_ref text,
  received_at timestamptz NOT NULL DEFAULT now(),
  entered_at timestamptz,
  entered_by uuid,
  verification_status text NOT NULL DEFAULT 'pending' CHECK (verification_status IN ('pending','verified','rejected','manual_review')),
  evidence jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_payment_customer_date ON billing_payments(customer_id, received_at DESC);
CREATE INDEX IF NOT EXISTS idx_payment_collector_date ON billing_payments(collected_by_employee_id, received_at DESC);
CREATE INDEX IF NOT EXISTS idx_payment_entry ON billing_payments(entered_by, entered_at DESC);

CREATE TABLE IF NOT EXISTS billing_reconciliation_cases (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  customer_id uuid NOT NULL,
  invoice_id uuid REFERENCES billing_invoices(id),
  payment_id uuid REFERENCES billing_payments(id),
  case_type text NOT NULL CHECK (case_type IN ('missing_payment_entry','unmatched_payment','duplicate_entry','underpayment','overpayment','late_entry','collector_mismatch','employee_entry_missing')),
  reason text NOT NULL,
  expected_amount_minor bigint,
  observed_amount_minor bigint,
  confidence numeric(6,5) CHECK (confidence >= 0 AND confidence <= 1),
  status text NOT NULL DEFAULT 'open' CHECK (status IN ('open','reviewing','resolved','dismissed')),
  evidence jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamptz NOT NULL DEFAULT now(),
  resolved_at timestamptz,
  resolved_by uuid
);
CREATE INDEX IF NOT EXISTS idx_recon_open ON billing_reconciliation_cases(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_recon_customer ON billing_reconciliation_cases(customer_id, created_at DESC);

CREATE TABLE IF NOT EXISTS security_audit_events (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  actor_type text NOT NULL CHECK (actor_type IN ('user','employee','admin','system','ai','external')),
  actor_id uuid,
  action text NOT NULL,
  object_type text NOT NULL,
  object_id uuid,
  source_ip inet,
  user_agent text,
  request_id text,
  before_state jsonb,
  after_state jsonb,
  evidence_refs jsonb NOT NULL DEFAULT '[]'::jsonb,
  risk_score numeric(6,5),
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_security_audit_actor ON security_audit_events(actor_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_security_audit_object ON security_audit_events(object_type, object_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_security_audit_risk ON security_audit_events(risk_score DESC, created_at DESC);

CREATE OR REPLACE VIEW billing_collection_audit AS
SELECT
  p.id AS payment_id,
  p.customer_id,
  p.invoice_id,
  p.amount_minor,
  p.currency,
  p.payment_method,
  p.collected_by_employee_id,
  p.entered_by,
  p.received_at,
  p.entered_at,
  CASE
    WHEN p.entered_at IS NULL THEN 'employee_entry_missing'
    WHEN p.entered_by IS NULL THEN 'entry_actor_missing'
    WHEN p.collected_by_employee_id IS NULL THEN 'collector_missing'
    WHEN p.entered_by IS DISTINCT FROM p.collected_by_employee_id THEN 'collector_entry_mismatch'
    WHEN p.entered_at < p.received_at THEN 'invalid_entry_time'
    ELSE 'ok'
  END AS audit_status
FROM billing_payments p;

CREATE OR REPLACE VIEW billing_outstanding_summary AS
SELECT
  i.customer_id,
  SUM(i.amount_minor) FILTER (WHERE i.status IN ('issued','partially_paid','overdue')) AS issued_minor,
  COALESCE((SELECT SUM(p.amount_minor) FROM billing_payments p WHERE p.invoice_id = i.id AND p.verification_status = 'verified'),0) AS verified_payment_minor
FROM billing_invoices i
GROUP BY i.customer_id;
