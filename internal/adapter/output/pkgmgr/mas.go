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
	loadOnce     sync.Once
	installedRaw string // cached stdout from `mas list`
	loadErr      error
}

// NewMAS creates a MASManager backed by the given CommandRunner.
func NewMAS(runner port.CommandRunner) *MASManager {
	return &MASManager{runner: runner}
}

func (m *MASManager) Name() string { return "mas" }

func (m *MASManager) IsAvailable(ctx context.Context) bool {
	res, err := m.runner.Run(ctx, "mas", []string{"version"})
	return err == nil && res.ExitCode == 0
}

func (m *MASManager) IsInstalled(ctx context.Context, ref domain.PackageRef) (bool, error) {
	if ref.MAS == "" {
		return false, nil
	}
	m.loadOnce.Do(func() {
		res, err := m.runner.Run(ctx, "mas", []string{"list"})
		if err != nil {
			m.loadErr = err
			return
		}
		if !res.DryRun {
			m.installedRaw = res.Stdout
		}
	})
	if m.loadErr != nil {
		return false, m.loadErr
	}
	// Match the first whitespace-delimited token (app ID) exactly to avoid
	// false positives where ref.MAS "123" would match a line containing "1234".
	for _, line := range strings.Split(m.installedRaw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) > 0 && parts[0] == ref.MAS {
			return true, nil
		}
	}
	return false, nil
}

func (m *MASManager) BuildCommand(ref domain.PackageRef) string {
	if ref.MAS == "" {
		return ""
	}
	return "mas install " + ref.MAS
}

func (m *MASManager) Install(ctx context.Context, ref domain.PackageRef) (*port.RunResult, error) {
	if ref.MAS == "" {
		return nil, fmt.Errorf("no MAS ID for this package")
	}
	return m.runner.Run(ctx, "mas", []string{"install", ref.MAS})
}
