-- FTN Data Governor foundation. Stores governance metadata and policy state;
-- application payloads and secrets remain outside these governance tables.
CREATE TABLE IF NOT EXISTS data_assets (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
  asset_key TEXT NOT NULL,
  name TEXT NOT NULL,
  domain TEXT NOT NULL,
  classification TEXT NOT NULL DEFAULT 'internal' CHECK (classification IN ('public','internal','confidential','restricted')),
  owner_principal_id UUID REFERENCES principals(id) ON DELETE SET NULL,
  steward_principal_id UUID REFERENCES principals(id) ON DELETE SET NULL,
  source_ref TEXT NOT NULL DEFAULT '',
  retention_days INTEGER CHECK (retention_days IS NULL OR retention_days >= 0),
  status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','archived','pending_deletion','deleted')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (tenant_id, asset_key)
);
CREATE INDEX IF NOT EXISTS data_assets_tenant_idx ON data_assets(tenant_id);
CREATE INDEX IF NOT EXISTS data_assets_classification_idx ON data_assets(classification);

CREATE TABLE IF NOT EXISTS data_policies (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
  policy_key TEXT NOT NULL,
  version INTEGER NOT NULL DEFAULT 1 CHECK (version > 0),
  policy_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('draft','active','retired')),
  created_by UUID REFERENCES principals(id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (tenant_id, policy_key, version)
);
CREATE INDEX IF NOT EXISTS data_policies_tenant_idx ON data_policies(tenant_id);

CREATE TABLE IF NOT EXISTS data_lineage (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
  source_asset_id UUID NOT NULL REFERENCES data_assets(id) ON DELETE CASCADE,
  target_asset_id UUID NOT NULL REFERENCES data_assets(id) ON DELETE CASCADE,
  transform_ref TEXT NOT NULL DEFAULT '',
  recorded_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CHECK (source_asset_id <> target_asset_id)
);
CREATE INDEX IF NOT EXISTS data_lineage_source_idx ON data_lineage(source_asset_id);
CREATE INDEX IF NOT EXISTS data_lineage_target_idx ON data_lineage(target_asset_id);

CREATE TABLE IF NOT EXISTS data_access_events (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
  asset_id UUID REFERENCES data_assets(id) ON DELETE SET NULL,
  principal_id UUID REFERENCES principals(id) ON DELETE SET NULL,
  action TEXT NOT NULL,
  decision TEXT NOT NULL CHECK (decision IN ('allow','deny','approval_required')),
  reason TEXT NOT NULL DEFAULT '',
  request_hash TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS data_access_events_asset_idx ON data_access_events(asset_id);
CREATE INDEX IF NOT EXISTS data_access_events_created_idx ON data_access_events(created_at);

CREATE TABLE IF NOT EXISTS data_requests (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
  asset_id UUID REFERENCES data_assets(id) ON DELETE SET NULL,
  request_type TEXT NOT NULL CHECK (request_type IN ('export','deletion','retention_override','classification_change','access')),
  status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','approved','rejected','completed','expired')),
  requested_by UUID REFERENCES principals(id) ON DELETE SET NULL,
  approval_id UUID REFERENCES change_approvals(id) ON DELETE SET NULL,
  request_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS data_requests_tenant_idx ON data_requests(tenant_id);
CREATE INDEX IF NOT EXISTS data_requests_status_idx ON data_requests(status);

INSERT INTO permissions(key, description) VALUES
 ('data.read','Read Data Governor metadata'),
 ('data.change','Request Data Governor policy changes'),
 ('data.export','Request governed data export'),
 ('data.delete','Request governed data deletion')
ON CONFLICT (key) DO NOTHING;
