CREATE TABLE IF NOT EXISTS isp_accounts (
  id UUID PRIMARY KEY,
  account_no TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  phone TEXT UNIQUE,
  email TEXT UNIQUE,
  status TEXT NOT NULL DEFAULT 'active',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS isp_services (
  id UUID PRIMARY KEY,
  account_id UUID NOT NULL REFERENCES isp_accounts(id) ON DELETE CASCADE,
  service_no TEXT UNIQUE NOT NULL,
  service_type TEXT NOT NULL,
  package_code TEXT,
  status TEXT NOT NULL DEFAULT 'active',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS isp_invoices (
  id UUID PRIMARY KEY,
  account_id UUID NOT NULL REFERENCES isp_accounts(id) ON DELETE RESTRICT,
  invoice_no TEXT UNIQUE NOT NULL,
  amount NUMERIC(14,2) NOT NULL CHECK (amount >= 0),
  currency CHAR(3) NOT NULL DEFAULT 'BDT',
  status TEXT NOT NULL DEFAULT 'unpaid',
  due_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS isp_payments (
  id UUID PRIMARY KEY,
  account_id UUID NOT NULL REFERENCES isp_accounts(id) ON DELETE RESTRICT,
  invoice_id UUID REFERENCES isp_invoices(id) ON DELETE SET NULL,
  provider TEXT NOT NULL,
  provider_reference TEXT,
  amount NUMERIC(14,2) NOT NULL CHECK (amount > 0),
  currency CHAR(3) NOT NULL DEFAULT 'BDT',
  status TEXT NOT NULL DEFAULT 'pending',
  idempotency_key TEXT UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS isp_support_tickets (
  id UUID PRIMARY KEY,
  account_id UUID NOT NULL REFERENCES isp_accounts(id) ON DELETE CASCADE,
  subject TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'open',
  priority TEXT NOT NULL DEFAULT 'normal',
  ai_classification TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS isp_services_account_idx ON isp_services(account_id);
CREATE INDEX IF NOT EXISTS isp_invoices_account_idx ON isp_invoices(account_id);
CREATE INDEX IF NOT EXISTS isp_payments_account_idx ON isp_payments(account_id);
CREATE INDEX IF NOT EXISTS isp_tickets_account_idx ON isp_support_tickets(account_id);
