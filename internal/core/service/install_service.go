package service

import (
	"context"
	"fmt"
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

	if err := s.checker.CheckPackageManager(ctx, s.pkgMgr.Name()); err != nil {
		result.Errors = append(result.Errors, err.Error())
	} else {
		result.PackageManagerOK = true
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
			Command: s.buildCommandString(ref),
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

	for i, step := range steps {
		// Resume: skip tools already completed in a previous run.
		if state.IsCompleted(step.Tool.Name) {
			if progress != nil {
				step.Status = domain.StatusDone
				progress <- domain.ProgressEvent{Step: step, Index: i, Total: total}
			}
			continue
		}

		if progress != nil {
			step.Status = domain.StatusInstalling
			progress <- domain.ProgressEvent{Step: step, Index: i, Total: total}
		}

		res, err := s.pkgMgr.Install(ctx, step.Ref)
		if res != nil {
			step.Output = res.Stdout
			step.Command = res.Command
		}
		if err != nil {
			step.Status = domain.StatusFailed
			step.Error = err
			state.FailedIDs = append(state.FailedIDs, step.Tool.Name)
		} else {
			step.Status = domain.StatusDone
			state.MarkCompleted(step.Tool.Name)
			if s.store != nil {
				_ = s.store.Save(state)
			}
		}

		if progress != nil {
			progress <- domain.ProgressEvent{Step: step, Index: i, Total: total}
		}
	}

	// Clear persisted state once everything has run.
	if s.store != nil {
		_ = s.store.Clear()
	}

	return nil
}

func (s *InstallService) buildCommandString(ref *domain.PackageRef) string {
	if ref.Cask != "" {
		return "brew install --cask " + ref.Cask
	}
	if ref.Formula != "" {
		return "brew install " + ref.Formula
	}
	if ref.MAS != "" {
		return "mas install " + ref.MAS
	}
	return ""
}
