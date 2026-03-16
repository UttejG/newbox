package checker

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"syscall"

	"github.com/uttejg/newbox/internal/core/port"
)

// SystemChecker implements port.SystemChecker using real OS calls.
type SystemChecker struct {
	Runner port.CommandRunner
}

func (c *SystemChecker) CheckInternet(ctx context.Context) error {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, "https://github.com", nil)
	if err != nil {
		return fmt.Errorf("building internet check request: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("no internet connection: %w", err)
	}
	resp.Body.Close()
	return nil
}

func (c *SystemChecker) CheckDiskSpace(ctx context.Context, minGB int) error {
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

func (c *SystemChecker) CheckPackageManager(ctx context.Context, name string) error {
	res, err := c.Runner.Run(ctx, name, []string{"--version"})
	if err != nil || (res != nil && res.ExitCode != 0) {
		return fmt.Errorf("package manager %q not available", name)
	}
	return nil
}
