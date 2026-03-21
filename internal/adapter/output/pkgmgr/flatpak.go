package pkgmgr

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/uttejg/newbox/internal/core/domain"
	"github.com/uttejg/newbox/internal/core/port"
)

// FlatpakManager implements PackageManager for Flatpak (cross-distro Linux fallback).
type FlatpakManager struct {
	runner       port.CommandRunner
	loadOnce     sync.Once
	installedRaw string // cached stdout from `flatpak list --app`
}

// NewFlatpak creates a FlatpakManager backed by the given CommandRunner.
func NewFlatpak(runner port.CommandRunner) *FlatpakManager {
	return &FlatpakManager{runner: runner}
}

func (f *FlatpakManager) Name() string { return "flatpak" }

func (f *FlatpakManager) IsAvailable(ctx context.Context) bool {
	res, err := f.runner.Run(ctx, "flatpak", []string{"--version"})
	return err == nil && res.ExitCode == 0
}

func (f *FlatpakManager) IsInstalled(ctx context.Context, ref domain.PackageRef) (bool, error) {
	if ref.Flatpak == "" {
		return false, nil
	}
	f.loadOnce.Do(func() {
		res, err := f.runner.Run(ctx, "flatpak", []string{"list", "--app"})
		if err != nil || res.DryRun {
			return
		}
		f.installedRaw = res.Stdout
	})
	return strings.Contains(f.installedRaw, ref.Flatpak), nil
}

func (f *FlatpakManager) BuildCommand(ref domain.PackageRef) string {
	if ref.Flatpak == "" {
		return ""
	}
	return "flatpak install -y flathub " + ref.Flatpak
}

func (f *FlatpakManager) Install(ctx context.Context, ref domain.PackageRef) (*port.RunResult, error) {
	if ref.Flatpak == "" {
		return nil, fmt.Errorf("no flatpak ID for this tool")
	}
	return f.runner.Run(ctx, "flatpak", []string{"install", "-y", "flathub", ref.Flatpak})
}
