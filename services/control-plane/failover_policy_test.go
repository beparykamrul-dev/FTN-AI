package main

import (
    "testing"
    "time"
)

func TestChooseFailoverPath(t *testing.T) {
    now := time.Now()
    paths := []FailoverPath{
        {PathID:"primary", Priority:1, Healthy:false, BGPUp:true, BFDUp:true, LastSeen:now},
        {PathID:"stale", Priority:2, Healthy:true, BGPUp:true, BFDUp:true, LastSeen:now.Add(-2*time.Minute)},
        {PathID:"backup", Priority:3, Healthy:true, BGPUp:true, BFDUp:true, LastSeen:now},
    }
    got, ok := chooseFailoverPath(paths, FailoverPolicy{RequireBGP:true, RequireBFD:true}, now)
    if !ok || got.PathID != "backup" { t.Fatalf("got %+v ok=%v", got, ok) }
}
