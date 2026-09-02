-- FTN QKD inventory/security metadata. Raw key material is intentionally excluded.
CREATE TABLE IF NOT EXISTS qkd_nodes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  role TEXT NOT NULL CHECK (role IN ('tx','rx','kme','kms','trusted-node')),
  vendor TEXT NOT NULL DEFAULT '',
  model TEXT NOT NULL DEFAULT '',
  endpoint_ref TEXT NOT NULL DEFAULT '',
  kms_ref TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'unknown' CHECK (status IN ('unknown','healthy','degraded','offline','maintenance')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS qkd_nodes_tenant_idx ON qkd_nodes(tenant_id);
CREATE INDEX IF NOT EXISTS qkd_nodes_status_idx ON qkd_nodes(status);

CREATE TABLE IF NOT EXISTS qkd_links (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
  source_node_id UUID NOT NULL REFERENCES qkd_nodes(id) ON DELETE CASCADE,
  target_node_id UUID NOT NULL REFERENCES qkd_nodes(id) ON DELETE CASCADE,
  quantum_channel TEXT NOT NULL,
  classical_channel TEXT NOT NULL,
  authenticated BOOLEAN NOT NULL DEFAULT false,
  status TEXT NOT NULL DEFAULT 'unknown' CHECK (status IN ('unknown','up','degraded','down','maintenance')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CHECK (source_node_id <> target_node_id)
);
CREATE INDEX IF NOT EXISTS qkd_links_tenant_idx ON qkd_links(tenant_id);
CREATE INDEX IF NOT EXISTS qkd_links_source_idx ON qkd_links(source_node_id);
CREATE INDEX IF NOT EXISTS qkd_links_target_idx ON qkd_links(target_node_id);

CREATE TABLE IF NOT EXISTS qkd_status (
  node_id UUID PRIMARY KEY REFERENCES qkd_nodes(id) ON DELETE CASCADE,
  pool_state TEXT NOT NULL DEFAULT 'unknown' CHECK (pool_state IN ('unknown','healthy','degraded','empty','offline')),
  available_keys BIGINT NOT NULL DEFAULT 0 CHECK (available_keys >= 0),
  generation_rate_bps BIGINT NOT NULL DEFAULT 0 CHECK (generation_rate_bps >= 0),
  consumption_rate_bps BIGINT NOT NULL DEFAULT 0 CHECK (consumption_rate_bps >= 0),
  healthy BOOLEAN NOT NULL DEFAULT false,
  fallback_mode TEXT NOT NULL DEFAULT 'deny' CHECK (fallback_mode IN ('deny','pqc','hybrid')),
  observed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS qkd_kms (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  endpoint_ref TEXT NOT NULL,
  api_profile TEXT NOT NULL DEFAULT 'etsi-gs-qkd-020-v1.1.1',
  status TEXT NOT NULL DEFAULT 'unknown' CHECK (status IN ('unknown','healthy','degraded','offline')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (tenant_id, name)
);
CREATE INDEX IF NOT EXISTS qkd_kms_tenant_idx ON qkd_kms(tenant_id);

INSERT INTO permissions(key, description) VALUES
 ('qkd.read','Read QKD node, link, KMS and status metadata'),
 ('qkd.change','Request approval-gated QKD policy or consumer changes')
ON CONFLICT (key) DO NOTHING;
