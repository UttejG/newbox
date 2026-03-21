package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/uttejg/newbox/internal/core/domain"
	"github.com/uttejg/newbox/internal/core/port"
)

// InstallService orchestrates preflight checks, planning, and package installation.
type InstallService struct {
	pkgMgr  port.PackageManager
	checker port.SystemChecker
	store   port.StateStore // may be nil for no-resume mode
	dryRun  bool
}

// NewInstallService creates an InstallService. Pass nil for store to disable resume support.
func NewInstallService(pkgMgr port.PackageManager, checker port.SystemChecker, store port.StateStore, dryRun bool) *InstallService {
	return &InstallService{pkgMgr: pkgMgr, checker: checker, store: store, dryRun: dryRun}
}

func (s *InstallService) Preflight(ctx context.Context) (*domain.PreflightResult, error) {
	result := &domain.PreflightResult{}

	if err := s.checker.CheckInternet(ctx); err != nil {
		result.Errors = append(result.Errors, err.Error())
	} else {
		result.InternetOK = true
	}

	if err := s.checker.CheckDiskSpace(ctx, 5); err != nil {
		result.Errors = append(result.Errors, err.Error())
	} else {
		result.DiskSpaceOK = true
	}

	// Use the package manager's own availability check rather than shelling out
	// "<name> --version", which fails for composite managers that have no binary.
	if s.pkgMgr.IsAvailable(ctx) {
		result.PackageManagerOK = true
	} else {
		result.Errors = append(result.Errors, fmt.Sprintf("package manager %q not available", s.pkgMgr.Name()))
	}

	return result, nil
}

func (s *InstallService) Plan(ctx context.Context, selection *domain.UserSelection) (*domain.InstallPlan, error) {
	plan := &domain.InstallPlan{DryRun: s.dryRun}

	var targetOS domain.OS
	if selection.Platform != nil {
		targetOS = selection.Platform.OS
	}

	for _, tool := range selection.AllTools() {
		ref := tool.PackageRefFor(targetOS)
		if ref == nil || ref.IsEmpty() {
			continue
		}

		step := domain.ExecutionStep{
			Tool:    tool,
			Ref:     *ref,
			Command: s.pkgMgr.BuildCommand(*ref),
		}

		if s.dryRun {
			// In dry-run mode mark everything as "would install" — we can't
			// reliably query real install status through a dry-run runner.
			step.Status = domain.StatusDryRun
		} else {
			installed, err := s.pkgMgr.IsInstalled(ctx, *ref)
			if err != nil {
				return nil, fmt.Errorf("checking %s: %w", tool.Name, err)
			}
			if installed {
				step.Status = domain.StatusSkipped
			} else {
				step.Status = domain.StatusPending
			}
		}

		plan.Steps = append(plan.Steps, step)
	}

	return plan, nil
}

func (s *InstallService) Execute(ctx context.Context, plan *domain.InstallPlan, progress chan<- domain.ProgressEvent) error {
	// Load or initialise state for resume tracking.
	var state *domain.InstallState
	if s.store != nil {
		loaded, err := s.store.Load()
		if err == nil && loaded != nil {
			state = loaded
		}
	}
	if state == nil {
		state = &domain.InstallState{
			StartedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}

	steps := plan.PendingSteps()
	total := len(steps)
	var failedTools []string
	var saveWarnings []string

	for i := range steps {
		// Honour context cancellation between steps.
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		step := &steps[i]

		// Resume: skip tools already completed in a previous run.
		if state.IsCompleted(step.Tool.Name) {
			if progress != nil {
				step.Status = domain.StatusDone
				progress <- domain.ProgressEvent{Step: *step, Index: i, Total: total}
			}
			continue
		}

		if progress != nil {
			step.Status = domain.StatusInstalling
			progress <- domain.ProgressEvent{Step: *step, Index: i, Total: total}
		}

		res, err := s.pkgMgr.Install(ctx, step.Ref)
		if res != nil {
			step.Output = res.Stdout
			step.Command = res.Command
		}
		if err != nil {
			step.Status = domain.StatusFailed
			step.Error = err
			failedTools = append(failedTools, step.Tool.Name)
			state.FailedIDs = append(state.FailedIDs, step.Tool.Name)
			if s.store != nil {
				if saveErr := s.store.Save(state); saveErr != nil {
					saveWarnings = append(saveWarnings, saveErr.Error())
				}
			}
		} else {
			step.Status = domain.StatusDone
			state.MarkCompleted(step.Tool.Name)
			if s.store != nil {
				if saveErr := s.store.Save(state); saveErr != nil {
					saveWarnings = append(saveWarnings, saveErr.Error())
				}
			}
		}

		if progress != nil {
			progress <- domain.ProgressEvent{Step: *step, Index: i, Total: total}
		}
	}

	// Only clear persisted state when everything succeeded; preserve it on
	// partial failure so the user can resume.
	if s.store != nil && len(failedTools) == 0 {
		_ = s.store.Clear()
	}

	if len(failedTools) > 0 {
		return fmt.Errorf("%d tool(s) failed to install: %s", len(failedTools), strings.Join(failedTools, ", "))
	}
	if len(saveWarnings) > 0 {
		return fmt.Errorf("warning: resume state could not be persisted: %s", strings.Join(saveWarnings, "; "))
	}
	return nil
}
