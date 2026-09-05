#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

POSTGRES_IMAGE="postgres:17-alpine"
NAME="ftn-db-restore-${GITHUB_RUN_ID:-local}-$$"
DUMP="$(mktemp)"
cleanup() { docker rm -f "$NAME" >/dev/null 2>&1 || true; rm -f "$DUMP"; }
trap cleanup EXIT

for required in \
  services/control-plane/schema.sql \
  services/control-plane/migrations/002_platform_foundation.sql \
  services/control-plane/migrations/003_durable_execution.sql \
  services/control-plane/migrations/004_approval_execution.sql \
  services/control-plane/migrations/005_job_integrity_leases.sql \
  services/control-plane/migrations/006_job_event_automation.sql \
  services/control-plane/migrations/007_qkd_security.sql \
  services/control-plane/migrations/008_data_governor.sql \
  services/control-plane/migrations/009_data_governance_controls.sql \
  services/control-plane/migrations/010_dns_guard_filtering.sql \
  services/control-plane/migrations/011_migration_registry.sql \
  services/control-plane/migrations/012_control_nodes_tenant_scope.sql \
  services/control-plane/migrations/013_approval_payload_binding.sql \
  services/control-plane/migrations/014_tenant_scoped_approval_hash.sql \
  services/control-plane/migrations/015_approval_event_trigger.sql \
  services/control-plane/migrations/016_tenant_scoped_job_idempotency.sql \
  services/control-plane/migrations/017_active_defense_execution.sql \
  services/control-plane/migrations/018_tenant_scoped_active_defense_idempotency.sql \
  services/control-plane/migrations/019_tenant_scoped_service_requests.sql \
  services/control-plane/migrations/020_tenant_scoped_api_idempotency.sql \
  services/control-plane/migrations/021_verification_identity.sql \
  services/control-plane/migrations/022_data_request_idempotency.sql \
  services/control-plane/migrations/023_execution_attempt_integrity.sql \
  services/control-plane/migrations/024_outbound_queue_leases.sql; do test -f "$required"; done

echo "Starting isolated PostgreSQL validation instance"
docker run -d --name "$NAME" -e POSTGRES_DB=ftn -e POSTGRES_USER=ftn -e POSTGRES_PASSWORD=ci-only-password "$POSTGRES_IMAGE" >/dev/null
ready_deadline=$((SECONDS + 90))
until docker exec "$NAME" pg_isready -U ftn -d ftn >/dev/null 2>&1; do
  if (( SECONDS >= ready_deadline )); then echo "PostgreSQL readiness timeout" >&2; exit 1; fi
  sleep 1
done

for migration in \
  services/control-plane/schema.sql \
  services/control-plane/migrations/002_platform_foundation.sql \
  services/control-plane/migrations/003_durable_execution.sql \
  services/control-plane/migrations/004_approval_execution.sql \
  services/control-plane/migrations/005_job_integrity_leases.sql \
  services/control-plane/migrations/006_job_event_automation.sql \
  services/control-plane/migrations/007_qkd_security.sql \
  services/control-plane/migrations/008_data_governor.sql \
  services/control-plane/migrations/009_data_governance_controls.sql \
  services/control-plane/migrations/010_dns_guard_filtering.sql \
  services/control-plane/migrations/011_migration_registry.sql \
  services/control-plane/migrations/012_control_nodes_tenant_scope.sql \
  services/control-plane/migrations/013_approval_payload_binding.sql \
  services/control-plane/migrations/014_tenant_scoped_approval_hash.sql \
  services/control-plane/migrations/015_approval_event_trigger.sql \
  services/control-plane/migrations/016_tenant_scoped_job_idempotency.sql \
  services/control-plane/migrations/017_active_defense_execution.sql \
  services/control-plane/migrations/018_tenant_scoped_active_defense_idempotency.sql \
  services/control-plane/migrations/019_tenant_scoped_service_requests.sql \
  services/control-plane/migrations/020_tenant_scoped_api_idempotency.sql \
  services/control-plane/migrations/021_verification_identity.sql \
  services/control-plane/migrations/022_data_request_idempotency.sql \
  services/control-plane/migrations/023_execution_attempt_integrity.sql \
  services/control-plane/migrations/024_outbound_queue_leases.sql; do
  echo "Applying $migration"
  docker exec -i "$NAME" psql -v ON_ERROR_STOP=1 -U ftn -d ftn < "$migration" >/dev/null
