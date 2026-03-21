package pkgmgr

import (
	"context"
	"fmt"
	"strings"

	"github.com/uttejg/newbox/internal/core/domain"
	"github.com/uttejg/newbox/internal/core/port"
)

// WingetManager implements PackageManager for Windows Package Manager (winget).
type WingetManager struct {
	runner port.CommandRunner
}

// NewWinget creates a WingetManager backed by the given CommandRunner.
func NewWinget(runner port.CommandRunner) *WingetManager {
	return &WingetManager{runner: runner}
}

func (w *WingetManager) Name() string { return "winget" }

func (w *WingetManager) CanHandle(ref domain.PackageRef) bool {
	return ref.Winget != ""
}

func (w *WingetManager) IsAvailable(ctx context.Context) error {
	res, err := w.runner.Run(ctx, "winget", []string{"--version"})
	if err != nil {
		return fmt.Errorf("winget: %w", err)
	}
	if res.ExitCode != 0 {
		return fmt.Errorf("winget: exited with code %d", res.ExitCode)
	}
	return nil
}

func (w *WingetManager) IsInstalled(ctx context.Context, ref domain.PackageRef) (bool, error) {
	if ref.Winget == "" {
		return false, nil
	}
	res, err := w.runner.Run(ctx, "winget", []string{"list", "--id", ref.Winget, "--exact"})
	if res != nil && res.DryRun {
		return false, nil
	}
	if res != nil && res.ExitCode != 0 {
		// winget list exits non-zero when the package is not installed — not an error.
		return false, nil
	}
	if err != nil {
		// Command failed to execute entirely (winget not on PATH, context cancelled, etc.).
		return false, fmt.Errorf("checking %s: %w", ref.Winget, err)
	}
	return strings.Contains(res.Stdout, ref.Winget), nil
}

func (w *WingetManager) Install(ctx context.Context, ref domain.PackageRef) (*port.RunResult, error) {
	if ref.Winget == "" {
		return nil, fmt.Errorf("no winget ID for this package")
	}
	return w.runner.Run(ctx, "winget", []string{
		"install",
		"--id", ref.Winget,
		"--exact",
		"--silent",
		"--disable-interactivity",
		"--accept-package-agreements",
		"--accept-source-agreements",
	})
}

// BuildCommand returns the winget install command string for plan display.
func (w *WingetManager) BuildCommand(ref domain.PackageRef) string {
	if ref.Winget == "" {
		return ""
	}
	return fmt.Sprintf("winget install --id %s --exact --silent --disable-interactivity --accept-package-agreements --accept-source-agreements", ref.Winget)
}
