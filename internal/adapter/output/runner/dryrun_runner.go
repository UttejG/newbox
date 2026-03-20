package runner

import (
	"context"
	"strings"

	"github.com/uttejg/newbox/internal/core/port"
)

// DryRunRunner logs commands instead of executing them.
type DryRunRunner struct {
	Log []port.RunResult
}

func (r *DryRunRunner) Run(ctx context.Context, cmd string, args []string) (*port.RunResult, error) {
	full := cmd + " " + strings.Join(args, " ")
	result := &port.RunResult{
		Stdout:   "[dry-run] would execute: " + full,
		Command:  full,
		DryRun:   true,
		ExitCode: 0,
	}
	r.Log = append(r.Log, *result)
	return result, nil
}
