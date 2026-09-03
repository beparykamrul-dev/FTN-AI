package main

import (
	"context"
	"encoding/binary"
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type FlowListenerConfig struct {
	Address       string
	MaxPacketSize int
	QueueSize     int
}

type FlowIngestStats struct {
	Packets uint64
	Accepted uint64
	Rejected uint64
	Records uint64
	QueueDrops uint64
}

type FlowIngress struct {
	packet []byte
	exporter string
}

type FlowListener struct {
	cfg FlowListenerConfig
	conn *net.UDPConn
	collector *FlowTelemetryCollector
	queue chan FlowIngress
	stats FlowIngestStats
	wg sync.WaitGroup
	cancel context.CancelFunc
}

func NewFlowListener(cfg FlowListenerConfig, collector *FlowTelemetryCollector) (*FlowListener, error) {
	if cfg.Address == "" { cfg.Address = ":2055" }
	if cfg.MaxPacketSize <= 0 || cfg.MaxPacketSize > maxFlowPacketSize { cfg.MaxPacketSize = maxFlowPacketSize }
	if cfg.QueueSize <= 0 || cfg.QueueSize > 65536 { cfg.QueueSize = 4096 }
	if collector == nil { collector = NewFlowTelemetryCollector() }
	return &FlowListener{cfg: cfg, collector: collector, queue: make(chan FlowIngress, cfg.QueueSize)}, nil
}

func (l *FlowListener) Start(ctx context.Context) error {
	if l.conn != nil { return errors.New("flow listener already started") }
	addr, err := net.ResolveUDPAddr("udp", l.cfg.Address); if err != nil { return err }
	conn, err := net.ListenUDP("udp", addr); if err != nil { return err }
	l.conn = conn
	workerCtx, cancel := context.WithCancel(ctx); l.cancel = cancel
	l.wg.Add(1); go l.readLoop(workerCtx)
	l.wg.Add(1); go l.processLoop(workerCtx)
	return nil
}

func (l *FlowListener) Close() error {
	if l.cancel != nil { l.cancel() }
	if l.conn != nil { _ = l.conn.Close() }
	l.wg.Wait(); return nil
}

func (l *FlowListener) readLoop(ctx context.Context) {
	defer l.wg.Done()
	buf := make([]byte, l.cfg.MaxPacketSize)
	for {
		_ = l.conn.SetReadDeadline(time.Now().Add(2*time.Second))
		n, addr, err := l.conn.ReadFromUDP(buf)
		if err != nil {
			if ctx.Err() != nil { return }
			if ne, ok := err.(net.Error); ok && ne.Timeout() { continue }
			atomic.AddUint64(&l.stats.Rejected, 1); continue
		}
		atomic.AddUint64(&l.stats.Packets, 1)
		if n < 2 || n > l.cfg.MaxPacketSize { atomic.AddUint64(&l.stats.Rejected, 1); continue }
		p := make([]byte, n); copy(p, buf[:n])
		select {
		case l.queue <- FlowIngress{packet:p, exporter:addr.IP.String()}:
			atomic.AddUint64(&l.stats.Accepted, 1)
		default:
			atomic.AddUint64(&l.stats.QueueDrops, 1)
		}
	}
}

func (l *FlowListener) processLoop(ctx context.Context) {
	defer l.wg.Done()
	for {
		select {
		case <-ctx.Done(): return
		case in := <-l.queue:
			version := binary.BigEndian.Uint16(in.packet[:2])
			var records []FlowRecord
			var err error
		switch version {
			case 5:
				records, err = DecodeNetFlowV5(in.packet, in.exporter)
			default:
				err = errors.New("template-aware v9/ipfix ingest requires protocol-specific set handling")
			}
			if err != nil { atomic.AddUint64(&l.stats.Rejected, 1); continue }
			atomic.AddUint64(&l.stats.Records, uint64(len(records)))
		}
	}
}

func (l *FlowListener) Stats() FlowIngestStats {
	return FlowIngestStats{
		Packets: atomic.LoadUint64(&l.stats.Packets), Accepted: atomic.LoadUint64(&l.stats.Accepted),
		Rejected: atomic.LoadUint64(&l.stats.Rejected), Records: atomic.LoadUint64(&l.stats.Records), QueueDrops: atomic.LoadUint64(&l.stats.QueueDrops),
	}
}