done

docker exec -i "$NAME" psql -v ON_ERROR_STOP=1 -U ftn -d ftn <<'SQL' >/dev/null
INSERT INTO tenants (slug,name) VALUES ('ci-restore','CI Restore Validation') ON CONFLICT (slug) DO NOTHING;
INSERT INTO tenants (slug,name) VALUES ('ci-restore-b','CI Restore Validation B') ON CONFLICT (slug) DO NOTHING;
SQL

test "$(docker exec "$NAME" psql -At -U ftn -d ftn -c "SELECT count(*) FROM tenants WHERE slug='ci-restore';")" = "1"

echo "Validating tenant-scoped service requests"
docker exec -i "$NAME" psql -v ON_ERROR_STOP=1 -U ftn -d ftn <<'SQL' >/dev/null
INSERT INTO service_requests(tenant_id,service_id,status) SELECT id,'internet','accepted' FROM tenants WHERE slug='ci-restore';
INSERT INTO service_requests(tenant_id,service_id,status) SELECT id,'internet','accepted' FROM tenants WHERE slug='ci-restore-b';
SQL
test "$(docker exec "$NAME" psql -At -U ftn -d ftn -c "SELECT count(DISTINCT tenant_id) FROM service_requests;")" = "2"

echo "Validating tenant-scoped API idempotency storage"
docker exec -i "$NAME" psql -v ON_ERROR_STOP=1 -U ftn -d ftn <<'SQL' >/dev/null
DO $$
DECLARE t1 uuid; t2 uuid; p1 uuid; p2 uuid;
BEGIN
  SELECT id INTO t1 FROM tenants WHERE slug='ci-restore';
  SELECT id INTO t2 FROM tenants WHERE slug='ci-restore-b';
  INSERT INTO principals(tenant_id,subject,kind,issuer) VALUES(t1,'ci-p1','user','ci') ON CONFLICT DO NOTHING;
  INSERT INTO principals(tenant_id,subject,kind,issuer) VALUES(t2,'ci-p2','user','ci') ON CONFLICT DO NOTHING;
  SELECT id INTO p1 FROM principals WHERE tenant_id=t1 AND subject='ci-p1';
  SELECT id INTO p2 FROM principals WHERE tenant_id=t2 AND subject='ci-p2';
  INSERT INTO idempotency_keys(key,principal_id,request_hash,response_status,response_body,expires_at) VALUES(t1::text||':'||p1::text||':shared',p1,'hash-a',202,'{}',now()+interval '1 hour');
  INSERT INTO idempotency_keys(key,principal_id,request_hash,response_status,response_body,expires_at) VALUES(t2::text||':'||p2::text||':shared',p2,'hash-b',202,'{}',now()+interval '1 hour');
END $$;
SQL
test "$(docker exec "$NAME" psql -At -U ftn -d ftn -c "SELECT count(*) FROM idempotency_keys WHERE key LIKE '%:shared';")" = "2"

