package agent

import (
	"context"
	"errors"
	"testing"
)

type testRuntime struct {
	response Response
	err      error
	seen     Request
}

func (r *testRuntime) Run(_ context.Context, request Request) (Response, error) {
	r.seen = request
	return r.response, r.err
}

type testPolicy struct {
	err  error
	seen Request
}

func (p *testPolicy) Allow(_ context.Context, request Request, _ Response) error {
	p.seen = request
	return p.err
}

func TestFleetRequiresRuntimeAndPolicy(t *testing.T) {
	if _, err := NewFleet(nil, &testPolicy{}); err == nil {
		t.Fatal("expected missing runtime error")
	}
	if _, err := NewFleet(&testRuntime{}, nil); err == nil {
		t.Fatal("expected missing policy error")
	}
}

func TestFleetAppliesDefaultModeAndPolicy(t *testing.T) {
	runtime := &testRuntime{response: Response{Text: "ok"}}
	policy := &testPolicy{}
	fleet, err := NewFleet(runtime, policy)
	if err != nil {
		t.Fatalf("NewFleet: %v", err)
	}

	request := Request{Scope: Scope{ServiceID: "ftn-service"}, Input: "status"}
	response, err := fleet.Handle(context.Background(), request)
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}
	if response.Text != "ok" {
		t.Fatalf("unexpected response: %#v", response)
	}
	if runtime.seen.Mode != ModeAssistant {
		t.Fatalf("expected default assistant mode, got %q", runtime.seen.Mode)
	}
	if policy.seen.Scope.ServiceID != "ftn-service" {
		t.Fatalf("policy saw wrong service: %#v", policy.seen.Scope)
	}
}

func TestFleetRejectsMissingScope(t *testing.T) {
	fleet, err := NewFleet(&testRuntime{}, &testPolicy{})
	if err != nil {
		t.Fatalf("NewFleet: %v", err)
	}
	if _, err := fleet.Handle(context.Background(), Request{Input: "status"}); err == nil {
		t.Fatal("expected scope validation error")
	}
}

func TestFleetPropagatesRuntimeAndPolicyErrors(t *testing.T) {
	runtimeErr := errors.New("runtime failed")
	fleet, err := NewFleet(&testRuntime{err: runtimeErr}, &testPolicy{})
	if err != nil {
		t.Fatalf("NewFleet: %v", err)
	}
	if _, err := fleet.Handle(context.Background(), Request{Scope: Scope{UserID: "user"}}); err == nil {
		t.Fatal("expected runtime error")
	}

	policyErr := errors.New("policy denied")
	fleet, err = NewFleet(&testRuntime{response: Response{NeedsApproval: true}}, &testPolicy{err: policyErr})
	if err != nil {
		t.Fatalf("NewFleet: %v", err)
	}
	if _, err := fleet.Handle(context.Background(), Request{Scope: Scope{UserID: "user"}}); err == nil {
		t.Fatal("expected policy error")
	}
}
