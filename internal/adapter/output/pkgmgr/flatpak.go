package pkgmgr

import (
	"context"
	"fmt"
	"strings"

	"github.com/uttejg/newbox/internal/core/domain"
	"github.com/uttejg/newbox/internal/core/port"
)

// FlatpakManager implements PackageManager for Flatpak (cross-distro Linux fallback).
type FlatpakManager struct {
	runner port.CommandRunner
}

// NewFlatpak creates a FlatpakManager backed by the given CommandRunner.
func NewFlatpak(runner port.CommandRunner) *FlatpakManager {
	return &FlatpakManager{runner: runner}
}

func (f *FlatpakManager) Name() string { return "flatpak" }

func (f *FlatpakManager) CanHandle(ref domain.PackageRef) bool { return ref.Flatpak != "" }

func (f *FlatpakManager) IsAvailable(ctx context.Context) error {
	res, err := f.runner.Run(ctx, "flatpak", []string{"--version"})
	if err != nil {
		return fmt.Errorf("flatpak: %w", err)
	}
	if res.ExitCode != 0 {
		return fmt.Errorf("flatpak: exited with code %d", res.ExitCode)
	}
	return nil
}

func (f *FlatpakManager) BuildCommand(ref domain.PackageRef) string {
	if ref.Flatpak == "" {
		return ""
	}
	return "flatpak install -y flathub " + ref.Flatpak
}

func (f *FlatpakManager) IsInstalled(ctx context.Context, ref domain.PackageRef) (bool, error) {
	if ref.Flatpak == "" {
		return false, nil
	}
	res, err := f.runner.Run(ctx, "flatpak", []string{"list", "--app"})
	if err != nil {
		return false, nil
	}
	if res.DryRun {
		return false, nil
	}
	return strings.Contains(res.Stdout, ref.Flatpak), nil
}

func (f *FlatpakManager) Install(ctx context.Context, ref domain.PackageRef) (*port.RunResult, error) {
	if ref.Flatpak == "" {
		return nil, fmt.Errorf("no flatpak ID for this tool")
	}
	return f.runner.Run(ctx, "flatpak", []string{"install", "-y", "flathub", ref.Flatpak})
}
