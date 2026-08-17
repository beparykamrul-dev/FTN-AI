package controlplane

import (
	"testing"
	"time"
)

func TestDurableJobLifecycleAndCheckpoint(t *testing.T) {
	store := NewMemoryJobStore()
	now := time.Unix(100, 0)
	job := DurableJob{ID: "job-1", TenantID: "tenant-1", IdempotencyKey: "idem-1", State: JobPending, CreatedAt: now, UpdatedAt: now}
	if err := store.Create(job); err != nil {
		t.Fatal(err)
	}
	if err := store.Create(job); err != ErrDuplicateJob {
		t.Fatalf("expected duplicate idempotency rejection, got %v", err)
	}
	job, _ = store.Get(job.ID)
	job, err := TransitionJob(job, JobRunning, now.Add(time.Second))
	if err != nil || job.Attempt != 1 {
		t.Fatalf("running transition failed: %+v %v", job, err)
	}
	job = CheckpointJob(job, "step-2", now.Add(2*time.Second))
	if err := store.Update(job); err != nil {
		t.Fatal(err)
	}
	job, _ = store.Get(job.ID)
	if job.Checkpoint != "step-2" {
		t.Fatalf("checkpoint not persisted: %q", job.Checkpoint)
	}
	job, err = TransitionJob(job, JobSucceeded, now.Add(3*time.Second))
	if err != nil || job.State != JobSucceeded {
		t.Fatalf("success transition failed: %+v %v", job, err)
	}
}

func TestDurableJobRejectsInvalidTerminalTransition(t *testing.T) {
	job := DurableJob{ID: "job-2", State: JobSucceeded}
	if _, err := TransitionJob(job, JobRunning, time.Now()); err != ErrInvalidJobState {
		t.Fatalf("expected invalid transition, got %v", err)
	}
}
