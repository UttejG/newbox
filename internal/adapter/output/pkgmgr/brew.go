package pkgmgr

import (
	"context"
	"fmt"

	"github.com/uttejg/newbox/internal/core/domain"
	"github.com/uttejg/newbox/internal/core/port"
)

// BrewManager implements PackageManager for macOS Homebrew.
type BrewManager struct {
	runner port.CommandRunner
}

// NewBrew creates a BrewManager backed by the given CommandRunner.
func NewBrew(runner port.CommandRunner) *BrewManager {
	return &BrewManager{runner: runner}
}

func (b *BrewManager) Name() string { return "brew" }

func (b *BrewManager) CanHandle(ref domain.PackageRef) bool {
	return ref.Formula != "" || ref.Cask != ""
}

func (b *BrewManager) IsAvailable(ctx context.Context) bool {
	res, err := b.runner.Run(ctx, "brew", []string{"--version"})
	return err == nil && res.ExitCode == 0
}

func (b *BrewManager) IsInstalled(ctx context.Context, ref domain.PackageRef) (bool, error) {
	// In dry-run mode the runner returns a fake result; treat that as "not installed"
	// so all tools appear in the install plan.
	var args []string
	if ref.Cask != "" {
		args = []string{"list", "--cask", ref.Cask}
	} else if ref.Formula != "" {
		args = []string{"list", ref.Formula}
	} else {
		return false, nil
	}

	res, err := b.runner.Run(ctx, "brew", args)
	if res != nil && res.DryRun {
		return false, nil
	}
	if res != nil && res.ExitCode != 0 {
		// brew list exits non-zero when the formula/cask is not installed — not an error.
		return false, nil
	}
	if err != nil {
		// Command failed to execute entirely (brew not on PATH, context cancelled, etc.).
		return false, fmt.Errorf("checking %s: %w", ref.Formula+ref.Cask, err)
	}
	return true, nil
}

func (b *BrewManager) Install(ctx context.Context, ref domain.PackageRef) (*port.RunResult, error) {
	args := b.installArgs(ref)
	if args == nil {
		return nil, fmt.Errorf("brew: no formula or cask in package ref for %q", ref.Formula+ref.Cask)
	}
	return b.runner.Run(ctx, "brew", args)
}

func (b *BrewManager) installArgs(ref domain.PackageRef) []string {
	if ref.Cask != "" {
		return []string{"install", "--cask", ref.Cask}
	}
	if ref.Formula != "" {
		return []string{"install", ref.Formula}
	}
	return nil
}

// BuildCommand returns the brew install command string for plan display.
func (b *BrewManager) BuildCommand(ref domain.PackageRef) string {
	if ref.Cask != "" {
		return "brew install --cask " + ref.Cask
	}
	if ref.Formula != "" {
		return "brew install " + ref.Formula
	}
	return ""
}
