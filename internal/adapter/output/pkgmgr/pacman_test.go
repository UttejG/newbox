package pkgmgr_test

import (
	"context"
	"testing"

	"github.com/uttejg/newbox/internal/adapter/output/pkgmgr"
	"github.com/uttejg/newbox/internal/core/domain"
	"github.com/uttejg/newbox/internal/core/port"
	"github.com/uttejg/newbox/internal/testutil"
)

func TestPacmanManager_Name(t *testing.T) {
	p := pkgmgr.NewPacman(&testutil.FakeRunner{})
	if got := p.Name(); got != "pacman" {
		t.Errorf("Name() = %q, want \"pacman\"", got)
	}
}

func TestPacmanManager_IsAvailable(t *testing.T) {
	tests := []struct {
		name    string
		results []*port.RunResult
		wantErr bool
	}{
		{
			name:    "pacman available",
			results: []*port.RunResult{{ExitCode: 0}},
			wantErr: false,
		},
		{
			name:    "pacman not available",
			results: []*port.RunResult{{ExitCode: 1}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &testutil.FakeRunner{Results: tt.results}
			p := pkgmgr.NewPacman(fake)
			got := p.IsAvailable(context.Background())
			if (got != nil) != tt.wantErr {
				t.Errorf("IsAvailable() error = %v, wantErr %v", got, tt.wantErr)
			}
			if len(fake.Calls) != 1 {
				t.Fatalf("expected 1 call, got %d", len(fake.Calls))
			}
			call := fake.Calls[0]
			if call.Cmd != "pacman" || call.Args[0] != "--version" {
				t.Errorf("unexpected call: %s %v", call.Cmd, call.Args)
			}
		})
	}
}

func TestPacmanManager_IsInstalled(t *testing.T) {
	tests := []struct {
		name    string
		ref     domain.PackageRef
		results []*port.RunResult
		want    bool
	}{
		{
			name:    "installed (exit 0)",
			ref:     domain.PackageRef{Pacman: "firefox"},
			results: []*port.RunResult{{ExitCode: 0}},
			want:    true,
		},
		{
			name:    "not installed (exit 1)",
			ref:     domain.PackageRef{Pacman: "firefox"},
			results: []*port.RunResult{{ExitCode: 1}},
			want:    false,
		},
		{
			name: "empty ref returns false",
			ref:  domain.PackageRef{},
			want: false,
		},
		{
			name:    "dry-run treated as not installed",
			ref:     domain.PackageRef{Pacman: "firefox"},
			results: []*port.RunResult{{ExitCode: 0, DryRun: true}},
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &testutil.FakeRunner{Results: tt.results}
			p := pkgmgr.NewPacman(fake)
			got, err := p.IsInstalled(context.Background(), tt.ref)
			if err != nil {
				t.Fatalf("IsInstalled() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("IsInstalled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacmanManager_Install(t *testing.T) {
	tests := []struct {
		name     string
		ref      domain.PackageRef
		wantArgs []string
		wantErr  bool
	}{
		{
			name:     "install package",
			ref:      domain.PackageRef{Pacman: "firefox"},
			wantArgs: []string{"-S", "--noconfirm", "firefox"},
		},
		{
			name:    "empty ref returns error",
			ref:     domain.PackageRef{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &testutil.FakeRunner{}
			p := pkgmgr.NewPacman(fake)
			_, err := p.Install(context.Background(), tt.ref)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
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
			if call.Cmd != "pacman" {
				t.Errorf("cmd = %q, want \"pacman\"", call.Cmd)
			}
			for i, a := range tt.wantArgs {
				if call.Args[i] != a {
					t.Errorf("args[%d] = %q, want %q", i, call.Args[i], a)
				}
			}
		})
	}
}
