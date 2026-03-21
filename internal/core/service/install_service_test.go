package service_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/uttejg/newbox/internal/core/domain"
	"github.com/uttejg/newbox/internal/core/service"
	"github.com/uttejg/newbox/internal/testutil"
)

// makeSelection builds a UserSelection for the given tools on macOS.
func makeSelection(tools ...domain.Tool) *domain.UserSelection {
	return &domain.UserSelection{
		Platform: &domain.Platform{OS: domain.OSMacOS},
		ToolsByCategory: map[string][]domain.Tool{
			"test": tools,
		},
	}
}

// ── Preflight ───────────────────────────────────────────────────────────────

func TestInstallService_Preflight_AllOK(t *testing.T) {
	svc := service.NewInstallService(
		&testutil.FakePackageManager{},
		&testutil.FakeSystemChecker{},
		nil,
		false,
		io.Discard,
	)
	result, err := svc.Preflight(context.Background())
	if err != nil {
		t.Fatalf("Preflight() error = %v", err)
	}
	if !result.InternetOK {
		t.Error("expected InternetOK")
	}
	if !result.DiskSpaceOK {
		t.Error("expected DiskSpaceOK")
	}
	if !result.PackageManagerOK {
		t.Error("expected PackageManagerOK")
	}
	if !result.OK() {
		t.Error("expected OK()")
	}
	if len(result.Errors) != 0 {
		t.Errorf("expected no errors, got %v", result.Errors)
	}
}

func TestInstallService_Preflight_Failures(t *testing.T) {
	tests := []struct {
		name       string
		checker    testutil.FakeSystemChecker
		pkg        testutil.FakePackageManager
		wantErrors int
		wantOK     bool
	}{
		{
			name:       "internet failure",
			checker:    testutil.FakeSystemChecker{InternetErr: errTest},
			wantErrors: 1,
			wantOK:     false,
		},
		{
			name:       "disk failure",
			checker:    testutil.FakeSystemChecker{DiskErr: errTest},
			wantErrors: 1,
			wantOK:     false,
		},
		{
			name: "pkgmgr failure",
			pkg: testutil.FakePackageManager{
				IsAvailableFunc: func(_ context.Context) error { return errTest },
			},
			wantErrors: 1,
			wantOK:     false,
		},
		{
			name:    "all failures",
			checker: testutil.FakeSystemChecker{InternetErr: errTest, DiskErr: errTest},
			pkg: testutil.FakePackageManager{
				IsAvailableFunc: func(_ context.Context) error { return errTest },
			},
			wantErrors: 3,
			wantOK:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewInstallService(
				&tt.pkg,
				&tt.checker,
				nil,
				false,
				io.Discard,
			)
			result, err := svc.Preflight(context.Background())
			if err != nil {
				t.Fatalf("Preflight() error = %v", err)
			}
			if len(result.Errors) != tt.wantErrors {
				t.Errorf("Errors count = %d, want %d; got %v", len(result.Errors), tt.wantErrors, result.Errors)
			}
			if result.OK() != tt.wantOK {
				t.Errorf("OK() = %v, want %v", result.OK(), tt.wantOK)
			}
		})
	}
}

// ── Plan ────────────────────────────────────────────────────────────────────

func TestInstallService_Plan_DryRun(t *testing.T) {
	tool := domain.Tool{
		Name:  "Signal",
		MacOS: &domain.PackageRef{Cask: "signal"},
	}
	svc := service.NewInstallService(
		&testutil.FakePackageManager{},
		&testutil.FakeSystemChecker{},
		nil,
		true, // dry-run
		io.Discard,
	)
	plan, err := svc.Plan(context.Background(), makeSelection(tool))
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
	}
	if !plan.DryRun {
		t.Error("expected plan.DryRun=true")
	}
	if len(plan.Steps) != 1 {
		t.Fatalf("expected 1 step, got %d", len(plan.Steps))
	}
	if plan.Steps[0].Status != domain.StatusDryRun {
		t.Errorf("status = %q, want %q", plan.Steps[0].Status, domain.StatusDryRun)
	}
	if plan.Steps[0].Command != "fake install --cask signal" {
		t.Errorf("command = %q", plan.Steps[0].Command)
	}
}

