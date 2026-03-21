package pkgmgr

import (
	"context"
	"fmt"
	"strings"

	"github.com/uttejg/newbox/internal/core/domain"
	"github.com/uttejg/newbox/internal/core/port"
)

// DnfManager implements PackageManager for Fedora/RHEL (dnf).
type DnfManager struct {
	runner port.CommandRunner
}

// NewDnf creates a DnfManager backed by the given CommandRunner.
func NewDnf(runner port.CommandRunner) *DnfManager {
	return &DnfManager{runner: runner}
}

func (d *DnfManager) Name() string { return "dnf" }

func (d *DnfManager) CanHandle(ref domain.PackageRef) bool { return ref.Dnf != "" }

func (d *DnfManager) IsAvailable(ctx context.Context) error {
	res, err := d.runner.Run(ctx, "dnf", []string{"--version"})
	if err != nil {
		return fmt.Errorf("dnf: %w", err)
	}
	if res.ExitCode != 0 {
		return fmt.Errorf("dnf: exited with code %d", res.ExitCode)
	}
	return nil
}

func (d *DnfManager) BuildCommand(ref domain.PackageRef) string {
	if ref.Dnf == "" {
		return ""
	}
	return "dnf install -y " + ref.Dnf
}

func (d *DnfManager) IsInstalled(ctx context.Context, ref domain.PackageRef) (bool, error) {
	if ref.Dnf == "" {
		return false, nil
	}
	res, err := d.runner.Run(ctx, "dnf", []string{"list", "installed", ref.Dnf})
	if err != nil {
		return false, nil
	}
	if res.DryRun {
		return false, nil
	}
	return strings.Contains(res.Stdout, ref.Dnf), nil
}

func (d *DnfManager) Install(ctx context.Context, ref domain.PackageRef) (*port.RunResult, error) {
	if ref.Dnf == "" {
		return nil, fmt.Errorf("no dnf package for this tool")
	}
	return d.runner.Run(ctx, "dnf", []string{"install", "-y", ref.Dnf})
}
