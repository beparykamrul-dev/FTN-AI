package dns

import (
	"reflect"
	"sort"
	"sync"
)

type ZoneChange struct {
	Zone string `json:"zone"`
	Version uint64 `json:"version"`
	Records []ZoneRecord `json:"records"`
}

type SyncResult struct {
	Zone string `json:"zone"`
	Applied bool `json:"applied"`
	Conflict bool `json:"conflict"`
	Version uint64 `json:"version"`
}

// SyncEngine reconciles versioned zones between FTN DNS nodes. It is transport
// agnostic: the caller supplies the remote snapshot and decides how it is sent.
type SyncEngine struct {
	mu sync.RWMutex
	zones map[string]Zone
}

func NewSyncEngine() *SyncEngine { return &SyncEngine{zones: make(map[string]Zone)} }

func (e *SyncEngine) Apply(change ZoneChange) SyncResult {
	e.mu.Lock(); defer e.mu.Unlock()
	current, exists := e.zones[normalizeZone(change.Zone)]
	result := SyncResult{Zone: normalizeZone(change.Zone), Version: current.Version}
	if exists && change.Version < current.Version { result.Conflict = true; return result }
	if exists && change.Version == current.Version && !reflect.DeepEqual(current.Records, change.Records) { result.Conflict = true; return result }
	if !exists || change.Version > current.Version {
		records := append([]ZoneRecord(nil), change.Records...)
		sort.SliceStable(records, func(i, j int) bool { return records[i].Name+records[i].Type+records[i].Value < records[j].Name+records[j].Type+records[j].Value })
		e.zones[result.Zone] = Zone{Name: result.Zone, Records: records, Version: change.Version}
		result.Applied = true; result.Version = change.Version
	}
	return result
}

func (e *SyncEngine) Get(zone string) (Zone, bool) {
	e.mu.RLock(); defer e.mu.RUnlock()
	z, ok := e.zones[normalizeZone(zone)]
	return z, ok
}
