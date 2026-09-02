-- FTN DNS Guard: policy-driven DNS block/filter profiles.
-- Stores policy metadata only; raw DNS payloads and secrets are not stored.
CREATE TABLE IF NOT EXISTS dns_filter_profiles (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
  profile_key TEXT NOT NULL,
  name TEXT NOT NULL,
  mode TEXT NOT NULL CHECK (mode IN ('basic','family','strict','custom')),
  enabled BOOLEAN NOT NULL DEFAULT true,
  safe_search BOOLEAN NOT NULL DEFAULT false,
  youtube_restricted BOOLEAN NOT NULL DEFAULT false,
  block_ads BOOLEAN NOT NULL DEFAULT false,
  block_trackers BOOLEAN NOT NULL DEFAULT false,
  block_malware BOOLEAN NOT NULL DEFAULT true,
  block_phishing BOOLEAN NOT NULL DEFAULT true,
  block_adult BOOLEAN NOT NULL DEFAULT false,
  block_gambling BOOLEAN NOT NULL DEFAULT false,
  block_dating BOOLEAN NOT NULL DEFAULT false,
  block_distraction BOOLEAN NOT NULL DEFAULT false,
  custom_allowlist JSONB NOT NULL DEFAULT '[]'::jsonb,
  custom_blocklist JSONB NOT NULL DEFAULT '[]'::jsonb,
  created_by UUID REFERENCES principals(id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (tenant_id, profile_key)
);
CREATE INDEX IF NOT EXISTS dns_filter_profiles_tenant_idx ON dns_filter_profiles(tenant_id);
CREATE INDEX IF NOT EXISTS dns_filter_profiles_enabled_idx ON dns_filter_profiles(tenant_id, enabled);

CREATE TABLE IF NOT EXISTS dns_filter_bindings (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
  profile_id UUID NOT NULL REFERENCES dns_filter_profiles(id) ON DELETE CASCADE,
  subject_type TEXT NOT NULL CHECK (subject_type IN ('tenant','customer','device','network','vlan','pppoe')),
  subject_ref TEXT NOT NULL,
  enabled BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (tenant_id, subject_type, subject_ref)
);
CREATE INDEX IF NOT EXISTS dns_filter_bindings_profile_idx ON dns_filter_bindings(profile_id);
CREATE INDEX IF NOT EXISTS dns_filter_bindings_tenant_idx ON dns_filter_bindings(tenant_id);

CREATE TABLE IF NOT EXISTS dns_filter_events (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
  profile_id UUID REFERENCES dns_filter_profiles(id) ON DELETE SET NULL,
  subject_ref TEXT NOT NULL DEFAULT '',
  domain_hash TEXT NOT NULL DEFAULT '',
  category TEXT NOT NULL DEFAULT 'unknown',
  decision TEXT NOT NULL CHECK (decision IN ('allow','block','override')),
  reason TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS dns_filter_events_tenant_created_idx ON dns_filter_events(tenant_id, created_at DESC);

INSERT INTO permissions(key, description) VALUES
 ('dns.filter.read','Read DNS filtering policy metadata'),
 ('dns.filter.change','Request DNS filtering policy changes')
ON CONFLICT (key) DO NOTHING;

INSERT INTO dns_filter_profiles
  (tenant_id, profile_key, name, mode, safe_search, youtube_restricted, block_ads, block_trackers,
   block_malware, block_phishing, block_adult, block_gambling, block_dating, block_distraction)
SELECT t.id, v.profile_key, v.name, v.mode, v.safe_search, v.youtube_restricted, v.block_ads, v.block_trackers,
       v.block_malware, v.block_phishing, v.block_adult, v.block_gambling, v.block_dating, v.block_distraction
FROM tenants t
CROSS JOIN (VALUES
 ('guard-basic','FTN Guard Basic','basic',false,false,true,true,true,true,false,false,false,false),
 ('guard-family','FTN Guard Family','family',true,true,true,true,true,true,true,true,true,false),
 ('guard-strict','FTN Guard Strict','strict',true,true,true,true,true,true,true,true,true,true),
 ('kahf-guard-family','Kahf Guard Family Compatible','family',true,true,true,true,true,true,true,true,true,false)
) AS v(profile_key,name,mode,safe_search,youtube_restricted,block_ads,block_trackers,block_malware,block_phishing,block_adult,block_gambling,block_dating,block_distraction)
ON CONFLICT (tenant_id, profile_key) DO NOTHING;