echo "Validating data request idempotency"
docker exec -i "$NAME" psql -v ON_ERROR_STOP=1 -U ftn -d ftn <<'SQL' >/dev/null
DO $$
DECLARE t1 uuid; p1 uuid; a1 uuid; r1 uuid; r2 uuid;
BEGIN
  SELECT id INTO t1 FROM tenants WHERE slug='ci-restore';
  SELECT id INTO p1 FROM principals WHERE tenant_id=t1 AND subject='ci-p1';
  INSERT INTO change_approvals(tenant_id,requested_by,action,resource,request_hash,status,expires_at)
  VALUES(t1,p1,'data.export','data-governor/request','ci-data-hash','pending',now()+interval '1 hour')
  RETURNING id INTO a1;
  INSERT INTO data_requests(tenant_id,request_type,status,requested_by,approval_id,request_json,request_hash)
  VALUES(t1,'export','pending',p1,a1,'{}','ci-data-hash') RETURNING id INTO r1;
  INSERT INTO data_requests(tenant_id,request_type,status,requested_by,approval_id,request_json,request_hash)
  VALUES(t1,'export','pending',p1,a1,'{}','ci-data-hash')
  ON CONFLICT(tenant_id,request_hash) DO UPDATE SET updated_at=data_requests.updated_at
  RETURNING id INTO r2;
  IF r1 <> r2 THEN RAISE EXCEPTION 'data request replay created duplicate row'; END IF;
END $$;
SQL
test "$(docker exec "$NAME" psql -At -U ftn -d ftn -c "SELECT count(*) FROM data_requests WHERE request_hash='ci-data-hash';")" = "1"

echo "Validating verifier identity column"
test "$(docker exec "$NAME" psql -At -U ftn -d ftn -c "SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema='public' AND table_name='durable_jobs' AND column_name='verified_by');")" = "t"

echo "Validating execution-attempt integrity trigger"
test "$(docker exec "$NAME" psql -At -U ftn -d ftn -c "SELECT count(*) FROM pg_trigger WHERE tgname='durable_jobs_execution_integrity';")" = "1"

echo "Validating outbound queue lease schema"
test "$(docker exec "$NAME" psql -At -U ftn -d ftn -c "SELECT count(*) FROM information_schema.columns WHERE table_schema='public' AND table_name='ftn_mail_outbound_queue' AND column_name IN ('lease_token','lease_expires_at');")" = "2"
test "$(docker exec "$NAME" psql -At -U ftn -d ftn -c "SELECT EXISTS (SELECT 1 FROM pg_indexes WHERE schemaname='public' AND indexname='idx_ftn_mail_outbound_processing_lease');")" = "t"

echo "Creating custom-format backup"
docker exec "$NAME" pg_dump -U ftn -d ftn -Fc > "$DUMP"
test -s "$DUMP"
docker exec "$NAME" createdb -U ftn ftn_restore
docker exec -i "$NAME" pg_restore -U ftn -d ftn_restore --exit-on-error < "$DUMP"
test "$(docker exec "$NAME" psql -At -U ftn -d ftn_restore -c "SELECT count(*) FROM tenants WHERE slug='ci-restore';")" = "1"

test "$(docker exec "$NAME" psql -At -U ftn -d ftn -c "SELECT count(*) FROM pg_trigger WHERE tgname='durable_jobs_event_journal';")" = "1"
test "$(docker exec "$NAME" psql -At -U ftn -d ftn -c "SELECT count(*) FROM pg_trigger WHERE tgname='change_approvals_lifecycle_event';")" = "1"
test "$(docker exec "$NAME" psql -At -U ftn -d ftn -c "SELECT count(*) FROM schema_migrations WHERE version >= 11 AND version <= 24;")" = "14"
test "$(docker exec "$NAME" psql -At -U ftn -d ftn -c "SELECT EXISTS (SELECT 1 FROM pg_indexes WHERE schemaname='public' AND indexname='service_requests_tenant_idx');")" = "t"
test "$(docker exec "$NAME" psql -At -U ftn -d ftn -c "SELECT EXISTS (SELECT 1 FROM pg_indexes WHERE schemaname='public' AND indexname='data_requests_tenant_request_hash_uidx');")" = "t"

echo "PostgreSQL migration + backup + restore validation passed."
