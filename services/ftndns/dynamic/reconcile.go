package dynamic

import "sort"

// Record is the desired FTNDNS record state.
type Record struct {
	Name    string
	Type    string
	Address string
	NodeID  string
	TTL     uint32
}

// Reconcile produces the minimal deterministic set of changes needed to make
// the advertised records match the selected healthy endpoints.
func Reconcile(current, desired []Record) (add, remove []Record) {
	cur := make(map[string]Record, len(current))
	for _, r := range current {
		cur[key(r)] = r
	}
	want := make(map[string]Record, len(desired))
	for _, r := range desired {
		want[key(r)] = r
	}
	for k, r := range want {
		if old, ok := cur[k]; !ok || old != r {
			add = append(add, r)
		}
	}
	for k, r := range cur {
		if _, ok := want[k]; !ok {
			remove = append(remove, r)
		}
	}
	sort.Slice(add, func(i, j int) bool { return key(add[i]) < key(add[j]) })
	sort.Slice(remove, func(i, j int) bool { return key(remove[i]) < key(remove[j]) })
	return add, remove
}

func key(r Record) string {
	return r.Name + "|" + r.Type + "|" + r.Address + "|" + r.NodeID
}
