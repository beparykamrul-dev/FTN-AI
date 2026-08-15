package dynamic

import (
	"sort"
	"time"
)

// StateRepair describes a safe repair operation from the last confirmed
// state toward the desired state. It never mutates a DNS backend directly.
type StateRepair struct {
	Version   uint64
	Source    string
	CreatedAt time.Time
	Add       []Record
	Remove    []Record
}

// BuildRepair creates a deterministic repair plan only when desired state is
// newer than the confirmed state.
func BuildRepair(confirmedVersion, desiredVersion uint64, source string, current, desired []Record) StateRepair {
	if desiredVersion <= confirmedVersion {
		return StateRepair{Version: confirmedVersion, Source: source, CreatedAt: time.Now().UTC()}
	}
	add, remove := Reconcile(current, desired)
	sort.Slice(add, func(i, j int) bool { return key(add[i]) < key(add[j]) })
	sort.Slice(remove, func(i, j int) bool { return key(remove[i]) < key(remove[j]) })
	return StateRepair{Version: desiredVersion, Source: source, CreatedAt: time.Now().UTC(), Add: add, Remove: remove}
}
