package pkgmgr_test

import (
	"context"
	"testing"

	"github.com/uttejg/newbox/internal/adapter/output/pkgmgr"
	"github.com/uttejg/newbox/internal/core/domain"
	"github.com/uttejg/newbox/internal/core/port"
	"github.com/uttejg/newbox/internal/testutil"
)

func TestBrewManager_Name(t *testing.T) {
	b := pkgmgr.NewBrew(&testutil.FakeRunner{})
	if got := b.Name(); got != "brew" {
		t.Errorf("Name() = %q, want \"brew\"", got)
	}
}

func TestBrewManager_IsAvailable(t *testing.T) {
	tests := []struct {
		name    string
		results []*port.RunResult
		want    bool
	}{
		{
			name:    "brew available",
			results: []*port.RunResult{{ExitCode: 0}},
			want:    true,
		},
		{
			name:    "brew not available",
			results: []*port.RunResult{{ExitCode: 1}},
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &testutil.FakeRunner{Results: tt.results}
			b := pkgmgr.NewBrew(fake)
			got := b.IsAvailable(context.Background())
			if got != tt.want {
				t.Errorf("IsAvailable() = %v, want %v", got, tt.want)
			}
			if len(fake.Calls) != 1 {
				t.Fatalf("expected 1 call, got %d", len(fake.Calls))
			}
			call := fake.Calls[0]
			if call.Cmd != "brew" || call.Args[0] != "--version" {
				t.Errorf("unexpected call: %s %v", call.Cmd, call.Args)
			}
		})
	}
}

func TestBrewManager_IsInstalled(t *testing.T) {
	tests := []struct {
		name    string
		ref     domain.PackageRef
		results []*port.RunResult
		want    bool
		wantCmd []string // expected args
	}{
		{
			name:    "formula installed",
			ref:     domain.PackageRef{Formula: "git"},
			results: []*port.RunResult{{ExitCode: 0}},
			want:    true,
			wantCmd: []string{"list", "git"},
		},
		{
			name:    "formula not installed",
			ref:     domain.PackageRef{Formula: "git"},
			results: []*port.RunResult{{ExitCode: 1}},
			want:    false,
			wantCmd: []string{"list", "git"},
		},
		{
			name:    "cask installed",
			ref:     domain.PackageRef{Cask: "signal"},
			results: []*port.RunResult{{ExitCode: 0}},
			want:    true,
			wantCmd: []string{"list", "--cask", "signal"},
		},
		{
			name:    "cask not installed",
			ref:     domain.PackageRef{Cask: "signal"},
			results: []*port.RunResult{{ExitCode: 1}},
			want:    false,
			wantCmd: []string{"list", "--cask", "signal"},
		},
		{
			name: "dry-run result treated as not installed",
			ref:  domain.PackageRef{Formula: "git"},
			results: []*port.RunResult{{ExitCode: 0, DryRun: true,
				Stdout: "[dry-run] would execute: brew list git"}},
			want:    false,
			wantCmd: []string{"list", "git"},
		},
		{
			name: "empty ref returns false",
			ref:  domain.PackageRef{},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &testutil.FakeRunner{Results: tt.results}
			b := pkgmgr.NewBrew(fake)
			got, err := b.IsInstalled(context.Background(), tt.ref)
			if err != nil {
				t.Fatalf("IsInstalled() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("IsInstalled() = %v, want %v", got, tt.want)
			}
			if tt.wantCmd != nil {
				if len(fake.Calls) != 1 {
					t.Fatalf("expected 1 call, got %d", len(fake.Calls))
				}
				for i, a := range tt.wantCmd {
					if fake.Calls[0].Args[i] != a {
						t.Errorf("args[%d] = %q, want %q", i, fake.Calls[0].Args[i], a)
					}
				}
			}
		})
	}
}

func TestBrewManager_Install(t *testing.T) {
	tests := []struct {
		name     string
		ref      domain.PackageRef
		wantArgs []string
		wantNil  bool // nil result when ref is empty
	}{
		{
			name:     "install formula",
			ref:      domain.PackageRef{Formula: "git"},
			wantArgs: []string{"install", "git"},
		},
		{
			name:     "install cask",
			ref:      domain.PackageRef{Cask: "signal"},
			wantArgs: []string{"install", "--cask", "signal"},
		},
		{
			name:    "empty ref returns error",
			ref:     domain.PackageRef{},
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &testutil.FakeRunner{}
			b := pkgmgr.NewBrew(fake)
			res, err := b.Install(context.Background(), tt.ref)
			if tt.wantNil {
				if err == nil {
					t.Fatal("expected error for empty ref, got nil")
				}
				if res != nil {
					t.Error("expected nil result for empty ref")
				}
				return
			}
			if err != nil {
				t.Fatalf("Install() error = %v", err)
			}
			if len(fake.Calls) != 1 {
				t.Fatalf("expected 1 call, got %d", len(fake.Calls))
			}
			call := fake.Calls[0]
			if call.Cmd != "brew" {
				t.Errorf("cmd = %q, want \"brew\"", call.Cmd)
			}
			for i, a := range tt.wantArgs {
				if call.Args[i] != a {
					t.Errorf("args[%d] = %q, want %q", i, call.Args[i], a)
				}
			}
		})
	}
}
