package main

import "sync/atomic"

type FlowRuntimeMetrics struct {
	Packets uint64 `json:"packets"`
	Accepted uint64 `json:"accepted"`
	Rejected uint64 `json:"rejected"`
	Records uint64 `json:"records"`
	QueueDrops uint64 `json:"queue_drops"`
	SequenceGaps uint64 `json:"sequence_gaps"`
}

type FlowMetrics struct {
	packets uint64
	accepted uint64
	rejected uint64
	records uint64
	queueDrops uint64
	sequenceGaps uint64
}

func (m *FlowMetrics) AddPackets() { atomic.AddUint64(&m.packets, 1) }
func (m *FlowMetrics) AddAccepted() { atomic.AddUint64(&m.accepted, 1) }
func (m *FlowMetrics) AddRejected() { atomic.AddUint64(&m.rejected, 1) }
func (m *FlowMetrics) AddRecords(n uint64) { atomic.AddUint64(&m.records, n) }
func (m *FlowMetrics) AddQueueDrops() { atomic.AddUint64(&m.queueDrops, 1) }
func (m *FlowMetrics) AddSequenceGaps(n uint64) { atomic.AddUint64(&m.sequenceGaps, n) }
func (m *FlowMetrics) Snapshot() FlowRuntimeMetrics {
	return FlowRuntimeMetrics{
		Packets: atomic.LoadUint64(&m.packets),
		Accepted: atomic.LoadUint64(&m.accepted),
		Rejected: atomic.LoadUint64(&m.rejected),
		Records: atomic.LoadUint64(&m.records),
		QueueDrops: atomic.LoadUint64(&m.queueDrops),
		SequenceGaps: atomic.LoadUint64(&m.sequenceGaps),
	}
}
