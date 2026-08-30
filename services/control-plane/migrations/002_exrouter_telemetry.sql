ALTER TABLE control_nodes ADD COLUMN IF NOT EXISTS jitter_ms DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE control_nodes ADD COLUMN IF NOT EXISTS retransmissions DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE control_nodes ADD COLUMN IF NOT EXISTS capacity_mbps DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE control_nodes ADD COLUMN IF NOT EXISTS utilization_percent DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE control_nodes ADD COLUMN IF NOT EXISTS bgp_up BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE control_nodes ADD COLUMN IF NOT EXISTS bfd_up BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE control_nodes ADD COLUMN IF NOT EXISTS isis_up BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE control_nodes ADD COLUMN IF NOT EXISTS evpn_ready BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE control_nodes ADD COLUMN IF NOT EXISTS anycast_ready BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE control_nodes ADD COLUMN IF NOT EXISTS rpki_valid BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE control_nodes ADD COLUMN IF NOT EXISTS prefix_count INTEGER NOT NULL DEFAULT 0;

CREATE INDEX IF NOT EXISTS control_nodes_bgp_health_idx ON control_nodes(bgp_up, bfd_up);
CREATE INDEX IF NOT EXISTS control_nodes_rpki_idx ON control_nodes(rpki_valid);
CREATE INDEX IF NOT EXISTS control_nodes_utilization_idx ON control_nodes(utilization_percent);

CREATE TABLE IF NOT EXISTS exrouter_route_events (
  id BIGSERIAL PRIMARY KEY,
  service_id TEXT NOT NULL,
  path_id TEXT NOT NULL,
  score DOUBLE PRECISION NOT NULL,
  decision TEXT NOT NULL,
  reason TEXT,
  actor TEXT,
  approval_required BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS exrouter_route_events_service_idx ON exrouter_route_events(service_id, created_at DESC);
CREATE INDEX IF NOT EXISTS exrouter_route_events_path_idx ON exrouter_route_events(path_id, created_at DESC);
