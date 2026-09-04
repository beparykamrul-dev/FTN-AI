-- Scope API idempotency storage to the authenticated tenant and principal.
-- The logical Idempotency-Key remains unchanged at the API boundary; storage is prefixed.

UPDATE idempotency_keys i
SET key = p.tenant_id::text || ':' || p.id::text || ':' || i.key
FROM principals p
WHERE i.principal_id = p.id
  AND left(i.key, length(p.tenant_id::text) + length(p.id::text) + 2) <> p.tenant_id::text || ':' || p.id::text || ':';

CREATE OR REPLACE FUNCTION ftn_require_api_idempotency_principal() RETURNS trigger
LANGUAGE plpgsql AS $$
BEGIN
  IF NEW.principal_id IS NULL THEN
    RAISE EXCEPTION 'api idempotency principal_id is required';
  END IF;
  RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS idempotency_keys_require_principal ON idempotency_keys;
CREATE TRIGGER idempotency_keys_require_principal
BEFORE INSERT OR UPDATE OF principal_id
ON idempotency_keys
FOR EACH ROW EXECUTE FUNCTION ftn_require_api_idempotency_principal();

INSERT INTO schema_migrations(version,name)
VALUES (20,'020_tenant_scoped_api_idempotency.sql')
ON CONFLICT (version) DO NOTHING;
