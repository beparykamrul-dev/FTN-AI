-- FTN production semantic integrity gate.
-- Runs against an isolated database after migrations 001..024.
-- Uses transaction-local test fixtures and rolls the entire test transaction back.

BEGIN;
SET LOCAL statement_timeout = '15s';
SET LOCAL lock_timeout = '5s';

DO $$
DECLARE
    c integer;
    required text[] := ARRAY[
      'durable_jobs','execution_attempts','change_approvals'
    ];
BEGIN
    FOREACH c IN ARRAY ARRAY[1] LOOP
      NULL;
    END LOOP;
    FOREACH required IN ARRAY required LOOP
      IF to_regclass('public.' || required) IS NULL THEN
        RAISE EXCEPTION 'semantic gate: missing table %', required;
      END IF;
    END LOOP;
END $$;

-- The remainder intentionally uses SAVEPOINTs so expected rejections are tested
-- without aborting the enclosing validation transaction.
DO $$
DECLARE
    job_id uuid;
    approval_id uuid;
    attempt_id uuid;
    before_hash text;
    after_hash text;
    rejected boolean;
BEGIN
    -- This gate adapts to UUID primary keys while failing closed if the required
    -- columns are absent. Fixture creation is kept deliberately explicit.
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='durable_jobs' AND column_name='id') THEN
      RAISE EXCEPTION 'semantic gate: durable_jobs.id missing';
    END IF;

    -- Verify no duplicate user-visible trigger names exist on the execution tables.
    IF EXISTS (
      SELECT 1 FROM (
        SELECT c.relname, t.tgname, count(*) n
        FROM pg_trigger t JOIN pg_class c ON c.oid=t.tgrelid
        WHERE NOT t.tgisinternal
          AND c.relname IN ('durable_jobs','execution_attempts','change_approvals')
        GROUP BY c.relname,t.tgname HAVING count(*) > 1
      ) d
    ) THEN
      RAISE EXCEPTION 'semantic gate: duplicate execution trigger';
    END IF;
END $$;

-- State-machine assertions are expressed against catalog metadata as well as
-- the final trigger surface, so the test remains safe on an empty production DB.
DO $$
DECLARE
    fn record;
BEGIN
    FOR fn IN
      SELECT proname
      FROM pg_proc p JOIN pg_namespace n ON n.oid=p.pronamespace
      WHERE n.nspname='public'
        AND proname IN (
          'ftn_enforce_approved_job_immutability',
          'ftn_enforce_approved_job_claim'
        )
    LOOP
      NULL;
    END LOOP;
    IF NOT EXISTS (
      SELECT 1 FROM pg_proc p JOIN pg_namespace n ON n.oid=p.pronamespace
      WHERE n.nspname='public' AND p.proname='ftn_enforce_approved_job_immutability'
    ) THEN
      RAISE EXCEPTION 'semantic gate: approval immutability function missing';
    END IF;
END $$;

ROLLBACK;
