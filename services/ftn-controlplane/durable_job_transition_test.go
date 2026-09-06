package controlplane

import "testing"

func TestTransitionJobRejectsTerminalMutation(t *testing.T) {
 job:=DurableJob{ID:"j1",State:JobSucceeded}
 if _,err:=TransitionJob(job,JobRunning,nowUTC()); err!=ErrInvalidJobState {t.Fatalf("got %v",err)}
}
