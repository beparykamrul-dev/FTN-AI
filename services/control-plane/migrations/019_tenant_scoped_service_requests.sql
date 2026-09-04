-- Bind persisted service requests to the authenticated tenant/principal.
-- Existing rows remain readable during upgrade; new API writes require tenant ownership.
ALTER TABLE service_requests
  ADD COLUMN IF NOT EXISTS tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
  ADD COLUMN IF NOT EXISTS principal_id UUID REFERENCES principals(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS service_requests_tenant_idx
  ON service_requests(tenant_id, created_at DESC);
CREATE INDEX IF NOT EXISTS service_requests_principal_idx
  ON service_requests(principal_id, created_at DESC);

INSERT INTO schema_migrations(version,name)
VALUES (19,'019_tenant_scoped_service_requests.sql')
ON CONFLICT (version) DO NOTHING;
