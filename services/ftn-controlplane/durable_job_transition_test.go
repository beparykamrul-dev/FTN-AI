package controlplane

import (
 "testing"
 "time"
)

func TestTransitionJobRejectsTerminalMutation(t *testing.T) {
 job:=DurableJob{ID:"j1",State:JobSucceeded}
 if _,err:=TransitionJob(job,JobRunning,time.Now().UTC()); err!=ErrInvalidJobState {t.Fatalf("got %v",err)}
}
