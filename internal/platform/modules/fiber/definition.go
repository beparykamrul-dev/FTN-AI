package fiber

import (
	"sort"
	"strings"
	"github.com/beparykamrul-dev/FTN-AI/internal/platform/module
)

func Definition() module.Definition { caps:=[]string{"fiber-monitoring","gis","topology","fault-events","impact-analysis"}; sort.Strings(caps); return module.Definition{Name:"fiber",Version:"v1",Capabilities:caps} }
func ValidDefinition(d module.Definition) bool { return strings.TrimSpace(d.Name)=="fiber" && strings.TrimSpace(d.Version)!="" && len(d.Capabilities)>0 }
