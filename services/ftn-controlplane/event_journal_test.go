package controlplane

import "testing"

func TestMemoryEventJournalSequenceReplayAndOffsets(t *testing.T) {
	j := NewMemoryEventJournal()
	first, err := j.Append(JournalEvent{ID: "e1", TenantID: "t1", Type: "job.created", CorrelationID: "c1", Payload: []byte(`{"job":"j1"}`)})
	if err != nil || first.Sequence != 1 {
		t.Fatalf("first event: %+v %v", first, err)
	}
	second, err := j.Append(JournalEvent{ID: "e2", TenantID: "t1", Type: "job.running", CorrelationID: "c1", CausationID: "e1"})
	if err != nil || second.Sequence != 2 {
		t.Fatalf("second event: %+v %v", second, err)
	}
	other, err := j.Append(JournalEvent{ID: "e3", TenantID: "t2", Type: "job.created"})
	if err != nil || other.Sequence != 1 {
		t.Fatalf("tenant isolation sequence: %+v %v", other, err)
	}
	items, err := j.ReadAfter("t1", 1, 10)
	if err != nil || len(items) != 1 || items[0].ID != "e2" {
		t.Fatalf("replay: %+v %v", items, err)
	}
	if err := j.CommitOffset("worker-1", "t1", 2); err != nil {
		t.Fatal(err)
	}
	if got := j.Offset("worker-1", "t1"); got != 2 {
		t.Fatalf("offset = %d, want 2", got)
	}
	if err := j.CommitOffset("worker-1", "t1", 1); err == nil {
		t.Fatal("expected backwards offset rejection")
	}
	if err := j.CommitOffset("worker-1", "t1", 99); err == nil {
		t.Fatal("expected future offset rejection")
	}
}
