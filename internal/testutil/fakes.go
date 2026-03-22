// Package testutil provides shared fakes and builders for unit tests.
package testutil

import (
	"context"
	"fmt"
	"strings"

	"github.com/uttejg/newbox/internal/core/domain"
	"github.com/uttejg/newbox/internal/core/port"
)

// RunCall records a single call to FakeRunner.Run.
type RunCall struct {
	Cmd  string
	Args []string
}

// FakeRunner records Run calls and returns preconfigured results.
type FakeRunner struct {
	Results []*port.RunResult
	Err     error
	Calls   []RunCall
}

func (f *FakeRunner) Run(_ context.Context, cmd string, args []string) (*port.RunResult, error) {
	f.Calls = append(f.Calls, RunCall{Cmd: cmd, Args: args})
	if f.Err != nil {
		return nil, f.Err
	}
	if len(f.Results) > 0 {
		r := f.Results[0]
		f.Results = f.Results[1:]
		return r, nil
	}
	return &port.RunResult{Command: cmd + " " + strings.Join(args, " ")}, nil
}

// FakePackageManager is a test double for port.PackageManager.
type FakePackageManager struct {
	AvailableResult bool
	InstalledTools  map[string]bool
	InstallErr      error
	InstallCalls    []domain.PackageRef
}

func (f *FakePackageManager) Name() string { return "fake" }

func (f *FakePackageManager) IsAvailable(_ context.Context) bool { return f.AvailableResult }

func (f *FakePackageManager) IsInstalled(_ context.Context, ref domain.PackageRef) (bool, error) {
	key := ref.Formula
	if ref.Cask != "" {
		key = ref.Cask
	}
	return f.InstalledTools[key], nil
}

func (f *FakePackageManager) BuildCommand(ref domain.PackageRef) string {
	if ref.Cask != "" {
		return "brew install --cask " + ref.Cask
	}
	if ref.Formula != "" {
		return "brew install " + ref.Formula
	}
	return ""
}

func (f *FakePackageManager) Install(_ context.Context, ref domain.PackageRef) (*port.RunResult, error) {
	f.InstallCalls = append(f.InstallCalls, ref)
	return &port.RunResult{}, f.InstallErr
}

// FakeSystemChecker is a test double for port.SystemChecker.
type FakeSystemChecker struct {
	InternetErr error
	DiskErr     error
	PkgMgrErr   error
	SudoErr     error
}

func (f *FakeSystemChecker) CheckInternet(_ context.Context) error { return f.InternetErr }

func (f *FakeSystemChecker) CheckDiskSpace(_ context.Context, _ int) error { return f.DiskErr }

func (f *FakeSystemChecker) CheckPackageManager(_ context.Context, _ string) error {
	return f.PkgMgrErr
}

func (f *FakeSystemChecker) CheckSudo(_ context.Context) error { return f.SudoErr }

// FakeCatalogProvider implements port.CatalogProvider with in-memory data.
type FakeCatalogProvider struct {
	Categories []domain.Category
	Profiles   []domain.Profile
	Err        error // if non-nil, all calls return this error
}

func (f *FakeCatalogProvider) LoadCategories() ([]domain.Category, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	return f.Categories, nil
}

func (f *FakeCatalogProvider) LoadProfiles() ([]domain.Profile, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	return f.Profiles, nil
}

// FakeStateStore is a test double for port.StateStore.
type FakeStateStore struct {
	State   *domain.InstallState
	SaveErr error
	LoadErr error
}

func (f *FakeStateStore) Save(state *domain.InstallState) error {
	if f.SaveErr != nil {
		return f.SaveErr
	}
	f.State = state
	return nil
}

func (f *FakeStateStore) Load() (*domain.InstallState, error) {
	return f.State, f.LoadErr
}

func (f *FakeStateStore) Clear() error {
	f.State = nil
	return nil
}

func (f *FakeStateStore) Exists() bool {
	return f.State != nil
}

// NewTestTool creates a Tool available on the given OSes for testing.
func NewTestTool(name string, oses ...domain.OS) domain.Tool {
	t := domain.Tool{Name: name, Description: "test tool " + name}
	for _, os := range oses {
		ref := &domain.PackageRef{Formula: "test-" + name}
		switch os {
		case domain.OSMacOS:
			t.MacOS = ref
		case domain.OSLinux:
			t.Linux = ref
		case domain.OSWindows:
			t.Windows = ref
		}
	}
	return t
}

// NewTestCategory creates a Category with the given tools.
func NewTestCategory(id, name string, tools ...domain.Tool) domain.Category {
	return domain.Category{
		ID:          id,
		Name:        name,
		Description: "test category " + id,
		Tools:       tools,
	}
}

// NewTestProfile creates a Profile for testing.
func NewTestProfile(id string, categories ...string) domain.Profile {
	return domain.Profile{
		ID:          id,
		Name:        fmt.Sprintf("Test Profile %s", id),
		Description: "test profile",
		Categories:  categories,
	}
}

// DefaultTestCatalog returns a small catalog useful in most tests.
func DefaultTestCatalog() ([]domain.Category, []domain.Profile) {
	cats := []domain.Category{
		NewTestCategory("messaging", "💬 Messaging",
			NewTestTool("Signal", domain.OSMacOS, domain.OSLinux),
			NewTestTool("Telegram", domain.OSMacOS, domain.OSLinux, domain.OSWindows),
			NewTestTool("WinOnly", domain.OSWindows),
		),
		NewTestCategory("browsers", "🌐 Browsers",
			NewTestTool("Firefox", domain.OSMacOS, domain.OSLinux, domain.OSWindows),
			NewTestTool("Arc", domain.OSMacOS), // macOS only
		),
		NewTestCategory("cli", "🔧 CLI Essentials",
			NewTestTool("git", domain.OSMacOS, domain.OSLinux, domain.OSWindows),
			NewTestTool("jq", domain.OSMacOS, domain.OSLinux, domain.OSWindows),
		),
	}
	profiles := []domain.Profile{
		NewTestProfile("developer", "messaging", "browsers", "cli"),
		NewTestProfile("minimal", "cli"),
		{ID: "full", Name: "🚀 Full", Description: "Everything", AllCategories: true},
		{ID: "custom", Name: "🔧 Custom", Description: "Pick your own"},
	}
	return cats, profiles
}
