package main

import (
	"context"
	"testing"
)

type fakeCoreFailoverAdapter struct { switched, verified, rolledBack string }
func (f *fakeCoreFailoverAdapter) SwitchActiveCore(_ context.Context, node string) error { f.switched=node; return nil }
func (f *fakeCoreFailoverAdapter) VerifyActiveCore(_ context.Context, node string) error { f.verified=node; return nil }
func (f *fakeCoreFailoverAdapter) RollbackActiveCore(_ context.Context, node string) error { f.rolledBack=node; return nil }

func TestFTNFailoverExecutorRequiresAdapter(t *testing.T) {
	e := NewFTNFailoverExecutor(nil)
	payload := FTNFailoverJobPayload{Intent:FTNFailoverIntent{TargetNode:"core-b",DecisionHash:"hash"},PrechangeSnapshotRequired:true,VerificationRequired:true}
	if err := e.Execute(context.Background(), payload); err == nil { t.Fatal("expected adapter requirement") }
}

func TestFTNFailoverExecutorLifecycle(t *testing.T) {
	a := &fakeCoreFailoverAdapter{}
	e := NewFTNFailoverExecutor(a)
	payload := FTNFailoverJobPayload{Intent:FTNFailoverIntent{TargetNode:"core-b",DecisionHash:"hash"},PrechangeSnapshotRequired:true,VerificationRequired:true,RollbackWhenSafe:true}
	ctx := context.Background()
	if err := e.Execute(ctx, payload); err != nil { t.Fatal(err) }
	if err := e.Verify(ctx, payload); err != nil { t.Fatal(err) }
	if err := e.Rollback(ctx, payload); err != nil { t.Fatal(err) }
	if a.switched!="core-b" || a.verified!="core-b" || a.rolledBack!="core-b" { t.Fatalf("adapter=%+v", a) }
}
