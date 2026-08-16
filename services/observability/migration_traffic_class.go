package observability

// MigrationTrafficClass assigns a relative priority to background storage movement.
type MigrationTrafficClass string

const (
	MigrationBulk MigrationTrafficClass = "bulk"
	MigrationCritical MigrationTrafficClass = "critical"
)

func (c MigrationTrafficClass) Valid() bool {
	return c == MigrationBulk || c == MigrationCritical
}
