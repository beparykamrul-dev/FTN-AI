package main

import (
	"context"
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type FlowListenerConfig struct { Address string; MaxPacketSize int; QueueSize int }
type FlowIngestStats struct { Packets uint64; Accepted uint64; Rejected uint64; Records uint64; Templates uint64; QueueDrops uint64 }
type FlowIngress struct { packet []byte; exporter string }
type FlowListener struct { cfg FlowListenerConfig; conn *net.UDPConn; collector *FlowTelemetryCollector; queue chan FlowIngress; stats FlowIngestStats; wg sync.WaitGroup; cancel context.CancelFunc }
func NewFlowListener(cfg FlowListenerConfig,collector *FlowTelemetryCollector)(*FlowListener,error){if cfg.Address==""{cfg.Address=":2055"};if cfg.MaxPacketSize<=0||cfg.MaxPacketSize>maxFlowPacketSize{cfg.MaxPacketSize=maxFlowPacketSize};if cfg.QueueSize<=0||cfg.QueueSize>65536{cfg.QueueSize=4096};if collector==nil{collector=NewFlowTelemetryCollector()};return &FlowListener{cfg:cfg,collector:collector,queue:make(chan FlowIngress,cfg.QueueSize)},nil}
func(l *FlowListener)Start(ctx context.Context)error{if l==nil{return errors.New("flow listener is required")};if ctx==nil{return errors.New("context is required")};if err:=ctx.Err();err!=nil{return err};if l.conn!=nil{return errors.New("flow listener already started")};addr,err:=net.ResolveUDPAddr("udp",l.cfg.Address);if err!=nil{return err};conn,err:=net.ListenUDP("udp",addr);if err!=nil{return err};l.conn=conn;workerCtx,cancel:=context.WithCancel(ctx);l.cancel=cancel;l.wg.Add(2);go l.readLoop(workerCtx);go l.processLoop(workerCtx);return nil}
func(l *FlowListener)Close()error{if l==nil{return nil};if l.cancel!=nil{l.cancel()};if l.conn!=nil{_ = l.conn.Close()};l.wg.Wait();return nil}
func(l *FlowListener)readLoop(ctx context.Context){defer l.wg.Done();buf:=make([]byte,l.cfg.MaxPacketSize);for{_ = l.conn.SetReadDeadline(time.Now().Add(2*time.Second));n,addr,err:=l.conn.ReadFromUDP(buf);if err!=nil{if ctx.Err()!=nil{return};if ne,ok:=err.(net.Error);ok&&ne.Timeout(){continue};atomic.AddUint64(&l.stats.Rejected,1);continue};atomic.AddUint64(&l.stats.Packets,1);if n<2||n>l.cfg.MaxPacketSize{atomic.AddUint64(&l.stats.Rejected,1);continue};p:=append([]byte(nil),buf[:n]...);select{case l.queue<-FlowIngress{packet:p,exporter:addr.IP.String()}:atomic.AddUint64(&l.stats.Accepted,1);default:atomic.AddUint64(&l.stats.QueueDrops,1)}}}
func(l *FlowListener)processLoop(ctx context.Context){defer l.wg.Done();for{select{case<-ctx.Done():return;case in:=<-l.queue:result,err:=DispatchFlowPacket(in.packet,in.exporter,l.collector.Templates);if err!=nil{atomic.AddUint64(&l.stats.Rejected,1);continue};atomic.AddUint64(&l.stats.Records,uint64(len(result.Records)));atomic.AddUint64(&l.stats.Templates,uint64(result.Templates))}}}
func(l *FlowListener)Stats()FlowIngestStats{if l==nil{return FlowIngestStats{}};return FlowIngestStats{Packets:atomic.LoadUint64(&l.stats.Packets),Accepted:atomic.LoadUint64(&l.stats.Accepted),Rejected:atomic.LoadUint64(&l.stats.Rejected),Records:atomic.LoadUint64(&l.stats.Records),Templates:atomic.LoadUint64(&l.stats.Templates),QueueDrops:atomic.LoadUint64(&l.stats.QueueDrops)}}
