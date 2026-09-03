-- FTN execution lease integrity boundary. Additive/idempotent.
-- A reclaimed/expired lease must invalidate the previous worker's authority.

CREATE OR REPLACE FUNCTION ftn_validate_running_job_lease()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    IF NEW.status = 'running' THEN
        IF NULLIF(BTRIM(NEW.locked_by), '') IS NULL OR NEW.locked_at IS NULL THEN
            RAISE EXCEPTION 'running_job_requires_active_lease';
        END IF;
    END IF;

    IF NEW.status IN ('queued','succeeded','failed','cancelled') THEN
        IF NEW.locked_by IS NOT NULL AND NULLIF(BTRIM(NEW.locked_by), '') IS NOT NULL
           AND NEW.status IN ('succeeded','failed') THEN
            -- Terminal completion retains the worker identity for audit, but the
            -- lease timestamp must still identify the execution that completed it.
            IF NEW.locked_at IS NULL THEN
                RAISE EXCEPTION 'terminal_job_requires_lease_provenance';
            END IF;
        END IF;
    END IF;

    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS durable_jobs_lease_integrity ON durable_jobs;
CREATE TRIGGER durable_jobs_lease_integrity
BEFORE INSERT OR UPDATE OF status, locked_by, locked_at
ON durable_jobs
FOR EACH ROW
EXECUTE FUNCTION ftn_validate_running_job_lease();

-- An execution attempt can only finish while its job is still owned by the same
-- worker. Once the lease has been reclaimed, the old worker loses authority.
CREATE OR REPLACE FUNCTION ftn_validate_attempt_lease_owner()
RETURNS trigger
LANGUAGE plpgsql
AS $$
DECLARE
    job_status TEXT;
    job_worker TEXT;
BEGIN
    IF NEW.status NOT IN ('succeeded','failed') OR OLD.status <> 'running' THEN
        RETURN NEW;
    END IF;

    SELECT status, COALESCE(locked_by, '')
      INTO job_status, job_worker
      FROM durable_jobs
     WHERE id = NEW.job_id
     FOR SHARE;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'durable_job_not_found';
    END IF;
    IF job_status IS DISTINCT FROM 'running' THEN
        RAISE EXCEPTION 'execution_attempt_lease_no_longer_active';
    END IF;
    IF job_worker IS DISTINCT FROM NEW.worker_id THEN
        RAISE EXCEPTION 'execution_attempt_lease_owner_mismatch';
    END IF;

    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS execution_attempt_lease_owner ON execution_attempts;
CREATE TRIGGER execution_attempt_lease_owner
BEFORE UPDATE OF status
ON execution_attempts
FOR EACH ROW
EXECUTE FUNCTION ftn_validate_attempt_lease_owner();

-- Requeueing is a lease revocation boundary: the old worker must no longer be
-- represented as the active owner of the queued job.
CREATE OR REPLACE FUNCTION ftn_requeue_revokes_execution_lease()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    IF OLD.status = 'running' AND NEW.status = 'queued' THEN
        NEW.locked_by := NULL;
        NEW.locked_at := NULL;
        NEW.started_at := NULL;
    END IF;
    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS durable_jobs_requeue_lease_revoke ON durable_jobs;
CREATE TRIGGER durable_jobs_requeue_lease_revoke
BEFORE UPDATE OF status
ON durable_jobs
FOR EACH ROW
EXECUTE FUNCTION ftn_requeue_revokes_execution_lease();

CREATE INDEX IF NOT EXISTS durable_jobs_lease_idx
    ON durable_jobs(status, locked_at, locked_by);
