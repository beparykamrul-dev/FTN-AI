-- FTN verification/rollback race boundary. Additive/idempotent.
-- Serialize terminal verification and rollback/execute approval transitions per job.

CREATE OR REPLACE FUNCTION ftn_lock_job_verification_boundary()
RETURNS trigger
LANGUAGE plpgsql
AS $$
DECLARE
    job_approval UUID;
    job_status TEXT;
    approval_status TEXT;
BEGIN
    IF NEW.job_id IS NULL THEN
        RAISE EXCEPTION 'verification_job_required';
    END IF;

    PERFORM pg_advisory_xact_lock(hashtextextended(NEW.job_id::text, 0));

    SELECT approval_id, status
      INTO job_approval, job_status
      FROM durable_jobs
     WHERE id = NEW.job_id
     FOR SHARE;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'verification_job_not_found';
    END IF;

    -- Verification cannot be attached to an execution that has already moved
    -- to a terminal state unless the caller is updating the same verification.
    IF TG_OP = 'INSERT' AND job_status IN ('succeeded','failed','cancelled') THEN
        RAISE EXCEPTION 'verification_after_terminal_job_forbidden';
    END IF;

    IF job_approval IS NOT NULL THEN
        SELECT status INTO approval_status
          FROM change_approvals
         WHERE id = job_approval
         FOR SHARE;
        IF NOT FOUND THEN
            RAISE EXCEPTION 'verification_approval_not_found';
        END IF;
        IF approval_status NOT IN ('approved','executed','rolled_back') THEN
            RAISE EXCEPTION 'verification_approval_state_invalid';
        END IF;
    END IF;

    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS verification_job_boundary ON job_verifications;
CREATE TRIGGER verification_job_boundary
BEFORE INSERT OR UPDATE OF job_id
ON job_verifications
FOR EACH ROW
EXECUTE FUNCTION ftn_lock_job_verification_boundary();

-- Approval transitions that finalize execution/rollback are serialized with
-- the same job lock when their resource identifies a durable job.
CREATE OR REPLACE FUNCTION ftn_lock_job_approval_transition()
RETURNS trigger
LANGUAGE plpgsql
AS $$
DECLARE
    job_id_text TEXT;
    job_status TEXT;
BEGIN
    IF NEW.status IS NOT DISTINCT FROM OLD.status THEN
        RETURN NEW;
    END IF;

    IF OLD.status NOT IN ('approved','executed','rolled_back') THEN
        RETURN NEW;
    END IF;

    IF NEW.status IN ('executed','rolled_back') THEN
        IF OLD.resource LIKE 'durable_job:%' THEN
            job_id_text := substring(OLD.resource from length('durable_job:') + 1);
            IF NULLIF(BTRIM(job_id_text), '') IS NULL THEN
                RAISE EXCEPTION 'approval_job_resource_invalid';
            END IF;
            PERFORM pg_advisory_xact_lock(hashtextextended(job_id_text, 0));
            SELECT status INTO job_status FROM durable_jobs WHERE id::text = job_id_text FOR SHARE;
            IF NOT FOUND THEN
                RAISE EXCEPTION 'approval_job_not_found';
            END IF;
            IF NEW.status = 'executed' AND job_status <> 'succeeded' THEN
                RAISE EXCEPTION 'approval_execute_requires_successful_job';
            END IF;
            IF NEW.status = 'rolled_back' AND job_status NOT IN ('succeeded','failed','cancelled') THEN
                RAISE EXCEPTION 'approval_rollback_requires_terminal_job';
            END IF;
        END IF;
    END IF;

    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS approval_job_transition_boundary ON change_approvals;
CREATE TRIGGER approval_job_transition_boundary
BEFORE UPDATE OF status
ON change_approvals
FOR EACH ROW
EXECUTE FUNCTION ftn_lock_job_approval_transition();
