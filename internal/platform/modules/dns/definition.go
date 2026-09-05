package dns

import (
	"sort"
	"strings"
	"github.com/beparykamrul-dev/FTN-AI/internal/platform/module"
)

func Definition() module.Definition {
	caps:=[]string{"powerdns","technitium","coredns","unbound","dnsdist","dns-sync"}; sort.Strings(caps)
	return module.Definition{Name:"dns",Version:"v1",Capabilities:caps}
}

func ValidDefinition(d module.Definition) bool { return strings.TrimSpace(d.Name)=="dns" && strings.TrimSpace(d.Version)!="" && len(d.Capabilities)>0 }
