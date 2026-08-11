package fiber

import "github.com/beparykamrul-dev/FTN-AI/internal/platform/module"

// Definition returns the public contract exposed by the Fiber module.
// Runtime implementations stay behind the module boundary.
func Definition() module.Definition {
	return module.Definition{
		Name: "fiber",
		Version: "v1",
		Capabilities: []string{
			"fiber-monitoring",
			"gis",
			"topology",
			"fault-events",
			"impact-analysis",
		},
	}
}
