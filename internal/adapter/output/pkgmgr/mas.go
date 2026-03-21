package pkgmgr

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/uttejg/newbox/internal/core/domain"
	"github.com/uttejg/newbox/internal/core/port"
)

// MASManager implements PackageManager for Mac App Store via the `mas` CLI.
type MASManager struct {
	runner       port.CommandRunner
	cacheOnce    sync.Once
	installedIDs map[string]struct{}
	cacheErr     error
}

// NewMAS creates a MASManager backed by the given CommandRunner.
func NewMAS(runner port.CommandRunner) *MASManager {
	return &MASManager{runner: runner}
}

func (m *MASManager) Name() string { return "mas" }

func (m *MASManager) CanHandle(ref domain.PackageRef) bool { return ref.MAS != "" }

func (m *MASManager) IsAvailable(ctx context.Context) error {
	res, err := m.runner.Run(ctx, "mas", []string{"version"})
	if err != nil {
		return fmt.Errorf("mas: %w", err)
	}
	if res.ExitCode != 0 {
		return fmt.Errorf("mas: exited with code %d", res.ExitCode)
	}
	return nil
}

func (m *MASManager) loadCache(ctx context.Context) {
	m.cacheOnce.Do(func() {
		m.installedIDs = make(map[string]struct{})
		res, err := m.runner.Run(ctx, "mas", []string{"list"})
		if err != nil {
			m.cacheErr = fmt.Errorf("mas list: %w", err)
			return
		}
		if res.DryRun {
			return
		}
		for _, line := range strings.Split(res.Stdout, "\n") {
			parts := strings.SplitN(strings.TrimSpace(line), " ", 2)
			if len(parts) >= 1 && parts[0] != "" {
				m.installedIDs[parts[0]] = struct{}{}
			}
		}
	})
}

func (m *MASManager) IsInstalled(ctx context.Context, ref domain.PackageRef) (bool, error) {
	if ref.MAS == "" {
		return false, nil
	}
	m.loadCache(ctx)
	if m.cacheErr != nil {
		return false, m.cacheErr
	}
	_, found := m.installedIDs[ref.MAS]
	return found, nil
}

func (m *MASManager) Install(ctx context.Context, ref domain.PackageRef) (*port.RunResult, error) {
	if ref.MAS == "" {
		return nil, fmt.Errorf("no MAS ID for this package")
	}
	return m.runner.Run(ctx, "mas", []string{"install", ref.MAS})
}

// BuildCommand returns the mas install command string for plan display.
func (m *MASManager) BuildCommand(ref domain.PackageRef) string {
	if ref.MAS == "" {
		return ""
	}
	return "mas install " + ref.MAS
}
