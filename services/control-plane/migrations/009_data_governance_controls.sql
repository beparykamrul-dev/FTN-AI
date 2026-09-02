-- FTN Data Governor controls: quality rules/results and governance indexes.
CREATE TABLE IF NOT EXISTS data_quality_rules (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
  asset_id UUID NOT NULL REFERENCES data_assets(id) ON DELETE CASCADE,
  rule_key TEXT NOT NULL,
  dimension TEXT NOT NULL CHECK (dimension IN ('completeness','validity','uniqueness','timeliness','consistency')),
  rule_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('draft','active','retired')),
  created_by UUID REFERENCES principals(id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (tenant_id, asset_id, rule_key)
);
CREATE INDEX IF NOT EXISTS data_quality_rules_tenant_idx ON data_quality_rules(tenant_id);
CREATE INDEX IF NOT EXISTS data_quality_rules_asset_idx ON data_quality_rules(asset_id);

CREATE TABLE IF NOT EXISTS data_quality_results (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
  rule_id UUID NOT NULL REFERENCES data_quality_rules(id) ON DELETE CASCADE,
  score NUMERIC(6,3) NOT NULL CHECK (score >= 0 AND score <= 100),
  passed BOOLEAN NOT NULL,
  sample_size BIGINT NOT NULL DEFAULT 0 CHECK (sample_size >= 0),
  observed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  details JSONB NOT NULL DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS data_quality_results_tenant_idx ON data_quality_results(tenant_id);
CREATE INDEX IF NOT EXISTS data_quality_results_rule_idx ON data_quality_results(rule_id);
CREATE INDEX IF NOT EXISTS data_quality_results_observed_idx ON data_quality_results(observed_at DESC);

CREATE INDEX IF NOT EXISTS data_requests_asset_idx ON data_requests(asset_id);
CREATE INDEX IF NOT EXISTS data_requests_approval_idx ON data_requests(approval_id);
CREATE INDEX IF NOT EXISTS data_access_events_tenant_created_idx ON data_access_events(tenant_id, created_at DESC);
CREATE INDEX IF NOT EXISTS data_lineage_tenant_recorded_idx ON data_lineage(tenant_id, recorded_at DESC);
