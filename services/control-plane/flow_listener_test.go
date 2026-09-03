package main

import (
    "context"
    "encoding/binary"
    "net"
    "testing"
    "time"
)

func TestFlowListenerReceivesNetFlowV5(t *testing.T) {
    runtime := NewTrafficRuntime()
    _ = runtime.UpsertEndpoint(ManagedEndpoint{ServiceID: "pubg", CIDR: "198.51.100.0/24"})
    collector := NewFlowTelemetryCollector()
    l, err := NewFlowListener(FlowListenerConfig{Address: "127.0.0.1:0", MaxPacketSize: 4096, QueueSize: 8, Exporters: []string{"127.0.0.1"}}, collector, runtime)
    if err != nil { t.Fatal(err) }
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    if err := l.Start(ctx); err != nil { t.Fatal(err) }
    defer l.Close()

    addr := l.conn.LocalAddr().(*net.UDPAddr)
    conn, err := net.DialUDP("udp", nil, addr)
    if err != nil { t.Fatal(err) }
    defer conn.Close()
    p := make([]byte, 72)
    binary.BigEndian.PutUint16(p[0:2], 5)
    binary.BigEndian.PutUint16(p[2:4], 1)
    copy(p[24:28], []byte{192, 0, 2, 1})
    copy(p[28:32], []byte{198, 51, 100, 1})
    binary.BigEndian.PutUint32(p[40:44], 10)
    binary.BigEndian.PutUint32(p[44:48], 100)
    if _, err := conn.Write(p); err != nil { t.Fatal(err) }

    deadline := time.Now().Add(2 * time.Second)
    for time.Now().Before(deadline) {
        if l.Stats().Records == 1 && runtime.FlowCount() == 1 { return }
        time.Sleep(10 * time.Millisecond)
    }
    t.Fatalf("listener did not process flow: stats=%+v", l.Stats())
}

func TestFlowListenerRejectsUnconfiguredExporter(t *testing.T) {
    runtime := NewTrafficRuntime()
    l, err := NewFlowListener(FlowListenerConfig{Address: "127.0.0.1:0", MaxPacketSize: 4096, QueueSize: 8, Exporters: []string{"192.0.2.10"}}, NewFlowTelemetryCollector(), runtime)
    if err != nil { t.Fatal(err) }
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    if err := l.Start(ctx); err != nil { t.Fatal(err) }
    defer l.Close()

    conn, err := net.DialUDP("udp", nil, l.conn.LocalAddr().(*net.UDPAddr))
    if err != nil { t.Fatal(err) }
    defer conn.Close()
    packet := make([]byte, 24)
    binary.BigEndian.PutUint16(packet[:2], 5)
    binary.BigEndian.PutUint16(packet[2:4], 0)
    if _, err := conn.Write(packet); err != nil { t.Fatal(err) }
    time.Sleep(50 * time.Millisecond)
    if l.Stats().Accepted != 0 || l.Stats().Records != 0 {
        t.Fatalf("unconfigured exporter was accepted: %+v", l.Stats())
    }
}
