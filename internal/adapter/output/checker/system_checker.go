package checker

import (
	"context"
	"fmt"
	"net/http"

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

func (c *SystemChecker) CheckPackageManager(ctx context.Context, name string) error {
	res, err := c.Runner.Run(ctx, name, []string{"--version"})
	if err != nil || (res != nil && res.ExitCode != 0) {
		return fmt.Errorf("package manager %q not available", name)
	}
	return nil
}
