package pkgmgr

import (
	"context"
	"fmt"
	"strings"

	"github.com/uttejg/newbox/internal/core/domain"
	"github.com/uttejg/newbox/internal/core/port"
)

// MASManager implements PackageManager for Mac App Store via the `mas` CLI.
type MASManager struct {
	runner port.CommandRunner
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
	res, err := m.runner.Run(ctx, "mas", []string{"list"})
	if err != nil {
		return false, err
	}
	if res.DryRun {
		return false, nil
	}
	return strings.Contains(res.Stdout, ref.MAS), nil
}

func (m *MASManager) Install(ctx context.Context, ref domain.PackageRef) (*port.RunResult, error) {
	if ref.MAS == "" {
		return nil, fmt.Errorf("no MAS ID for this package")
	}
	return m.runner.Run(ctx, "mas", []string{"install", ref.MAS})
}
