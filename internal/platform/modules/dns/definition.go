package dns

import "github.com/beparykamrul-dev/FTN-AI/internal/platform/module"

func Definition() module.Definition {
	return module.Definition{
		Name: "dns", Version: "v1",
		Capabilities: []string{"powerdns", "technitium", "coredns", "unbound", "dnsdist", "dns-sync"},
	}
}
