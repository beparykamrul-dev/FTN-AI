package fiber

import "context"

type NodeType string
const (
	NodeOLT NodeType = "olt"
	NodeFiber NodeType = "fiber"
	NodeTable NodeType = "table"
	NodeJoint NodeType = "joint"
	NodeSplitter NodeType = "splitter"
	NodeONU NodeType = "onu"
	NodeRouter NodeType = "router"
	NodeUser NodeType = "pppoe-user"
)

type GraphNode struct { ID string; Type NodeType; Location Point }
type GraphEdge struct { ID, From, To string; DistanceMeters, LossDb float64; Status string }

type FiberGraph struct { Nodes []GraphNode; Edges []GraphEdge }

// GraphRepository supports topology impact analysis without coupling the map
// layer to a particular database implementation.
type GraphRepository interface {
	LoadGraph(context.Context, string) (FiberGraph, error)
	SaveGraph(context.Context, FiberGraph) error
}

// Impact summarizes the service blast radius of a failed segment/node.
type Impact struct {
	NodeID string
	AffectedEdges int
	AffectedONUs int
	AffectedUsers int
	AffectedServices int
}

type ImpactAnalyzer interface { Analyze(context.Context, FiberGraph, string) (Impact, error) }
