package monitoring

import "github.com/beparykamrul-dev/FTN-AI/internal/platform/module"

func Definition() module.Definition {
	return module.Definition{
		Name: "monitoring", Version: "v1",
		Capabilities: []string{"metrics", "events", "health", "telemetry", "alerts", "observability"},
	}
}
