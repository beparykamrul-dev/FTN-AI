package ftnftp

import "time"

type TransferState string
const(
    TransferQueued TransferState="queued"
    TransferRunning TransferState="running"
    TransferCompleted TransferState="completed"
    TransferFailed TransferState="failed"
)

type Transfer struct { ID string `json:"id"`; ObjectID string `json:"object_id"`; Direction string `json:"direction"`; Protocol string `json:"protocol"`; EdgeNode string `json:"edge_node,omitempty"`; NodeID string `json:"node_id,omitempty"`; Bytes int64 `json:"bytes"`; Offset int64 `json:"offset,omitempty"`; Total int64 `json:"total,omitempty"`; ThroughputBps int64 `json:"throughput_bps"`; State TransferState `json:"state"`; StartedAt time.Time `json:"started_at,omitempty"`; FinishedAt time.Time `json:"finished_at,omitempty"`; UpdatedAt time.Time `json:"updated_at,omitempty"`; Error string `json:"error,omitempty"` }

func (t Transfer) Healthy() bool { return t.State == TransferCompleted && t.Error == "" }
