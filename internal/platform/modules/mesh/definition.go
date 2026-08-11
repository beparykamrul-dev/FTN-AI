package mesh

import "github.com/beparykamrul-dev/FTN-AI/internal/platform/module"

func Definition() module.Definition {
	return module.Definition{
		Name: "mesh", Version: "v1",
		Capabilities: []string{"topology", "link-state", "routing", "failover", "multipath"},
	}
}
