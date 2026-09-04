CREATE TABLE IF NOT EXISTS service_requests (
  id BIGSERIAL PRIMARY KEY,
  service_id TEXT NOT NULL,
  device_brand TEXT,
  model TEXT,
  mac TEXT,
  serial TEXT,
  scope TEXT,
  status TEXT NOT NULL DEFAULT 'accepted',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS service_requests_service_id_idx ON service_requests(service_id);
CREATE INDEX IF NOT EXISTS service_requests_mac_idx ON service_requests(mac);
CREATE INDEX IF NOT EXISTS service_requests_created_at_idx ON service_requests(created_at DESC);

CREATE TABLE IF NOT EXISTS control_nodes (
  id TEXT PRIMARY KEY,
  tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
  provider TEXT NOT NULL,
  region TEXT NOT NULL DEFAULT '',
  services TEXT[] NOT NULL DEFAULT '{}',
  cpu_percent DOUBLE PRECISION NOT NULL DEFAULT 0,
  ram_percent DOUBLE PRECISION NOT NULL DEFAULT 0,
  ssd_percent DOUBLE PRECISION NOT NULL DEFAULT 0,
  hdd_percent DOUBLE PRECISION NOT NULL DEFAULT 0,
  net_mbps DOUBLE PRECISION NOT NULL DEFAULT 0,
  latency_ms DOUBLE PRECISION NOT NULL DEFAULT 0,
  packet_loss_percent DOUBLE PRECISION NOT NULL DEFAULT 0,
  healthy BOOLEAN NOT NULL DEFAULT false,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS control_nodes_provider_idx ON control_nodes(provider);
CREATE INDEX IF NOT EXISTS control_nodes_region_idx ON control_nodes(region);
CREATE INDEX IF NOT EXISTS control_nodes_health_idx ON control_nodes(healthy);
CREATE INDEX IF NOT EXISTS control_nodes_tenant_idx ON control_nodes(tenant_id);
CREATE INDEX IF NOT EXISTS control_nodes_tenant_health_idx ON control_nodes(tenant_id, healthy);
