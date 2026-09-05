package mesh

import("testing";"time")

func TestSelectPathsKeepsBestDuplicateNextHop(t *testing.T){got:=SelectPaths([]Path{{NextHop:"edge-a",Metric:20,Hops:2},{NextHop:"edge-a",Metric:10,Hops:1},{NextHop:"edge-b",Metric:30,Hops:3}},2);if len(got)!=2||got[0].NextHop!="edge-a"||got[0].Metric!=10{t.Fatalf("unexpected paths: %#v",got)}}
func TestChooseFailoverRejectsFutureLink(t *testing.T){now:=time.Now().UTC();_,ok:=ChooseFailover("edge-a",[]Link{{ID:"edge-b",State:LinkUp,Metric:1,UpdatedAt:now.Add(time.Minute)}},now,5*time.Minute);if ok{t.Fatal("future-dated link must not be selected")}}
