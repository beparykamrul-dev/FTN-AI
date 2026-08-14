package fiber

import "context"

type Point struct{ Lat, Lon float64 }

type FiberPath struct {
	ID string
	Name string
	Points []Point
	DistanceMeters float64
	Status string
}

type Table struct { ID, Name string; Location Point }
type Joint struct { ID string; Location Point; Status string }
type Splitter struct { ID string; Location Point; Ratio string; Capacity int }
type ONU struct { ID, Serial, Status string; Location Point }
type Router struct { ID, Name, Status string; Location Point }
type PPPoEUser struct { ID, Username, Status string; ONUID, RouterID string; Location Point }
type CutEvent struct { ID string; Location Point; StartedAt string; RestoredAt string; Cause string }

type FiberTopology struct {
	Paths []FiberPath
	Tables []Table
	Joints []Joint
	Splitters []Splitter
	Users []PPPoEUser
	ONUs []ONU
	Routers []Router
	Cuts []CutEvent
}

// Repository abstracts PostgreSQL/PostGIS-backed topology storage.
type Repository interface {
	GetTopology(context.Context, string) (FiberTopology, error)
	UpsertPath(context.Context, FiberPath) error
	UpsertTable(context.Context, Table) error
	UpsertJoint(context.Context, Joint) error
	UpsertSplitter(context.Context, Splitter) error
	UpsertONU(context.Context, ONU) error
	UpsertRouter(context.Context, Router) error
	UpsertUser(context.Context, PPPoEUser) error
	RecordCut(context.Context, CutEvent) error
}

// AIAnalysis is advisory: it identifies probable cuts, path risk and recovery
// candidates; recovery execution remains behind FTN policy/authorization.
type AIAnalysis interface {
	AnalyzePath(context.Context, FiberTopology) (string, error)
	DetectCut(context.Context, FiberTopology) ([]CutEvent, error)
	RecommendRecovery(context.Context, FiberTopology) (string, error)
}
