package orchestration

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Runtime string

const (
	RuntimeDocker Runtime = "docker"
	RuntimePodman Runtime = "podman"
	RuntimeKubernetes Runtime = "kubernetes"
)

type Host struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Address string `json:"address"`
	Runtime Runtime `json:"runtime"`
	AgentID string `json:"agent_id,omitempty"`
	Status string `json:"status"`
	LastSeen time.Time `json:"last_seen,omitempty"`
}

type Workload struct {
	ID string `json:"id"`
	HostID string `json:"host_id"`
	Name string `json:"name"`
	Image string `json:"image,omitempty"`
	DesiredState string `json:"desired_state"`
	ObservedState string `json:"observed_state"`
	Version string `json:"version,omitempty"`
}

type FleetStore struct { mu sync.RWMutex; hosts map[string]Host; workloads map[string]Workload }

func NewFleetStore() *FleetStore { return &FleetStore{hosts: make(map[string]Host), workloads: make(map[string]Workload)} }

func (s *FleetStore) UpsertHost(h Host) error {
	if strings.TrimSpace(h.ID) == "" || strings.TrimSpace(h.Address) == "" { return fmt.Errorf("host id and address are required") }
	if h.Runtime == "" { return fmt.Errorf("runtime is required") }
	s.mu.Lock(); s.hosts[h.ID] = h; s.mu.Unlock(); return nil
}

func (s *FleetStore) UpsertWorkload(w Workload) error {
	if strings.TrimSpace(w.ID) == "" || strings.TrimSpace(w.HostID) == "" { return fmt.Errorf("workload id and host id are required") }
	s.mu.Lock(); s.workloads[w.ID] = w; s.mu.Unlock(); return nil
}

func (s *FleetStore) Hosts() []Host { s.mu.RLock(); defer s.mu.RUnlock(); out:=make([]Host,0,len(s.hosts)); for _,h:=range s.hosts { out=append(out,h) }; return out }
func (s *FleetStore) Workloads() []Workload { s.mu.RLock(); defer s.mu.RUnlock(); out:=make([]Workload,0,len(s.workloads)); for _,w:=range s.workloads { out=append(out,w) }; return out }
