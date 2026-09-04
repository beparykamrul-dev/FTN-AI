#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

POSTGRES_IMAGE="postgres:17-alpine"
NAME="ftn-db-restore-${GITHUB_RUN_ID:-local}-$$"
DUMP="$(mktemp)"
RESTORE_DUMP="$(mktemp)"
cleanup() {
  docker rm -f "$NAME" >/dev/null 2>&1 || true
  rm -f "$DUMP" "$RESTORE_DUMP"
}
trap cleanup EXIT

for required in \
  services/control-plane/schema.sql \
  services/control-plane/migrations/002_platform_foundation.sql \
  services/control-plane/migrations/003_durable_execution.sql \
  services/control-plane/migrations/004_approval_execution.sql \
  services/control-plane/migrations/005_job_integrity_leases.sql \
  services/control-plane/migrations/006_job_event_automation.sql \
  services/control-plane/migrations/006_job_event_triggers.sql \
  services/control-plane/migrations/007_qkd_security.sql \
  services/control-plane/migrations/008_data_governor.sql \
  services/control-plane/migrations/009_data_governance_controls.sql \
  services/control-plane/migrations/010_dns_guard_filtering.sql \
  services/control-plane/migrations/011_migration_registry.sql \
  services/control-plane/migrations/012_control_nodes_tenant_scope.sql \
  services/control-plane/migrations/013_approval_payload_binding.sql \
  services/control-plane/migrations/014_tenant_scoped_approval_hash.sql \
  services/control-plane/migrations/015_approval_event_trigger.sql \
  services/control-plane/migrations/016_tenant_scoped_job_idempotency.sql; do
  test -f "$required"
done

echo "Starting isolated PostgreSQL validation instance"
docker run -d --name "$NAME" \
  -e POSTGRES_DB=ftn \
  -e POSTGRES_USER=ftn \
  -e POSTGRES_PASSWORD=ci-only-password \
  "$POSTGRES_IMAGE" >/dev/null

ready_deadline=$((SECONDS + 90))
until docker exec "$NAME" pg_isready -U ftn -d ftn >/dev/null 2>&1; do
  if (( SECONDS >= ready_deadline )); then
    echo "PostgreSQL readiness timeout" >&2
    exit 1
  fi
  sleep 1
done

for migration in \
  services/control-plane/schema.sql \
  services/control-plane/migrations/002_platform_foundation.sql \
  services/control-plane/migrations/003_durable_execution.sql \
  services/control-plane/migrations/004_approval_execution.sql \
  services/control-plane/migrations/005_job_integrity_leases.sql \
  services/control-plane/migrations/006_job_event_automation.sql \
  services/control-plane/migrations/006_job_event_triggers.sql \
  services/control-plane/migrations/007_qkd_security.sql \
  services/control-plane/migrations/008_data_governor.sql \
  services/control-plane/migrations/009_data_governance_controls.sql \
  services/control-plane/migrations/010_dns_guard_filtering.sql \
  services/control-plane/migrations/011_migration_registry.sql \
  services/control-plane/migrations/012_control_nodes_tenant_scope.sql \
  services/control-plane/migrations/013_approval_payload_binding.sql \
  services/control-plane/migrations/014_tenant_scoped_approval_hash.sql \
  services/control-plane/migrations/015_approval_event_trigger.sql \
  services/control-plane/migrations/016_tenant_scoped_job_idempotency.sql; do
  echo "Applying $migration"
  docker exec -i "$NAME" psql -v ON_ERROR_STOP=1 -U ftn -d ftn < "$migration" >/dev/null
done

echo "Writing deterministic integrity fixture"
docker exec -i "$NAME" psql -v ON_ERROR_STOP=1 -U ftn -d ftn <<'SQL' >/dev/null
INSERT INTO tenants (slug, name) VALUES ('ci-restore', 'CI Restore Validation') ON CONFLICT (slug) DO NOTHING;
INSERT INTO tenants (slug, name) VALUES ('ci-restore-b', 'CI Restore Validation B') ON CONFLICT (slug) DO NOTHING;
SQL

