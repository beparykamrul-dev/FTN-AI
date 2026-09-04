package dns

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
)

type Record struct {
	ID string `json:"id"`
	Zone string `json:"zone"`
	Name string `json:"name"`
	Type string `json:"type"`
	Value string `json:"value"`
	TTL uint32 `json:"ttl"`
	Priority uint16 `json:"priority,omitempty"`
}

type ZoneSnapshot struct {
	Provider ProviderType `json:"provider"`
	Zone string `json:"zone"`
	Version string `json:"version"`
	Hash string `json:"hash"`
	Records []Record `json:"records"`
}

type ZoneDrift struct {
	Zone string `json:"zone"`
	Provider ProviderType `json:"provider"`
	ExpectedHash string `json:"expected_hash"`
	ObservedHash string `json:"observed_hash"`
	Drifted bool `json:"drifted"`
}

func SnapshotZone(ctx context.Context, provider ProviderType, zone string, records []Record) (ZoneSnapshot, error) {
	select { case <-ctx.Done(): return ZoneSnapshot{}, ctx.Err(); default: }
	copyRecords := append([]Record(nil), records...)
	sort.Slice(copyRecords, func(i, j int) bool {
		a := strings.ToLower(copyRecords[i].Name)+"|"+strings.ToUpper(copyRecords[i].Type)+"|"+copyRecords[i].Value
		b := strings.ToLower(copyRecords[j].Name)+"|"+strings.ToUpper(copyRecords[j].Type)+"|"+copyRecords[j].Value
		return a < b
	})
	h := sha256.New()
	for _, r := range copyRecords {
		h.Write([]byte(strings.ToLower(strings.TrimSuffix(r.Name, "."))))
		h.Write([]byte("|")); h.Write([]byte(strings.ToUpper(r.Type)))
		h.Write([]byte("|")); h.Write([]byte(r.Value)); h.Write([]byte("|")); h.Write([]byte(string(r.TTL)))
	}
	return ZoneSnapshot{Provider: provider, Zone: zone, Hash: hex.EncodeToString(h.Sum(nil)), Records: copyRecords}, nil
}

func CompareZoneSnapshot(expected, observed ZoneSnapshot) ZoneDrift {
	return ZoneDrift{Zone: observed.Zone, Provider: observed.Provider, ExpectedHash: expected.Hash, ObservedHash: observed.Hash, Drifted: expected.Hash != observed.Hash}
}
