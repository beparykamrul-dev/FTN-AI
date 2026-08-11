package fiber

import "github.com/beparykamrul-dev/FTN-AI/internal/platform/fiber"

// Module is the FTN Fiber module boundary. The implementation delegates to
// the canonical fiber package so the Control Plane can load the feature as a
// module without coupling itself to GIS, fault, or impact internals.
type Module struct {
	Topology *fiber.Topology
}

func NewModule() *Module { return &Module{Topology: fiber.NewTopology()} }

func (m *Module) Name() string { return "fiber" }
func (m *Module) Version() string { return "v1" }
func (m *Module) Capabilities() []string {
	return []string{"gis", "topology", "fault-events", "impact-analysis", "fiber-monitoring"}
}
