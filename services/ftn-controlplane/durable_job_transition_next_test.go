package controlplane

import (
 "testing"
 "time"
)

func TestTransitionJobIncrementsAttemptOnlyWhenRunning(t *testing.T) {
 now:=time.Unix(10,0); j:=DurableJob{ID:"j",State:JobPending,Attempt:2}; j2,err:=TransitionJob(j,JobRunning,now); if err!=nil {t.Fatal(err)}; if j2.Attempt!=3 {t.Fatalf("attempt=%d",j2.Attempt)}
 j3,err:=TransitionJob(j2,JobSucceeded,now.Add(time.Second)); if err!=nil {t.Fatal(err)}; if j3.Attempt!=3 {t.Fatalf("attempt=%d",j3.Attempt)}
}
