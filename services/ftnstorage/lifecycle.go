package ftnstorage

// LifecycleState describes the safe lifecycle of a content-addressed chunk.
type LifecycleState string

const (
	ChunkHot     LifecycleState = "hot"
	ChunkWarm    LifecycleState = "warm"
	ChunkCold    LifecycleState = "cold"
	ChunkRetired LifecycleState = "retired"
)

// TransitionAllowed prevents accidental deletion/reclamation before a chunk
// reaches the explicit retired state.
func TransitionAllowed(from, to LifecycleState) bool {
	switch from {
	case ChunkHot:
		return to == ChunkWarm
	case ChunkWarm:
		return to == ChunkHot || to == ChunkCold
	case ChunkCold:
		return to == ChunkWarm || to == ChunkRetired
	default:
		return false
	}
}
