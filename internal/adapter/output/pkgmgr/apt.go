package pkgmgr

import (
	"context"
	"fmt"
	"strings"

	"github.com/uttejg/newbox/internal/core/domain"
	"github.com/uttejg/newbox/internal/core/port"
)

// AptManager implements PackageManager for Debian/Ubuntu (apt).
type AptManager struct {
	runner port.CommandRunner
}

// NewApt creates an AptManager backed by the given CommandRunner.
func NewApt(runner port.CommandRunner) *AptManager {
	return &AptManager{runner: runner}
}

func (a *AptManager) Name() string { return "apt" }

func (a *AptManager) IsAvailable(ctx context.Context) bool {
	res, err := a.runner.Run(ctx, "apt-get", []string{"--version"})
	return err == nil && res.ExitCode == 0
}

func (a *AptManager) IsInstalled(ctx context.Context, ref domain.PackageRef) (bool, error) {
	if ref.Apt == "" {
		return false, nil
	}
	res, err := a.runner.Run(ctx, "dpkg-query", []string{"-W", "-f=${Status}", ref.Apt})
	if err != nil {
		return false, nil
	}
	if res.DryRun {
		return false, nil
	}
	return strings.Contains(res.Stdout, "install ok installed"), nil
}

func (a *AptManager) Install(ctx context.Context, ref domain.PackageRef) (*port.RunResult, error) {
	if ref.Apt == "" {
		return nil, fmt.Errorf("no apt package for this tool")
	}
	return a.runner.Run(ctx, "apt-get", []string{"install", "-y", ref.Apt})
}
