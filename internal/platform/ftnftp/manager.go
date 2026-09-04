package ftnftp

import (
    "fmt"
    "sync"
)

type Object struct {
    ID string `json:"id"`
    TenantID string `json:"tenant_id"`
    Bucket string `json:"bucket"`
    Key string `json:"key"`
    Size int64 `json:"size"`
    ContentType string `json:"content_type,omitempty"`
    Checksum string `json:"checksum,omitempty"`
    OriginNode string `json:"origin_node,omitempty"`
    StorageClass string `json:"storage_class"`
    ReplicationClass string `json:"replication_class"`
    Status string `json:"status"`
    SHA256 string `json:"sha256,omitempty"`
    Version int64 `json:"version,omitempty"`
    PrimaryNode string `json:"primary_node,omitempty"`
    ReplicaNodes []string `json:"replica_nodes,omitempty"`
    State string `json:"state,omitempty"`
}

type Manager struct { mu sync.RWMutex; objects map[string]Object }
func NewManager()*Manager{return &Manager{objects:map[string]Object{}}}
func(m *Manager) Put(o Object)error{if o.ID==""||o.Bucket==""||o.Key==""{return fmt.Errorf("id, bucket and key are required")};m.mu.Lock();defer m.mu.Unlock();m.objects[o.ID]=o;return nil}
func(m *Manager) Get(id string)(Object,bool){m.mu.RLock();defer m.mu.RUnlock();o,ok:=m.objects[id];return o,ok}
func(m *Manager) List()[]Object{m.mu.RLock();defer m.mu.RUnlock();out:=make([]Object,0,len(m.objects));for _,o:=range m.objects{out=append(out,o)};return out}
