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
	runner port.CommandRunner

	cacheOnce    sync.Once
	installedIDs map[string]struct{}
	cacheErr     error
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

// loadCache fetches the list of installed Flatpak apps once and caches it.
func (f *FlatpakManager) loadCache(ctx context.Context) {
	f.cacheOnce.Do(func() {
		f.installedIDs = make(map[string]struct{})
		res, err := f.runner.Run(ctx, "flatpak", []string{"list", "--app", "--columns=application"})
		if err != nil {
			f.cacheErr = fmt.Errorf("flatpak list: %w", err)
			return
		}
		if res.DryRun {
			return
		}
		if res.ExitCode != 0 {
			f.cacheErr = fmt.Errorf("flatpak list: exited with code %d", res.ExitCode)
			return
		}
		for _, line := range strings.Split(res.Stdout, "\n") {
			id := strings.TrimSpace(line)
			if id != "" {
				f.installedIDs[id] = struct{}{}
			}
		}
	})
}

func (f *FlatpakManager) IsInstalled(ctx context.Context, ref domain.PackageRef) (bool, error) {
	if ref.Flatpak == "" {
		return false, nil
	}
	f.loadCache(ctx)
	if f.cacheErr != nil {
		return false, f.cacheErr
	}
	_, found := f.installedIDs[ref.Flatpak]
	return found, nil
}

func (f *FlatpakManager) Install(ctx context.Context, ref domain.PackageRef) (*port.RunResult, error) {
	if ref.Flatpak == "" {
		return nil, fmt.Errorf("no flatpak ID for this tool")
	}
	return f.runner.Run(ctx, "flatpak", []string{"install", "-y", "flathub", ref.Flatpak})
}
