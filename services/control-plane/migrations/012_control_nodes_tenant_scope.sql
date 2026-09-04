-- Bind control-plane nodes to a tenant. Existing unowned rows remain
-- intentionally inaccessible until re-registered by an authorized tenant.
ALTER TABLE control_nodes
  ADD COLUMN IF NOT EXISTS tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS control_nodes_tenant_idx
  ON control_nodes(tenant_id);

CREATE INDEX IF NOT EXISTS control_nodes_tenant_health_idx
  ON control_nodes(tenant_id, healthy);

INSERT INTO schema_migrations(version,name)
VALUES (12,'012_control_nodes_tenant_scope.sql')
ON CONFLICT (version) DO NOTHING;
