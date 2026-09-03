-- FTN verification/rollback race boundary. Additive/idempotent.
-- Serialize terminal approval execute/rollback transitions per durable job.
-- Verification state is persisted on durable_jobs.verification_payload in the
-- current execution model; no separate job_verifications table is assumed.

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

    IF NEW.status IN ('executed','rolled_back') AND OLD.resource LIKE 'durable_job:%' THEN
        job_id_text := substring(OLD.resource from length('durable_job:') + 1);
        IF NULLIF(BTRIM(job_id_text), '') IS NULL THEN
            RAISE EXCEPTION 'approval_job_resource_invalid';
        END IF;
        PERFORM pg_advisory_xact_lock(hashtextextended(job_id_text, 0));
        SELECT status INTO job_status
          FROM durable_jobs
         WHERE id::text = job_id_text
         FOR SHARE;
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

    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS approval_job_transition_boundary ON change_approvals;
CREATE TRIGGER approval_job_transition_boundary
BEFORE UPDATE OF status
ON change_approvals
FOR EACH ROW
EXECUTE FUNCTION ftn_lock_job_approval_transition();
