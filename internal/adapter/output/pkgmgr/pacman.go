package pkgmgr

import (
	"context"
	"fmt"

	"github.com/uttejg/newbox/internal/core/domain"
	"github.com/uttejg/newbox/internal/core/port"
)

// PacmanManager implements PackageManager for Arch Linux (pacman).
type PacmanManager struct {
	runner port.CommandRunner
}

// NewPacman creates a PacmanManager backed by the given CommandRunner.
func NewPacman(runner port.CommandRunner) *PacmanManager {
	return &PacmanManager{runner: runner}
}

func (p *PacmanManager) Name() string { return "pacman" }

func (p *PacmanManager) CanHandle(ref domain.PackageRef) bool { return ref.Pacman != "" }

func (p *PacmanManager) IsAvailable(ctx context.Context) error {
	res, err := p.runner.Run(ctx, "pacman", []string{"--version"})
	if err != nil {
		return fmt.Errorf("pacman: %w", err)
	}
	if res.ExitCode != 0 {
		return fmt.Errorf("pacman: exited with code %d", res.ExitCode)
	}
	return nil
}

func (p *PacmanManager) BuildCommand(ref domain.PackageRef) string {
	if ref.Pacman == "" {
		return ""
	}
	return "pacman -S --noconfirm " + ref.Pacman
}

func (p *PacmanManager) IsInstalled(ctx context.Context, ref domain.PackageRef) (bool, error) {
	if ref.Pacman == "" {
		return false, nil
	}
	res, err := p.runner.Run(ctx, "pacman", []string{"-Q", ref.Pacman})
	if err != nil {
		return false, nil
	}
	if res.DryRun {
		return false, nil
	}
	return res.ExitCode == 0, nil
}

func (p *PacmanManager) Install(ctx context.Context, ref domain.PackageRef) (*port.RunResult, error) {
	if ref.Pacman == "" {
		return nil, fmt.Errorf("no pacman package for this tool")
	}
	return p.runner.Run(ctx, "pacman", []string{"-S", "--noconfirm", ref.Pacman})
}
