package mesh

import (
	"sort"
	"strings"
	"github.com/beparykamrul-dev/FTN-AI/internal/platform/module
)

func Definition() module.Definition { caps:=[]string{"topology","link-state","routing","failover","multipath"}; sort.Strings(caps); return module.Definition{Name:"mesh",Version:"v1",Capabilities:caps} }
func ValidDefinition(d module.Definition) bool { return strings.TrimSpace(d.Name)=="mesh" && strings.TrimSpace(d.Version)!="" && len(d.Capabilities)>0 }