before="$(docker exec "$NAME" psql -At -U ftn -d ftn -c "SELECT count(*) FROM tenants WHERE slug='ci-restore';")"
test "$before" = "1"

echo "Validating tenant-scoped durable-job idempotency"
docker exec -i "$NAME" psql -v ON_ERROR_STOP=1 -U ftn -d ftn <<'SQL' >/dev/null
INSERT INTO durable_jobs(tenant_id, idempotency_key, job_type, correlation_id)
SELECT id, 'ci-shared-key', 'ci.test', 'ci-restore-a'
FROM tenants WHERE slug='ci-restore'
ON CONFLICT DO NOTHING;
INSERT INTO durable_jobs(tenant_id, idempotency_key, job_type, correlation_id)
SELECT id, 'ci-shared-key', 'ci.test', 'ci-restore-b'
FROM tenants WHERE slug='ci-restore-b'
ON CONFLICT DO NOTHING;
SQL

scoped_jobs="$(docker exec "$NAME" psql -At -U ftn -d ftn -c "SELECT count(*) FROM durable_jobs WHERE idempotency_key LIKE '%:ci-shared-key' AND tenant_id IN (SELECT id FROM tenants WHERE slug IN ('ci-restore','ci-restore-b'));" )"
test "$scoped_jobs" = "2"

idempotency_trigger="$(docker exec "$NAME" psql -At -U ftn -d ftn -c "SELECT count(*) FROM pg_trigger WHERE tgname='durable_jobs_scope_idempotency';")"
test "$idempotency_trigger" = "1"

echo "Creating custom-format backup"
docker exec "$NAME" pg_dump -U ftn -d ftn -Fc > "$DUMP"
test -s "$DUMP"

docker exec "$NAME" createdb -U ftn ftn_restore
docker exec -i "$NAME" pg_restore -U ftn -d ftn_restore --exit-on-error < "$DUMP"

after="$(docker exec "$NAME" psql -At -U ftn -d ftn_restore -c "SELECT count(*) FROM tenants WHERE slug='ci-restore';")"
test "$after" = "1"

event_trigger="$(docker exec "$NAME" psql -At -U ftn -d ftn -c "SELECT count(*) FROM pg_trigger WHERE tgname='durable_jobs_event_journal';")"
test "$event_trigger" = "1"

approval_event_trigger="$(docker exec "$NAME" psql -At -U ftn -d ftn -c "SELECT count(*) FROM pg_trigger WHERE tgname='change_approvals_lifecycle_event';")"
test "$approval_event_trigger" = "1"

schema_version="$(docker exec "$NAME" psql -At -U ftn -d ftn -c "SELECT to_regclass('public.event_journal') IS NOT NULL AND to_regclass('public.durable_jobs') IS NOT NULL;")"
test "$schema_version" = "t"

registry_count="$(docker exec "$NAME" psql -At -U ftn -d ftn -c "SELECT count(*) FROM schema_migrations WHERE version >= 11 AND version <= 16;")"
test "$registry_count" = "6"

control_nodes_tenant="$(docker exec "$NAME" psql -At -U ftn -d ftn -c "SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema='public' AND table_name='control_nodes' AND column_name='tenant_id');")"
test "$control_nodes_tenant" = "t"

approval_payload="$(docker exec "$NAME" psql -At -U ftn -d ftn -c "SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema='public' AND table_name='change_approvals' AND column_name='payload_hash');")"
test "$approval_payload" = "t"

approval_scoped_unique="$(docker exec "$NAME" psql -At -U ftn -d ftn -c "SELECT EXISTS (SELECT 1 FROM pg_indexes WHERE schemaname='public' AND indexname='change_approvals_tenant_request_hash_uq');")"
test "$approval_scoped_unique" = "t"

echo "PostgreSQL migration + backup + restore validation passed."
