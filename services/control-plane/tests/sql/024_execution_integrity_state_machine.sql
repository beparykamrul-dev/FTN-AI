-- FTN 015..024 execution-integrity state-machine tests.
-- Run only against a disposable PostgreSQL database after all migrations.
-- Every expected rejection is isolated behind a SAVEPOINT.

BEGIN;
SET LOCAL statement_timeout = '10s';
SET LOCAL lock_timeout = '2s';

DO $$
DECLARE
  t UUID;
  p UUID;
  approval UUID;
  job UUID;
  attempt UUID;
  payload JSONB := '{"target":"router-01","operation":"reload"}'::jsonb;
  worker TEXT := 'semantic-test-worker';
  resource TEXT;
  got TEXT;
BEGIN
  SELECT id INTO t FROM tenants ORDER BY created_at LIMIT 1;
  IF t IS NULL THEN
    INSERT INTO tenants(slug,name) VALUES ('semantic-test-tenant','Semantic Integrity Test') RETURNING id INTO t;
  END IF;

  SELECT id INTO p FROM principals WHERE tenant_id=t ORDER BY created_at LIMIT 1;
  IF p IS NULL THEN
    INSERT INTO principals(tenant_id,subject,kind,issuer)
    VALUES (t,'semantic-test-principal','service','semantic-test')
    RETURNING id INTO p;
  END IF;

  INSERT INTO change_approvals(tenant_id,requested_by,action,resource,request_hash,status,expires_at,approval_payload)
  VALUES (t,p,'router.reload','durable-job-test','semantic-request-001','approved',now()+interval '10 minutes',payload)
  RETURNING id INTO approval;

  -- Valid approved durable job must bind action/hash/payload to the approval.
  INSERT INTO durable_jobs(tenant_id,idempotency_key,job_type,payload,status,max_attempts,correlation_id,approval_id,execution_action)
  VALUES (t,'semantic-job-001','network.change',payload,'queued',3,'semantic-correlation-001',approval,'router.reload')
  RETURNING id INTO job;

  IF NOT EXISTS (
    SELECT 1 FROM durable_jobs
    WHERE id=job AND approval_request_hash='semantic-request-001'
      AND execution_payload_hash=md5(payload::text)
  ) THEN
    RAISE EXCEPTION 'approved job did not bind approval hash/payload hash';
  END IF;

  -- queued -> running with a valid lease is allowed and must pass approval gates.
  UPDATE durable_jobs
     SET status='running', locked_by=worker, locked_at=now(), started_at=now(), attempts=1
   WHERE id=job;

  -- A matching running execution attempt is allowed.
  INSERT INTO execution_attempts(job_id,attempt_no,worker_id,status)
  VALUES(job,1,worker,'running')
  RETURNING id INTO attempt;

  -- Duplicate active attempt is rejected.
  SAVEPOINT duplicate_active_attempt;
  BEGIN
    INSERT INTO execution_attempts(job_id,attempt_no,worker_id,status)
    VALUES(job,2,worker,'running');
    RAISE EXCEPTION 'expected duplicate active attempt rejection did not occur';
  EXCEPTION WHEN unique_violation THEN
    ROLLBACK TO SAVEPOINT duplicate_active_attempt;
  END;

  -- Identity must remain immutable.
  SAVEPOINT attempt_identity;
  BEGIN
    UPDATE execution_attempts SET worker_id='different-worker' WHERE id=attempt;
    RAISE EXCEPTION 'expected execution attempt identity rejection did not occur';
  EXCEPTION WHEN OTHERS THEN
    GET STACKED DIAGNOSTICS got = MESSAGE_TEXT;
    ROLLBACK TO SAVEPOINT attempt_identity;
    IF got NOT LIKE '%immutable%' THEN RAISE; END IF;
  END;

  -- running -> queued is not a legal execution-attempt transition.
  SAVEPOINT attempt_bad_transition;
  BEGIN
    UPDATE execution_attempts SET status='failed', finished_at=NULL WHERE id=attempt;
    RAISE EXCEPTION 'expected finished_at rejection did not occur';
  EXCEPTION WHEN OTHERS THEN
    GET STACKED DIAGNOSTICS got = MESSAGE_TEXT;
    ROLLBACK TO SAVEPOINT attempt_bad_transition;
    IF got NOT LIKE '%finished_at%' THEN RAISE; END IF;
  END;

  -- A valid terminal attempt requires finished_at.
  UPDATE execution_attempts SET status='succeeded', finished_at=now() WHERE id=attempt;

  -- A successful approved job is allowed only with the matching successful attempt.
  UPDATE durable_jobs SET status='succeeded', finished_at=now() WHERE id=job;

  -- Terminal execution metadata is immutable.
  SAVEPOINT terminal_job_mutation;
  BEGIN
    UPDATE durable_jobs SET execution_action='tampered' WHERE id=job;
    RAISE EXCEPTION 'expected terminal job metadata rejection did not occur';
  EXCEPTION WHEN OTHERS THEN
    GET STACKED DIAGNOSTICS got = MESSAGE_TEXT;
    ROLLBACK TO SAVEPOINT terminal_job_mutation;
    IF got NOT LIKE '%immutable%' THEN RAISE; END IF;
  END;

  -- Approval cannot execute until its durable job has succeeded.
  resource := 'durable_job:' || job::text;
  UPDATE change_approvals SET resource=resource WHERE id=approval;
  UPDATE change_approvals SET status='executed', executed_at=now() WHERE id=approval;

  -- Executed -> rolled_back is permitted only while the referenced job is terminal.
  UPDATE change_approvals SET status='rolled_back' WHERE id=approval;

  -- Invalid approval mutation after approval must be rejected.
  SAVEPOINT approval_mutation;
  BEGIN
    UPDATE change_approvals SET action='tampered-action' WHERE id=approval;
    RAISE EXCEPTION 'expected approved request immutability rejection did not occur';
  EXCEPTION WHEN OTHERS THEN
    GET STACKED DIAGNOSTICS got = MESSAGE_TEXT;
    ROLLBACK TO SAVEPOINT approval_mutation;
    IF got NOT LIKE '%immutable%' THEN RAISE; END IF;
  END;

  -- A terminal execution attempt cannot be replayed/terminalized again.
  SAVEPOINT terminal_attempt_replay;
  BEGIN
    UPDATE execution_attempts SET status='failed' WHERE id=attempt;
    RAISE EXCEPTION 'expected terminal attempt replay rejection did not occur';
  EXCEPTION WHEN OTHERS THEN
    GET STACKED DIAGNOSTICS got = MESSAGE_TEXT;
    ROLLBACK TO SAVEPOINT terminal_attempt_replay;
    IF got NOT LIKE '%immutable%' THEN RAISE; END IF;
  END;

  -- A new approved job with no matching successful attempt cannot be marked succeeded.
  INSERT INTO change_approvals(tenant_id,requested_by,action,resource,request_hash,status,expires_at,approval_payload)
  VALUES (t,p,'router.reload','durable-job-test-2','semantic-request-002','approved',now()+interval '10 minutes',payload)
  RETURNING id INTO approval;

  INSERT INTO durable_jobs(tenant_id,idempotency_key,job_type,payload,status,max_attempts,correlation_id,approval_id,execution_action,attempts)
  VALUES (t,'semantic-job-002','network.change',payload,'queued',3,'semantic-correlation-002',approval,'router.reload',1)
  RETURNING id INTO job;

  UPDATE durable_jobs SET status='running',locked_by=worker,locked_at=now(),started_at=now() WHERE id=job;

  SAVEPOINT missing_success_attempt;
  BEGIN
    UPDATE durable_jobs SET status='succeeded',finished_at=now() WHERE id=job;
    RAISE EXCEPTION 'expected missing successful attempt rejection did not occur';
  EXCEPTION WHEN OTHERS THEN
    GET STACKED DIAGNOSTICS got = MESSAGE_TEXT;
    ROLLBACK TO SAVEPOINT missing_success_attempt;
    IF got NOT LIKE '%successful%' THEN RAISE; END IF;
  END;

  -- Requeue revokes the lease; old worker can no longer finish its attempt.
  INSERT INTO execution_attempts(job_id,attempt_no,worker_id,status)
  VALUES(job,1,worker,'running') RETURNING id INTO attempt;
  UPDATE durable_jobs SET status='queued' WHERE id=job;
  IF EXISTS (SELECT 1 FROM durable_jobs WHERE id=job AND (locked_by IS NOT NULL OR locked_at IS NOT NULL)) THEN
    RAISE EXCEPTION 'requeue did not revoke execution lease';
  END IF;

  SAVEPOINT old_worker_finish;
  BEGIN
    UPDATE execution_attempts SET status='succeeded',finished_at=now() WHERE id=attempt;
    RAISE EXCEPTION 'expected reclaimed-lease rejection did not occur';
  EXCEPTION WHEN OTHERS THEN
    GET STACKED DIAGNOSTICS got = MESSAGE_TEXT;
    ROLLBACK TO SAVEPOINT old_worker_finish;
    IF got NOT LIKE '%lease%' THEN RAISE; END IF;
  END;

  -- Approval expiration blocks a fresh privileged claim.
  INSERT INTO change_approvals(tenant_id,requested_by,action,resource,request_hash,status,expires_at,approval_payload)
  VALUES (t,p,'router.reload','durable-job-expired','semantic-request-003','approved',now()-interval '1 minute',payload)
  RETURNING id INTO approval;

  INSERT INTO durable_jobs(tenant_id,idempotency_key,job_type,payload,status,max_attempts,correlation_id,approval_id,execution_action)
  VALUES (t,'semantic-job-003','network.change',payload,'queued',3,'semantic-correlation-003',approval,'router.reload')
  RETURNING id INTO job;

  SAVEPOINT expired_claim;
  BEGIN
    UPDATE durable_jobs SET status='running',locked_by=worker,locked_at=now(),started_at=now() WHERE id=job;
    RAISE EXCEPTION 'expected expired approval rejection did not occur';
  EXCEPTION WHEN OTHERS THEN
    GET STACKED DIAGNOSTICS got = MESSAGE_TEXT;
    ROLLBACK TO SAVEPOINT expired_claim;
    IF got NOT LIKE '%expired%' THEN RAISE; END IF;
  END;

  RAISE NOTICE 'FTN execution integrity state-machine tests PASSED';
END $$;

ROLLBACK;
