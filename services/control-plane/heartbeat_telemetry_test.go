package main

import "testing"

func TestAdvancedTelemetryValidation(t *testing.T) {
    good := Node{ID:"n1", Provider:"ftn", Services:[]string{"exrouter"}, CPUPercent:20, RAMPercent:30, SSDPercent:40, HDDPercent:50, NetMbps:100, LatencyMs:5, PacketLoss:0.1, JitterMs:1.5, Retransmissions:2}
    if !validNode(good) { t.Fatal("expected advanced telemetry snapshot to be valid") }
    bad := good
    bad.JitterMs = -1
    if validNode(bad) { t.Fatal("negative jitter must be rejected") }
    bad = good
    bad.Retransmissions = -1
    if validNode(bad) { t.Fatal("negative retransmissions must be rejected") }
}
