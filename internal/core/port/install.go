package port

import (
	"context"

	"github.com/uttejg/newbox/internal/core/domain"
)

// CommandRunner abstracts command execution — swapped for dry-run.
type CommandRunner interface {
	Run(ctx context.Context, cmd string, args []string) (*RunResult, error)
}

// RunResult captures command output.
type RunResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	DryRun   bool
	Command  string // the full command string for display
}

// PackageManager handles installing packages for a specific OS.
type PackageManager interface {
	Name() string
	// CanHandle reports whether this manager can install the given ref.
	CanHandle(ref domain.PackageRef) bool
	IsAvailable(ctx context.Context) bool
	IsInstalled(ctx context.Context, ref domain.PackageRef) (bool, error)
	Install(ctx context.Context, ref domain.PackageRef) (*RunResult, error)
	// BuildCommand returns the shell command string for the given ref (used for
	// plan display and dry-run output). Returns "" if the ref is unsupported.
	BuildCommand(ref domain.PackageRef) string
}

// SystemChecker runs pre-flight checks.
type SystemChecker interface {
	CheckInternet(ctx context.Context) error
	CheckDiskSpace(ctx context.Context, minGB int) error
	CheckPackageManager(ctx context.Context, name string) error
}

// InstallService orchestrates installation.
type InstallService interface {
	Preflight(ctx context.Context) (*domain.PreflightResult, error)
	Plan(ctx context.Context, selection *domain.UserSelection) (*domain.InstallPlan, error)
	Execute(ctx context.Context, plan *domain.InstallPlan, progress chan<- domain.ProgressEvent) error
}
