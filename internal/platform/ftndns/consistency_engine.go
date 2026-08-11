package ftndns

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
)

// Snapshot represents normalized DNS state obtained from one provider/mesh node.
type Snapshot struct {
	Provider ProviderID
	Version string
	Zones map[string][]string
}

type ConsistencyResult struct {
	Consistent bool
	ReferenceHash string
	Mismatches []string
}

// CompareSnapshots deterministically compares normalized provider state.
// It is intentionally provider-agnostic: reconciliation is a separate guarded step.
func CompareSnapshots(snapshots []Snapshot) ConsistencyResult {
	if len(snapshots) < 2 { return ConsistencyResult{Consistent: true} }
	ordered := append([]Snapshot(nil), snapshots...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].Provider < ordered[j].Provider })

	hashOf := func(s Snapshot) string {
		keys := make([]string, 0, len(s.Zones))
		for k := range s.Zones { keys = append(keys, k) }
		sort.Strings(keys)
		h := sha256.New()
		for _, k := range keys {
			values := append([]string(nil), s.Zones[k]...)
			sort.Strings(values)
			h.Write([]byte(k + "\x00"))
			for _, v := range values { h.Write([]byte(v + "\x00")) }
		}
		return hex.EncodeToString(h.Sum(nil))
	}

	ref := hashOf(ordered[0])
	result := ConsistencyResult{Consistent: true, ReferenceHash: ref}
	for _, s := range ordered[1:] {
		if hashOf(s) != ref {
			result.Consistent = false
			result.Mismatches = append(result.Mismatches, string(s.Provider))
		}
	}
	return result
}
