package dns

import (
	"sort"
	"strings"

	"github.com/beparykamrul-dev/FTN-AI/internal/platform/module"
)

func Definition() module.Definition {
	caps := []string{"powerdns", "technitium", "coredns", "unbound", "dnsdist", "dns-sync"}
	sort.Strings(caps)
	return module.Definition{Name: "dns", Version: "v1", Capabilities: caps}
}

func ValidDefinition(d module.Definition) bool {
	if strings.TrimSpace(d.Name) != "dns" || strings.TrimSpace(d.Version) == "" || len(d.Capabilities) == 0 {
		return false
	}
	seen := map[string]struct{}{}
	for _, c := range d.Capabilities {
		c = strings.TrimSpace(c)
		if c == "" {
			return false
		}
		if _, ok := seen[c]; ok {
			return false
		}
		seen[c] = struct{}{}
	}
	return true
}
