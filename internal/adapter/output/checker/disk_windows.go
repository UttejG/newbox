//go:build windows

package checker

import (
	"context"
	"fmt"
)

// CheckDiskSpace is not yet implemented on Windows.
func (c *SystemChecker) CheckDiskSpace(_ context.Context, _ int) error {
	return fmt.Errorf("disk space check is not supported on Windows")
}
