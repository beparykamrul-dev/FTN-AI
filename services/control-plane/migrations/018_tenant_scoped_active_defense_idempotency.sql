-- Scope active-defense idempotency to the owning tenant.
-- New active-defense executions must carry a tenant; legacy NULL rows remain readable.
UPDATE active_defense_executions
SET idempotency_key = tenant_id::text || ':' || idempotency_key
WHERE tenant_id IS NOT NULL
  AND left(idempotency_key, length(tenant_id::text) + 1) <> tenant_id::text || ':';

ALTER TABLE active_defense_executions
  DROP CONSTRAINT IF EXISTS active_defense_executions_idempotency_key_key;

-- Tenant ownership is part of the active-defense lifecycle. Cascading tenant deletion
-- avoids the previous FK/trigger conflict where ON DELETE SET NULL was rejected by
-- the tenant-required trigger below.
ALTER TABLE active_defense_executions
  DROP CONSTRAINT IF EXISTS active_defense_executions_tenant_id_fkey;
ALTER TABLE active_defense_executions
  ADD CONSTRAINT active_defense_executions_tenant_id_fkey
  FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE;

CREATE UNIQUE INDEX IF NOT EXISTS active_defense_executions_tenant_idempotency_uq
  ON active_defense_executions(tenant_id, idempotency_key);

CREATE OR REPLACE FUNCTION ftn_scope_active_defense_idempotency_key() RETURNS trigger
LANGUAGE plpgsql AS $$
DECLARE
  prefix TEXT;
BEGIN
  IF NEW.tenant_id IS NULL THEN
    RAISE EXCEPTION 'active-defense tenant_id is required';
  END IF;
  prefix := NEW.tenant_id::text || ':';
  IF left(NEW.idempotency_key, length(prefix)) <> prefix THEN
    NEW.idempotency_key := prefix || NEW.idempotency_key;
  END IF;
  RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS active_defense_scope_idempotency ON active_defense_executions;
CREATE TRIGGER active_defense_scope_idempotency
BEFORE INSERT OR UPDATE OF tenant_id, idempotency_key
ON active_defense_executions
FOR EACH ROW EXECUTE FUNCTION ftn_scope_active_defense_idempotency_key();

INSERT INTO schema_migrations(version,name)
VALUES (18,'018_tenant_scoped_active_defense_idempotency.sql')
ON CONFLICT (version) DO NOTHING;
