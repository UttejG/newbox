//go:build !windows

package checker

import (
	"context"
	"fmt"
	"os"
	"syscall"
)

func (c *SystemChecker) CheckDiskSpace(_ context.Context, minGB int) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot determine home directory: %w", err)
	}
	var stat syscall.Statfs_t
	if err := syscall.Statfs(home, &stat); err != nil {
		return fmt.Errorf("cannot check disk space: %w", err)
	}
	availBytes := stat.Bavail * uint64(stat.Bsize) //nolint:gosec
	availGB := availBytes / (1024 * 1024 * 1024)
	if int(availGB) < minGB {
		return fmt.Errorf("insufficient disk space: %dGB available, %dGB required", availGB, minGB)
	}
	return nil
}
