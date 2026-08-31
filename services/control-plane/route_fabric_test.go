package main

import "testing"

func TestNodeFabricState(t *testing.T) {
    n := Node{BGPUp:true, BFDUp:true, ISISUp:true, EVPNReady:true, AnycastReady:true, RPKIValid:true, PrefixCount:1200, CapacityMbps:10000, UtilizationPercent:40}
    s := nodeFabricState(n)
    if !s.BGPUp || !s.BFDUp || !s.ISISUp || !s.EVPNReady || !s.AnycastReady || !s.RPKIValid { t.Fatal("fabric readiness was not mapped") }
    if s.PrefixCount != 1200 || s.CapacityMbps != 10000 || s.UtilizationPercent != 40 { t.Fatal("fabric metrics were not mapped") }
}

func TestValidNodeRejectsInvalidFabricMetrics(t *testing.T) {
    n := Node{ID:"n1",Provider:"ftn",Services:[]string{"exrouter"},CapacityMbps:100,UtilizationPercent:101}
    if validNode(n) { t.Fatal("utilization over 100 must be rejected") }
    n.UtilizationPercent=50; n.PrefixCount=-1
    if validNode(n) { t.Fatal("negative prefix count must be rejected") }
}
