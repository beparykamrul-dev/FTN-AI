package builder

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

type BuildRunner struct {
	Timeout time.Duration
}

func NewBuildRunner(timeout time.Duration) *BuildRunner {
	if timeout <= 0 { timeout = 15 * time.Minute }
	return &BuildRunner{Timeout: timeout}
}

// Run executes only an explicitly supplied executable and arguments.
// Generated source is never passed to a shell. Callers must run this in an
// isolated container/VM with restricted filesystem, network and credentials.
func (r *BuildRunner) Run(ctx context.Context, dir string, executable string, args ...string) ([]byte, error) {
	if executable == "" { return nil, fmt.Errorf("build executable is required") }
	ctx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, executable, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return out, fmt.Errorf("build failed: %w", err)
	}
	return out, nil
}
