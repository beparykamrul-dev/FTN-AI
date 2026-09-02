-- FTN active-defense execution audit/state. Additive and idempotent.
CREATE TABLE IF NOT EXISTS active_defense_executions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID REFERENCES tenants(id) ON DELETE SET NULL,
  job_id UUID REFERENCES durable_jobs(id) ON DELETE SET NULL,
  alert_hash TEXT NOT NULL,
  idempotency_key TEXT NOT NULL UNIQUE,
  target_asset TEXT NOT NULL,
  target_scope TEXT NOT NULL CHECK (target_scope = 'ftn-owned-asset'),
  adapter TEXT NOT NULL CHECK (adapter IN ('nftables','xdp','mikrotik')),
  operation TEXT NOT NULL CHECK (operation IN ('temporary-containment','rate-limit','health-check','temporary-address-list')),
  duration_seconds INTEGER NOT NULL CHECK (duration_seconds > 0 AND duration_seconds <= 3600),
  snapshot_required BOOLEAN NOT NULL DEFAULT true,
  verification_required BOOLEAN NOT NULL DEFAULT true,
  status TEXT NOT NULL DEFAULT 'planned' CHECK (status IN ('planned','queued','running','succeeded','failed','rolled_back','cancelled')),
  snapshot_ref TEXT NOT NULL DEFAULT '',
  verification_ref TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS active_defense_executions_tenant_idx ON active_defense_executions(tenant_id, created_at DESC);
CREATE INDEX IF NOT EXISTS active_defense_executions_job_idx ON active_defense_executions(job_id);

INSERT INTO permissions(key, description) VALUES
 ('security.active-defense.read','Read FTN active-defense execution state'),
 ('security.active-defense.execute','Execute approved/bounded FTN-owned containment jobs')
ON CONFLICT (key) DO NOTHING;
