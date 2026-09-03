package main

import (
    "context"
    "errors"
    "net"
    "strings"
    "sync"
    "sync/atomic"
    "time"
)

type FlowListenerConfig struct {
    Address       string
    MaxPacketSize int
    QueueSize     int
    Exporters     []string
}

type FlowIngestStats struct {
    Packets    uint64 `json:"packets"`
    Accepted   uint64 `json:"accepted"`
    Rejected   uint64 `json:"rejected"`
    Records    uint64 `json:"records"`
    Templates  uint64 `json:"templates"`
    QueueDrops uint64 `json:"queue_drops"`
}

type FlowIngress struct {
    packet   []byte
    exporter string
}

type FlowListener struct {
    cfg       FlowListenerConfig
    conn      *net.UDPConn
    collector *FlowTelemetryCollector
    runtime   *TrafficRuntime
    exporters map[string]struct{}
    queue     chan FlowIngress
    stats     FlowIngestStats
    wg        sync.WaitGroup
    cancel    context.CancelFunc
}

func NewFlowListener(cfg FlowListenerConfig, collector *FlowTelemetryCollector, runtime *TrafficRuntime) (*FlowListener, error) {
    if strings.TrimSpace(cfg.Address) == "" { cfg.Address = ":2055" }
    if cfg.MaxPacketSize <= 0 || cfg.MaxPacketSize > maxFlowPacketSize { cfg.MaxPacketSize = maxFlowPacketSize }
    if cfg.QueueSize <= 0 || cfg.QueueSize > 65536 { cfg.QueueSize = 4096 }
    if collector == nil { collector = NewFlowTelemetryCollector() }
    if runtime == nil { return nil, errors.New("traffic runtime required") }
    allow := make(map[string]struct{}, len(cfg.Exporters))
    for _, exporter := range cfg.Exporters {
        exporter = strings.TrimSpace(exporter)
        if exporter != "" { allow[exporter] = struct{}{} }
    }
    if len(allow) == 0 { return nil, errors.New("flow exporter allowlist required") }
    return &FlowListener{cfg: cfg, collector: collector, runtime: runtime, exporters: allow, queue: make(chan FlowIngress, cfg.QueueSize)}, nil
}

func (l *FlowListener) Start(ctx context.Context) error {
    if l == nil || l.runtime == nil { return errors.New("traffic runtime required") }
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
    if l == nil { return nil }
    if l.cancel != nil { l.cancel() }
    if l.conn != nil { _ = l.conn.Close() }
    l.wg.Wait(); return nil
}

func (l *FlowListener) readLoop(ctx context.Context) {
    defer l.wg.Done()
    buf := make([]byte, l.cfg.MaxPacketSize)
    for {
        _ = l.conn.SetReadDeadline(time.Now().Add(2 * time.Second))
        n, addr, err := l.conn.ReadFromUDP(buf)
        if err != nil {
            if ctx.Err() != nil { return }
            if ne, ok := err.(net.Error); ok && ne.Timeout() { continue }
            atomic.AddUint64(&l.stats.Rejected, 1); continue
        }
        atomic.AddUint64(&l.stats.Packets, 1)
        if n < 2 || n > l.cfg.MaxPacketSize || !l.allowed(addr) { atomic.AddUint64(&l.stats.Rejected, 1); continue }
        packet := make([]byte, n); copy(packet, buf[:n])
        select {
        case l.queue <- FlowIngress{packet: packet, exporter: addr.IP.String()}: atomic.AddUint64(&l.stats.Accepted, 1)
        default: atomic.AddUint64(&l.stats.QueueDrops, 1)
        }
    }
}

func (l *FlowListener) allowed(addr *net.UDPAddr) bool {
    if addr == nil { return false }
    _, ok := l.exporters[addr.IP.String()]
    return ok
}

func (l *FlowListener) processLoop(ctx context.Context) {
    defer l.wg.Done()
    for {
        select {
        case <-ctx.Done(): return
        case in := <-l.queue:
            result, err := DispatchFlowPacket(in.packet, in.exporter, l.collector.Templates)
            if err != nil { atomic.AddUint64(&l.stats.Rejected, 1); continue }
            atomic.AddUint64(&l.stats.Records, uint64(len(result.Records)))
            atomic.AddUint64(&l.stats.Templates, uint64(result.Templates))
            if len(result.Records) > 0 {
                normalized := make([]FlowRecord, 0, len(result.Records))
                for _, record := range result.Records { normalized = append(normalized, NormalizeSampledCounters(record)) }
                l.runtime.Ingest(normalized, time.Now().UTC())
            }
        }
    }
}

func (l *FlowListener) Stats() FlowIngestStats {
    return FlowIngestStats{
        Packets: atomic.LoadUint64(&l.stats.Packets), Accepted: atomic.LoadUint64(&l.stats.Accepted), Rejected: atomic.LoadUint64(&l.stats.Rejected), Records: atomic.LoadUint64(&l.stats.Records), Templates: atomic.LoadUint64(&l.stats.Templates), QueueDrops: atomic.LoadUint64(&l.stats.QueueDrops),
    }
}
