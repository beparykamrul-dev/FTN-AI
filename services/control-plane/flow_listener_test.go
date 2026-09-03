package main

import (
	"context"
	"encoding/binary"
	"net"
	"testing"
	"time"
)

func TestFlowListenerReceivesNetFlowV5(t *testing.T) {
	collector := NewFlowTelemetryCollector()
	l, err := NewFlowListener(FlowListenerConfig{Address: "127.0.0.1:0", MaxPacketSize: 4096, QueueSize: 8}, collector)
	if err != nil { t.Fatal(err) }
	ctx, cancel := context.WithCancel(context.Background()); defer cancel()
	if err := l.Start(ctx); err != nil { t.Fatal(err) }
	defer l.Close()

	addr := l.conn.LocalAddr().(*net.UDPAddr)
	conn, err := net.DialUDP("udp", nil, addr); if err != nil { t.Fatal(err) }
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
		if l.Stats().Records == 1 { return }
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("listener did not process flow: %+v", l.Stats())
}

func TestFlowListenerBoundsPacket(t *testing.T) {
	l, err := NewFlowListener(FlowListenerConfig{Address: "127.0.0.1:0", MaxPacketSize: 64, QueueSize: 1}, nil)
	if err != nil { t.Fatal(err) }
	ctx, cancel := context.WithCancel(context.Background()); defer cancel()
	if err := l.Start(ctx); err != nil { t.Fatal(err) }
	defer l.Close()
	conn, err := net.DialUDP("udp", nil, l.conn.LocalAddr().(*net.UDPAddr)); if err != nil { t.Fatal(err) }
	defer conn.Close()
	p := make([]byte, 64); binary.BigEndian.PutUint16(p[:2], 5); binary.BigEndian.PutUint16(p[2:4], 1)
	if _, err := conn.Write(p); err != nil { t.Fatal(err) }
	for i := 0; i < 100; i++ { if l.Stats().Packets > 0 { break }; time.Sleep(time.Millisecond) }
	if l.Stats().Packets == 0 { t.Fatal("packet was not received") }
}
