CREATE TABLE IF NOT EXISTS accounts (
 id uuid PRIMARY KEY,
 name text NOT NULL,
 type text NOT NULL CHECK (type IN ('asset','liability','equity','income','expense')),
 currency char(3) NOT NULL,
 active boolean NOT NULL DEFAULT true,
 created_at timestamptz NOT NULL DEFAULT now(),
 updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS journal_entries (
 id uuid PRIMARY KEY,
 journal_id uuid NOT NULL,
 account_id uuid NOT NULL REFERENCES accounts(id),
 side text NOT NULL CHECK (side IN ('debit','credit')),
 amount_minor bigint NOT NULL CHECK (amount_minor > 0),
 currency char(3) NOT NULL,
 source_type text,
 source_id uuid,
 created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_journal_entries_journal ON journal_entries(journal_id);
CREATE INDEX IF NOT EXISTS idx_journal_entries_account ON journal_entries(account_id);
CREATE INDEX IF NOT EXISTS idx_journal_entries_source ON journal_entries(source_type, source_id);

CREATE TABLE IF NOT EXISTS mail_events (
 id uuid PRIMARY KEY,
 source_message_id uuid NOT NULL,
 event_type text NOT NULL,
 payload jsonb NOT NULL DEFAULT '{}'::jsonb,
 confidence numeric(5,4) CHECK (confidence >= 0 AND confidence <= 1),
 status text NOT NULL DEFAULT 'pending',
 created_at timestamptz NOT NULL DEFAULT now(),
 processed_at timestamptz,
 UNIQUE(source_message_id, event_type)
);
CREATE INDEX IF NOT EXISTS idx_mail_events_status ON mail_events(status, created_at);