func TestInstallService_Plan_SkipsInstalled(t *testing.T) {
	tool := domain.Tool{
		Name:  "git",
		MacOS: &domain.PackageRef{Formula: "git"},
	}
	svc := service.NewInstallService(
		&testutil.FakePackageManager{InstalledTools: map[string]bool{"git": true}},
		&testutil.FakeSystemChecker{},
		nil,
		false,
		io.Discard,
	)
	plan, err := svc.Plan(context.Background(), makeSelection(tool))
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
	}
	if len(plan.Steps) != 1 {
		t.Fatalf("expected 1 step, got %d", len(plan.Steps))
	}
	if plan.Steps[0].Status != domain.StatusSkipped {
		t.Errorf("status = %q, want %q", plan.Steps[0].Status, domain.StatusSkipped)
	}
}

func TestInstallService_Plan_PendingWhenNotInstalled(t *testing.T) {
	tool := domain.Tool{
		Name:  "git",
		MacOS: &domain.PackageRef{Formula: "git"},
	}
	svc := service.NewInstallService(
		&testutil.FakePackageManager{InstalledTools: map[string]bool{}},
		&testutil.FakeSystemChecker{},
		nil,
		false,
		io.Discard,
	)
	plan, err := svc.Plan(context.Background(), makeSelection(tool))
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
	}
	if len(plan.Steps) != 1 {
		t.Fatalf("expected 1 step, got %d", len(plan.Steps))
	}
	if plan.Steps[0].Status != domain.StatusPending {
		t.Errorf("status = %q, want %q", plan.Steps[0].Status, domain.StatusPending)
	}
}

func TestInstallService_Plan_SkipsToolWithNoRef(t *testing.T) {
	// Tool has no macOS ref — should not appear in plan.
	tool := domain.Tool{
		Name:    "WinOnlyTool",
		Windows: &domain.PackageRef{Winget: "some.tool"},
	}
	svc := service.NewInstallService(
		&testutil.FakePackageManager{},
		&testutil.FakeSystemChecker{},
		nil,
		false,
		io.Discard,
	)
	plan, err := svc.Plan(context.Background(), makeSelection(tool))
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
	}
	if len(plan.Steps) != 0 {
		t.Errorf("expected 0 steps, got %d", len(plan.Steps))
	}
}

// ── Execute ─────────────────────────────────────────────────────────────────

func TestInstallService_Execute_InstallsPendingSteps(t *testing.T) {
	tools := []domain.Tool{
		{Name: "git", MacOS: &domain.PackageRef{Formula: "git"}},
		{Name: "signal", MacOS: &domain.PackageRef{Cask: "signal"}},
	}
	fake := &testutil.FakePackageManager{}
	svc := service.NewInstallService(fake, &testutil.FakeSystemChecker{}, nil, false, io.Discard)

	plan, err := svc.Plan(context.Background(), makeSelection(tools...))
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
	}

	progress := make(chan domain.ProgressEvent, 20)
	if err := svc.Execute(context.Background(), plan, progress); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	close(progress)

	if len(fake.InstallCalls) != 2 {
		t.Errorf("expected 2 install calls, got %d", len(fake.InstallCalls))
	}
}

func TestInstallService_Execute_EmitsProgressEvents(t *testing.T) {
	tool := domain.Tool{Name: "git", MacOS: &domain.PackageRef{Formula: "git"}}
	svc := service.NewInstallService(
		&testutil.FakePackageManager{},
		&testutil.FakeSystemChecker{},
		nil,
		false,
		io.Discard,
	)
	plan, _ := svc.Plan(context.Background(), makeSelection(tool))

	progress := make(chan domain.ProgressEvent, 10)
	if err := svc.Execute(context.Background(), plan, progress); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	close(progress)

	var events []domain.ProgressEvent
	for ev := range progress {
		events = append(events, ev)
	}

	// Expect 2 events per step: "installing" then "done"
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].Step.Status != domain.StatusInstalling {
		t.Errorf("first event status = %q, want %q", events[0].Step.Status, domain.StatusInstalling)
	}
	if events[1].Step.Status != domain.StatusDone {
		t.Errorf("second event status = %q, want %q", events[1].Step.Status, domain.StatusDone)
	}
}

