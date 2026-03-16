//go:build windows

package checker

import "context"

// CheckDiskSpace on Windows skips the disk space check (syscall.Statfs not available).
// The check passes unconditionally; use system tools to verify disk space if needed.
func (c *SystemChecker) CheckDiskSpace(_ context.Context, _ int) error {
	return nil
}
