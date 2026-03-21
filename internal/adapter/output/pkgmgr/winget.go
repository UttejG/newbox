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

func (w *WingetManager) IsAvailable(ctx context.Context) bool {
	res, err := w.runner.Run(ctx, "winget", []string{"--version"})
	return err == nil && res != nil && res.ExitCode == 0
}

func (w *WingetManager) IsInstalled(ctx context.Context, ref domain.PackageRef) (bool, error) {
	if ref.Winget == "" {
		return false, nil
	}
	res, err := w.runner.Run(ctx, "winget", []string{"list", "--id", ref.Winget, "--exact"})
	if err != nil {
		return false, nil
	}
	return res.ExitCode == 0 && strings.Contains(res.Stdout, ref.Winget), nil
}

func (w *WingetManager) BuildCommand(ref domain.PackageRef) string {
	if ref.Winget == "" {
		return ""
	}
	return "winget install --id " + ref.Winget + " --exact --silent"
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
		"--accept-package-agreements",
	})
}
