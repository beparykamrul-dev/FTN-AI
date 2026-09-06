package mesh

import "testing"

func TestSelectPathsKeepsBestDuplicateNextHop(t *testing.T) {
	got := SelectPaths([]Path{{NextHop: "n1", Metric: 20, Hops: 2}, {NextHop: "n1", Metric: 5, Hops: 1}}, 2)
	if len(got) != 1 || got[0].Metric != 5 { t.Fatalf("best duplicate path was not retained: %#v", got) }
}

func TestSelectPathsRejectsInvalidHopCount(t *testing.T) {
	if got := SelectPaths([]Path{{NextHop: "n1", Hops: -1}}, 1); len(got) != 0 { t.Fatal("negative hop count must be rejected") }
}
