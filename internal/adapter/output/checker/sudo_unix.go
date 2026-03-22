//go:build !windows

package checker

import (
	"context"
	"fmt"
	"runtime"
)

// CheckSudo verifies that sudo credentials are cached (non-interactive).
// It is only meaningful on Linux, where apt/dnf/pacman require root; on macOS
// brew runs without root so the check is skipped.
func (c *SystemChecker) CheckSudo(ctx context.Context) error {
	if runtime.GOOS != "linux" {
		return nil
	}
	res, err := c.Runner.Run(ctx, "sudo", []string{"-n", "true"})
	if err != nil {
		return fmt.Errorf("failed to run sudo: %w", err)
	}
	if res == nil {
		return fmt.Errorf("failed to run sudo: no result returned")
	}
	if res.ExitCode != 0 {
		return fmt.Errorf("sudo credentials not cached — run 'sudo -v' first")
	}
	return nil
}
