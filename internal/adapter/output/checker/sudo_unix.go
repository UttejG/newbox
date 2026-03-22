//go:build !windows

package checker

import (
	"context"
	"fmt"
)

// CheckSudo verifies that sudo credentials are cached (non-interactive).
// On macOS and Linux, `sudo -n true` exits 0 only if credentials are valid.
func (c *SystemChecker) CheckSudo(ctx context.Context) error {
	res, err := c.Runner.Run(ctx, "sudo", []string{"-n", "true"})
	if err != nil || (res != nil && res.ExitCode != 0) {
		return fmt.Errorf("sudo credentials not cached — run 'sudo -v' first")
	}
	return nil
}
