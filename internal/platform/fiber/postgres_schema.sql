-- FTN GIS/Fiber persistence foundation for PostgreSQL + PostGIS.
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS fiber_nodes (
    id text PRIMARY KEY,
    kind text NOT NULL,
    latitude double precision NOT NULL,
    longitude double precision NOT NULL,
    status text NOT NULL DEFAULT 'unknown',
    geom geometry(Point, 4326) GENERATED ALWAYS AS (ST_SetSRID(ST_MakePoint(longitude, latitude), 4326)) STORED,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_fiber_nodes_geom ON fiber_nodes USING GIST (geom);
CREATE INDEX IF NOT EXISTS idx_fiber_nodes_kind_status ON fiber_nodes (kind, status);

CREATE TABLE IF NOT EXISTS fiber_links (
    id text PRIMARY KEY,
    from_node text NOT NULL REFERENCES fiber_nodes(id),
    to_node text NOT NULL REFERENCES fiber_nodes(id),
    distance_meters double precision NOT NULL DEFAULT 0,
    status text NOT NULL DEFAULT 'unknown',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_fiber_links_nodes ON fiber_links (from_node, to_node);
CREATE INDEX IF NOT EXISTS idx_fiber_links_status ON fiber_links (status);

CREATE TABLE IF NOT EXISTS ftn_customers (
    id text PRIMARY KEY,
    name text,
    service_id text,
    package_name text,
    pppoe_user text,
    onu_id text,
    router_id text,
    ip inet,
    service_status text NOT NULL DEFAULT 'unknown',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_ftn_customers_service ON ftn_customers (service_id);
CREATE INDEX IF NOT EXISTS idx_ftn_customers_onu ON ftn_customers (onu_id);
CREATE INDEX IF NOT EXISTS idx_ftn_customers_pppoe ON ftn_customers (pppoe_user);

CREATE TABLE IF NOT EXISTS fiber_incidents (
    id text PRIMARY KEY,
    link_id text NOT NULL REFERENCES fiber_links(id),
    status text NOT NULL,
    detected_at timestamptz NOT NULL DEFAULT now(),
    recovery_confidence double precision,
    recovery_reason text
);

CREATE INDEX IF NOT EXISTS idx_fiber_incidents_link_time ON fiber_incidents (link_id, detected_at DESC);
