-- Scope API idempotency storage to the authenticated tenant and principal.
-- The logical Idempotency-Key remains unchanged at the API boundary; storage is prefixed.

UPDATE idempotency_keys i
SET key = p.tenant_id::text || ':' || p.id::text || ':' || i.key
FROM principals p
WHERE i.principal_id = p.id
  AND left(i.key, length(p.tenant_id::text) + length(p.id::text) + 2) <> p.tenant_id::text || ':' || p.id::text || ':';

INSERT INTO schema_migrations(version,name)
VALUES (19,'019_tenant_scoped_api_idempotency.sql')
ON CONFLICT (version) DO NOTHING;
