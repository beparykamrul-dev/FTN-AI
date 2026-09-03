package main

import (
	"encoding/json"
	"testing"
)

func TestBuildFTNFailoverJobPayload(t *testing.T) {
	payload, err := BuildFTNFailoverJobPayload(FTNFailoverIntent{TargetNode:"core-b", DecisionHash:"abc", Reason:"alternate_core_healthy"})
	if err != nil { t.Fatal(err) }
	var got FTNFailoverJobPayload
	if err := json.Unmarshal(payload, &got); err != nil { t.Fatal(err) }
	if got.Intent.TargetNode != "core-b" || !got.PrechangeSnapshotRequired || !got.VerificationRequired || !got.RollbackWhenSafe {
		t.Fatalf("unexpected payload: %+v", got)
	}
}

func TestBuildFTNFailoverJobPayloadRejectsUnvalidatedIntent(t *testing.T) {
	if _, err := BuildFTNFailoverJobPayload(FTNFailoverIntent{TargetNode:"core-b"}); err == nil {
		t.Fatal("expected validation error")
	}
}
