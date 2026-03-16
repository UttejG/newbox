package runner

import (
	"context"
	"os/exec"
	"strings"

	"github.com/uttejg/newbox/internal/core/port"
)

// ExecRunner executes commands for real.
type ExecRunner struct{}

func (r *ExecRunner) Run(ctx context.Context, cmd string, args []string) (*port.RunResult, error) {
	c := exec.CommandContext(ctx, cmd, args...)
	out, err := c.CombinedOutput()
	result := &port.RunResult{
		Stdout:  string(out),
		Command: cmd + " " + strings.Join(args, " "),
	}
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		}
		return result, err
	}
	return result, nil
}
