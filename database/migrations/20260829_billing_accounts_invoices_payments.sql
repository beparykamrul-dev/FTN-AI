-- FTN billing foundation: service billing, consolidated home billing, payments and audit.
-- Provider credentials/secrets are intentionally not stored here.

CREATE TABLE IF NOT EXISTS billing_accounts (
    id BIGSERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    account_type TEXT NOT NULL CHECK (account_type IN ('personal','family','business','organization')),
    currency CHAR(3) NOT NULL DEFAULT 'BDT',
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_billing_accounts_user_id ON billing_accounts(user_id);

CREATE TABLE IF NOT EXISTS invoices (
    id BIGSERIAL PRIMARY KEY,
    billing_account_id BIGINT NOT NULL REFERENCES billing_accounts(id),
    invoice_number TEXT NOT NULL UNIQUE,
    service_scope TEXT NOT NULL DEFAULT 'consolidated',
    status TEXT NOT NULL DEFAULT 'issued' CHECK (status IN ('draft','issued','due','paid','overdue','cancelled','refunded')),
    subtotal_minor BIGINT NOT NULL DEFAULT 0,
    total_minor BIGINT NOT NULL DEFAULT 0,
    currency CHAR(3) NOT NULL DEFAULT 'BDT',
    issued_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    due_at TIMESTAMPTZ,
    paid_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_invoices_account_status ON invoices(billing_account_id,status);
CREATE INDEX IF NOT EXISTS idx_invoices_due_at ON invoices(due_at);

CREATE TABLE IF NOT EXISTS invoice_items (
    id BIGSERIAL PRIMARY KEY,
    invoice_id BIGINT NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    service_id TEXT NOT NULL,
    description TEXT NOT NULL,
    quantity BIGINT NOT NULL DEFAULT 1,
    unit_price_minor BIGINT NOT NULL DEFAULT 0,
    total_minor BIGINT NOT NULL DEFAULT 0,
    currency CHAR(3) NOT NULL DEFAULT 'BDT'
);

CREATE INDEX IF NOT EXISTS idx_invoice_items_invoice ON invoice_items(invoice_id);

CREATE TABLE IF NOT EXISTS payment_transactions (
    id BIGSERIAL PRIMARY KEY,
    billing_account_id BIGINT NOT NULL REFERENCES billing_accounts(id),
    invoice_id BIGINT REFERENCES invoices(id),
    provider TEXT NOT NULL,
    provider_transaction_id TEXT,
    idempotency_key TEXT NOT NULL UNIQUE,
    amount_minor BIGINT NOT NULL,
    currency CHAR(3) NOT NULL DEFAULT 'BDT',
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','verified','failed','refunded')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    verified_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_payment_provider_tx
    ON payment_transactions(provider,provider_transaction_id)
    WHERE provider_transaction_id IS NOT NULL;

CREATE TABLE IF NOT EXISTS billing_audit_events (
    id BIGSERIAL PRIMARY KEY,
    billing_account_id BIGINT REFERENCES billing_accounts(id),
    actor_id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    event_data JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_billing_audit_account_time
    ON billing_audit_events(billing_account_id,created_at DESC);
