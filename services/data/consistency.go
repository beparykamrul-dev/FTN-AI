package data

// ConsistencyMode defines the minimum consistency required by a service.
type ConsistencyMode string

const (
	ConsistencyStrong    ConsistencyMode = "strong"
	ConsistencyEventual  ConsistencyMode = "eventual"
)

func ValidConsistency(mode ConsistencyMode) bool {
	return mode == ConsistencyStrong || mode == ConsistencyEventual
}
