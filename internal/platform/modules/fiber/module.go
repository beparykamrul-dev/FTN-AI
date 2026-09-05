package fiber

import (
	"sort"
	"strings"
	"github.com/beparykamrul-dev/FTN-AI/internal/platform/fiber"
)

type Module struct { Topology *fiber.Topology }
func NewModule() *Module { return &Module{Topology:fiber.NewTopology()} }
func (m *Module) Name() string { return "fiber" }
func (m *Module) Version() string { return "v1" }
func (m *Module) Capabilities() []string { caps:=[]string{"gis","topology","fault-events","impact-analysis","fiber-monitoring"}; sort.Strings(caps); return caps }
func (m *Module) Ready() bool { return m!=nil && m.Topology!=nil }
func (m *Module) Valid() bool { return m!=nil && strings.TrimSpace(m.Name())!="" && strings.TrimSpace(m.Version())!="" && m.Topology!=nil }
