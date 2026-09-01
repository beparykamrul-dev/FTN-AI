-- FTN platform foundation. Additive/idempotent migration: preserves existing tables.
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS tenants (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  slug TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','suspended','disabled')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS principals (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  subject TEXT NOT NULL,
  kind TEXT NOT NULL,
  issuer TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','suspended','revoked')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (tenant_id, issuer, subject)
);
CREATE INDEX IF NOT EXISTS principals_tenant_idx ON principals(tenant_id);

CREATE TABLE IF NOT EXISTS roles (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (tenant_id, name)
);

CREATE TABLE IF NOT EXISTS permissions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  key TEXT NOT NULL UNIQUE,
  description TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS role_permissions (
  role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
  permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
  PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE IF NOT EXISTS principal_roles (
  principal_id UUID NOT NULL REFERENCES principals(id) ON DELETE CASCADE,
  role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
  PRIMARY KEY (principal_id, role_id)
);

CREATE TABLE IF NOT EXISTS api_credentials (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  principal_id UUID NOT NULL REFERENCES principals(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  token_hash TEXT NOT NULL UNIQUE,
  expires_at TIMESTAMPTZ,
  revoked_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS api_credentials_principal_idx ON api_credentials(principal_id);

CREATE TABLE IF NOT EXISTS service_registry (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  service_id TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  version TEXT NOT NULL DEFAULT '0.0.0',
  endpoint TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'unknown',
  health_url TEXT NOT NULL DEFAULT '',
  dependencies TEXT[] NOT NULL DEFAULT '{}',
  capabilities JSONB NOT NULL DEFAULT '{}'::jsonb,
  last_heartbeat_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS service_registry_status_idx ON service_registry(status);

CREATE TABLE IF NOT EXISTS service_entitlements (
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  principal_id UUID NOT NULL REFERENCES principals(id) ON DELETE CASCADE,
  service_id TEXT NOT NULL,
  active BOOLEAN NOT NULL DEFAULT true,
  expires_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (principal_id, service_id)
);
CREATE INDEX IF NOT EXISTS service_entitlements_tenant_idx ON service_entitlements(tenant_id);

CREATE TABLE IF NOT EXISTS change_approvals (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID REFERENCES tenants(id) ON DELETE SET NULL,
  requested_by UUID REFERENCES principals(id) ON DELETE SET NULL,
  action TEXT NOT NULL,
  resource TEXT NOT NULL,
  request_hash TEXT NOT NULL UNIQUE,
  status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','approved','rejected','expired','executed','rolled_back')),
  approved_by UUID REFERENCES principals(id) ON DELETE SET NULL,
  expires_at TIMESTAMPTZ,
  executed_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS change_approvals_status_idx ON change_approvals(status);
CREATE INDEX IF NOT EXISTS change_approvals_tenant_idx ON change_approvals(tenant_id);

CREATE TABLE IF NOT EXISTS audit_events (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID REFERENCES tenants(id) ON DELETE SET NULL,
  principal_id UUID REFERENCES principals(id) ON DELETE SET NULL,
  action TEXT NOT NULL,
  resource TEXT NOT NULL,
  outcome TEXT NOT NULL DEFAULT 'success',
  correlation_id TEXT NOT NULL,
  request_id TEXT NOT NULL,
  metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS audit_events_created_idx ON audit_events(created_at DESC);
CREATE INDEX IF NOT EXISTS audit_events_correlation_idx ON audit_events(correlation_id);
CREATE INDEX IF NOT EXISTS audit_events_principal_idx ON audit_events(principal_id);

CREATE TABLE IF NOT EXISTS idempotency_keys (
  key TEXT PRIMARY KEY,
  principal_id UUID REFERENCES principals(id) ON DELETE SET NULL,
  request_hash TEXT NOT NULL,
  response_status INTEGER,
  response_body JSONB,
  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idempotency_expires_idx ON idempotency_keys(expires_at);

INSERT INTO permissions(key, description) VALUES
 ('service.read','Read service catalog and state'),
 ('service.request','Create service requests'),
 ('service.manage','Manage service lifecycle'),
 ('node.read','Read node state'),
 ('node.manage','Manage node registration'),
 ('network.read','Read network health'),
 ('network.plan','Generate network recommendations'),
 ('network.change','Apply approved network changes'),
 ('approval.read','Read approvals'),
 ('approval.manage','Approve or reject privileged changes'),
 ('audit.read','Read audit events')
ON CONFLICT (key) DO NOTHING;

INSERT INTO service_registry(service_id,name,version,status,capabilities)
SELECT v.service_id,v.name,'1.0.0','registered',v.capabilities::jsonb
FROM (VALUES
 ('internet','FTN Internet','{"class":"core","domain":"connectivity"}'),
 ('ftndns','FTNDNS / Friendly DNS','{"class":"core","domain":"dns"}'),
 ('hosting','FTN Hosting','{"class":"platform","domain":"hosting"}'),
 ('cloud','FTN Cloud','{"class":"platform","domain":"storage"}'),
 ('drive','FTN Drive','{"class":"platform","domain":"storage"}'),
 ('cctv','CCTV Cloud','{"class":"platform","domain":"cctv"}'),
 ('fibermap','FTN FiberMap','{"class":"network","domain":"fiber"}'),
 ('ai','FTN AI Assistant','{"class":"control","domain":"ai"}'),
 ('media','FTN Media','{"class":"platform","domain":"media"}'),
 ('tv','FTN TV Player','{"class":"platform","domain":"media"}'),
 ('appstore','FTN App Store','{"class":"platform","domain":"apps"}'),
 ('mail','FTN Mail','{"class":"platform","domain":"messaging"}'),
 ('ecommerce','FTN E-Commerce','{"class":"platform","domain":"commerce"}'),
 ('developer','FTN Developer Platform','{"class":"platform","domain":"developer"}'),
 ('device-care','FTN Device Care','{"class":"platform","domain":"device"}'),
 ('codec','FTN Codec Fabric','{"class":"fabric","domain":"encoding"}'),
 ('media-processing','FTN Media Processing','{"class":"fabric","domain":"processing"}'),
 ('e2e-transfer','FTN E2E Transfer','{"class":"fabric","domain":"transfer"}')
) AS v(service_id,name,capabilities)
ON CONFLICT (service_id) DO UPDATE SET name=EXCLUDED.name, capabilities=EXCLUDED.capabilities, updated_at=now();
