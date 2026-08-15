package dynamic

import "sort"

// StateVersion identifies a monotonically ordered FTNDNS state snapshot.
type StateVersion struct {
	Epoch  uint64
	Source string
}

// SyncState reconciles global, regional and local record snapshots without
// making any tier authoritative until its state has been explicitly confirmed.
type SyncState struct {
	Version   StateVersion
	Confirmed bool
	Records   []Record
}

// MergeConfirmed accepts the newest confirmed snapshot. Equal epochs from
// different sources are rejected to avoid silent split-brain selection.
func MergeConfirmed(current, incoming SyncState) (SyncState, bool) {
	if !incoming.Confirmed || incoming.Version.Source == "" {
		return current, false
	}
	if incoming.Version.Epoch < current.Version.Epoch {
		return current, false
	}
	if incoming.Version.Epoch == current.Version.Epoch && incoming.Version.Source != current.Version.Source {
		return current, false
	}
	out := incoming
	out.Records = append([]Record(nil), incoming.Records...)
	sort.SliceStable(out.Records, func(i, j int) bool { return key(out.Records[i]) < key(out.Records[j]) })
	return out, true
}
