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
