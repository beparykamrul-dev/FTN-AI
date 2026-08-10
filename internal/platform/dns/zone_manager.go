package dns

import (
	"strings"
	"sync"
)

type ZoneRecord struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Value string `json:"value"`
	TTL uint32 `json:"ttl"`
}

type Zone struct {
	Name string `json:"name"`
	Records []ZoneRecord `json:"records"`
	Version uint64 `json:"version"`
}

type ZoneManager struct {
	mu sync.RWMutex
	zones map[string]Zone
}

func NewZoneManager() *ZoneManager { return &ZoneManager{zones: make(map[string]Zone)} }

func normalizeZone(name string) string { return strings.ToLower(strings.TrimSuffix(strings.TrimSpace(name), ".")) }

func (m *ZoneManager) Upsert(zone Zone) bool {
	zone.Name = normalizeZone(zone.Name)
	if zone.Name == "" { return false }
	m.mu.Lock(); defer m.mu.Unlock()
	old, ok := m.zones[zone.Name]
	zone.Version = old.Version + 1
	if !ok { zone.Version = 1 }
	m.zones[zone.Name] = zone
	return true
}

func (m *ZoneManager) Get(name string) (Zone, bool) {
	m.mu.RLock(); defer m.mu.RUnlock()
	z, ok := m.zones[normalizeZone(name)]
	return z, ok
}

func (m *ZoneManager) Snapshot() []Zone {
	m.mu.RLock(); defer m.mu.RUnlock()
	out := make([]Zone, 0, len(m.zones))
	for _, z := range m.zones { out = append(out, z) }
	return out
}
