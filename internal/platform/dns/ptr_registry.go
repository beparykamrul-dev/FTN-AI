package dns

import (
	"net"
	"strings"
	"sync"
)

// PTRRecord represents an authoritative reverse-DNS mapping. The registry
// stores validated names only; zone publication is handled by the DNS adapter.
type PTRRecord struct {
	IP string `json:"ip"`
	Hostname string `json:"hostname"`
	TTL uint32 `json:"ttl"`
}

type PTRRegistry struct {
	mu sync.RWMutex
	records map[string]PTRRecord
}

func NewPTRRegistry() *PTRRegistry { return &PTRRegistry{records: make(map[string]PTRRecord)} }

func (r *PTRRegistry) Upsert(record PTRRecord) bool {
	ip := net.ParseIP(record.IP)
	name := strings.TrimSuffix(strings.TrimSpace(record.Hostname), ".")
	if ip == nil || name == "" { return false }
	if record.TTL == 0 { record.TTL = 300 }
	record.IP = ip.String()
	record.Hostname = name
	r.mu.Lock(); defer r.mu.Unlock()
	r.records[record.IP] = record
	return true
}

func (r *PTRRegistry) Get(ip string) (PTRRecord, bool) {
	parsed := net.ParseIP(ip)
	if parsed == nil { return PTRRecord{}, false }
	r.mu.RLock(); defer r.mu.RUnlock()
	record, ok := r.records[parsed.String()]
	return record, ok
}

func (r *PTRRegistry) Delete(ip string) bool {
	parsed := net.ParseIP(ip)
	if parsed == nil { return false }
	r.mu.Lock(); defer r.mu.Unlock()
	if _, ok := r.records[parsed.String()]; !ok { return false }
	delete(r.records, parsed.String())
	return true
}