func TestInstallService_Execute_DryRun(t *testing.T) {
	tool := domain.Tool{Name: "signal", MacOS: &domain.PackageRef{Cask: "signal"}}
	fake := &testutil.FakePackageManager{}
	svc := service.NewInstallService(fake, &testutil.FakeSystemChecker{}, nil, true, io.Discard)
	plan, _ := svc.Plan(context.Background(), makeSelection(tool))

	if err := svc.Execute(context.Background(), plan, nil); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	// In dry-run, Install is still called (via DryRunRunner in production; fake here)
	if len(fake.InstallCalls) != 1 {
		t.Errorf("expected 1 install call, got %d", len(fake.InstallCalls))
	}
}

func TestInstallService_Execute_Resume_SkipsCompleted(t *testing.T) {
	tools := []domain.Tool{
		{Name: "git", MacOS: &domain.PackageRef{Formula: "git"}},
		{Name: "signal", MacOS: &domain.PackageRef{Cask: "signal"}},
	}
	fake := &testutil.FakePackageManager{}
	store := &testutil.FakeStateStore{
		State: &domain.InstallState{CompletedIDs: []string{"git"}},
	}
	svc := service.NewInstallService(fake, &testutil.FakeSystemChecker{}, store, false, io.Discard)

	plan, err := svc.Plan(context.Background(), makeSelection(tools...))
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
	}

	if err := svc.Execute(context.Background(), plan, nil); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	// Only "signal" should be installed; "git" was already completed.
	if len(fake.InstallCalls) != 1 {
		t.Errorf("expected 1 install call (skipping completed), got %d", len(fake.InstallCalls))
	}
	if len(fake.InstallCalls) == 1 && fake.InstallCalls[0].Cask != "signal" {
		t.Errorf("expected install for signal, got %+v", fake.InstallCalls[0])
	}
}

func TestInstallService_Execute_StateNotClearedOnFailure(t *testing.T) {
	tools := []domain.Tool{
		{Name: "git", MacOS: &domain.PackageRef{Formula: "git"}},
		{Name: "signal", MacOS: &domain.PackageRef{Cask: "signal"}},
	}

	// First run: second tool fails — state must NOT be cleared.
	failOnSecond := 0
	fake := &testutil.FakePackageManager{
		InstallErr: nil, // will be overridden via custom logic
	}
	_ = failOnSecond
	fakeWithErr := &testutil.FakePackageManager{}
	store := &testutil.FakeStateStore{}
	svc := service.NewInstallService(fakeWithErr, &testutil.FakeSystemChecker{}, store, false, io.Discard)

	plan, _ := svc.Plan(context.Background(), makeSelection(tools...))

	// Simulate second tool failure by setting InstallErr before Execute.
	fakeWithErr.InstallErr = errTest

	_ = svc.Execute(context.Background(), plan, nil)

	if store.ClearCalled {
		t.Error("expected ClearCalled = false when install fails")
	}

	// Second run: both succeed — state MUST be cleared.
	store2 := &testutil.FakeStateStore{}
	fakeOK := &testutil.FakePackageManager{}
	svc2 := service.NewInstallService(fakeOK, &testutil.FakeSystemChecker{}, store2, false, io.Discard)
	plan2, _ := svc2.Plan(context.Background(), makeSelection(tools...))

	if err := svc2.Execute(context.Background(), plan2, nil); err != nil {
		t.Fatalf("Execute() error on success run = %v", err)
	}
	if !store2.ClearCalled {
		t.Error("expected ClearCalled = true when all installs succeed")
	}
	_ = fake
}

func TestInstallService_Execute_DryRun_NoStateStore(t *testing.T) {
	tool := domain.Tool{Name: "signal", MacOS: &domain.PackageRef{Cask: "signal"}}
	fake := &testutil.FakePackageManager{}
	store := &testutil.FakeStateStore{}
	svc := service.NewInstallService(fake, &testutil.FakeSystemChecker{}, store, true, io.Discard)

	plan, _ := svc.Plan(context.Background(), makeSelection(tool))

	if err := svc.Execute(context.Background(), plan, nil); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if store.LoadCalled {
		t.Error("expected LoadCalled = false in dry-run")
	}
	if store.SaveCalled {
		t.Error("expected SaveCalled = false in dry-run")
	}
	if store.ClearCalled {
		t.Error("expected ClearCalled = false in dry-run")
	}
}

// ── helpers ──────────────────────────────────────────────────────────────────

var errTest = errors.New("test error")
