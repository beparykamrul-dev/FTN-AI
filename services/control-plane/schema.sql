create table if not exists service_requests (
  id bigserial primary key,
  service_id text not null,
  device_brand text,
  model text,
  mac text,
  serial text,
  scope text,
  status text not null default 'accepted',
  created_at timestamptz not null default now()
);
create index if not exists service_requests_service_id_idx on service_requests(service_id);
create index if not exists service_requests_mac_idx on service_requests(mac);
