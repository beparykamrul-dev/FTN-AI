-- FTN durable execution primitives. Additive/idempotent; no destructive changes.

CREATE TABLE IF NOT EXISTS event_journal (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID REFERENCES tenants(id) ON DELETE SET NULL,
  event_type TEXT NOT NULL,
  sequence BIGINT NOT NULL,
  correlation_id TEXT NOT NULL,
  causation_id TEXT NOT NULL DEFAULT '',
  aggregate_id TEXT NOT NULL DEFAULT '',
  payload JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (tenant_id, sequence)
);
CREATE INDEX IF NOT EXISTS event_journal_tenant_sequence_idx ON event_journal(tenant_id, sequence);
CREATE INDEX IF NOT EXISTS event_journal_correlation_idx ON event_journal(correlation_id);

CREATE TABLE IF NOT EXISTS event_consumer_offsets (
  consumer_id TEXT NOT NULL,
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  sequence BIGINT NOT NULL DEFAULT 0,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (consumer_id, tenant_id)
);

CREATE TABLE IF NOT EXISTS durable_jobs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID REFERENCES tenants(id) ON DELETE SET NULL,
  idempotency_key TEXT NOT NULL UNIQUE,
  job_type TEXT NOT NULL,
  payload JSONB NOT NULL DEFAULT '{}'::jsonb,
  status TEXT NOT NULL DEFAULT 'queued' CHECK (status IN ('queued','running','succeeded','failed','cancelled')),
  attempts INTEGER NOT NULL DEFAULT 0,
  max_attempts INTEGER NOT NULL DEFAULT 3 CHECK (max_attempts > 0),
  available_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  locked_at TIMESTAMPTZ,
  locked_by TEXT,
  started_at TIMESTAMPTZ,
  finished_at TIMESTAMPTZ,
  last_error TEXT NOT NULL DEFAULT '',
  correlation_id TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS durable_jobs_queue_idx ON durable_jobs(status, available_at);
CREATE INDEX IF NOT EXISTS durable_jobs_tenant_idx ON durable_jobs(tenant_id);
CREATE INDEX IF NOT EXISTS durable_jobs_correlation_idx ON durable_jobs(correlation_id);

CREATE TABLE IF NOT EXISTS execution_attempts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  job_id UUID NOT NULL REFERENCES durable_jobs(id) ON DELETE CASCADE,
  attempt_no INTEGER NOT NULL,
  worker_id TEXT NOT NULL,
  status TEXT NOT NULL CHECK (status IN ('running','succeeded','failed')),
  error TEXT NOT NULL DEFAULT '',
  started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  finished_at TIMESTAMPTZ,
  UNIQUE (job_id, attempt_no)
);
CREATE INDEX IF NOT EXISTS execution_attempts_job_idx ON execution_attempts(job_id, attempt_no DESC);

INSERT INTO permissions(key, description) VALUES
 ('event.read','Read durable event journal'),
 ('event.append','Append durable events'),
 ('event.commit','Commit consumer offsets'),
 ('job.read','Read durable job state'),
 ('job.submit','Submit durable jobs'),
 ('job.execute','Execute durable jobs')
ON CONFLICT (key) DO NOTHING;
