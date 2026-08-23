package kernel

import "testing"

func TestRequestIsStructured(t *testing.T) {
	r := Request{ID: "req-1", Tool: "network.inventory", Capability: "read", Arguments: map[string]string{"device": "edge-1"}}
	if r.Tool == "" || r.Capability == "" || r.Arguments == nil {
		t.Fatal("request must contain a structured tool contract")
	}
}

func TestMutationCarriesApprovalAndIdempotency(t *testing.T) {
	r := Request{Mutating: true, ApprovalID: "approval-1", Idempotency: "idem-1"}
	if r.Mutating && (r.ApprovalID == "" || r.Idempotency == "") {
		t.Fatal("mutating request requires approval and idempotency identifiers")
	}
}
