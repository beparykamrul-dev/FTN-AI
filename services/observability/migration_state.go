package observability

// MigrationState tracks a replica movement lifecycle without executing the movement itself.
type MigrationState string

const (
	MigrationPlanned MigrationState = "planned"
	MigrationApproved MigrationState = "approved"
	MigrationRunning MigrationState = "running"
	MigrationVerified MigrationState = "verified"
	MigrationRolledBack MigrationState = "rolled_back"
)

func (s MigrationState) Terminal() bool {
	return s == MigrationVerified || s == MigrationRolledBack
}
