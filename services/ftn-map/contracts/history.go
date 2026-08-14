package fiber

import "context"

type TopologyChange struct {
	ID string
	EntityID string
	Action string
	ActorID string
	Source string
	BeforeHash string
	AfterHash string
	CreatedAt string
}

type HistoryRepository interface {
	Append(context.Context, TopologyChange) error
	History(context.Context, string) ([]TopologyChange, error)
}
