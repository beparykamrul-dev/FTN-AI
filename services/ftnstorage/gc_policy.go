package ftnstorage

// GCPolicy controls when unreferenced chunks may be reclaimed.
type GCPolicy struct {
	MinAgeSeconds uint64
	RequireAudit  bool
	RequireLock   bool
}

func (p GCPolicy) Valid() bool {
	return p.MinAgeSeconds > 0 && p.RequireAudit && p.RequireLock
}
