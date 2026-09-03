package ftnftp

import "time"

type TransferState string
const(
	TransferQueued TransferState="queued"
	TransferRunning TransferState="running"
	TransferCompleted TransferState="completed"
	TransferFailed TransferState="failed"
)

type Transfer struct {
	ID string `json:"id"`
	ObjectID string `json:"object_id"`
	NodeID string `json:"node_id,omitempty"`
	Direction string `json:"direction,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	EdgeNode string `json:"edge_node,omitempty"`
	Offset int64 `json:"offset"`
	Total int64 `json:"total"`
	Bytes int64 `json:"bytes"`
	ThroughputBps int64 `json:"throughput_bps"`
	StartedAt time.Time `json:"started_at,omitempty"`
	FinishedAt time.Time `json:"finished_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	State TransferState `json:"state"`
	Error string `json:"error,omitempty"`
}

func (t Transfer) Healthy() bool { return t.State == TransferCompleted && t.Error == "" }
