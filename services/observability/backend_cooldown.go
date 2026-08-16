package observability

// BackendCooldown prevents an unhealthy backend from immediate re-selection.
type BackendCooldown struct {
	Name string
	UntilUnix int64
}

func (c BackendCooldown) Active(nowUnix int64) bool {
	return c.Name != "" && nowUnix < c.UntilUnix
}
