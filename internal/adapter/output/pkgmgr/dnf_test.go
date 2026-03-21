package pkgmgr_test

import (
	"context"
	"testing"

	"github.com/uttejg/newbox/internal/adapter/output/pkgmgr"
	"github.com/uttejg/newbox/internal/core/domain"
	"github.com/uttejg/newbox/internal/core/port"
	"github.com/uttejg/newbox/internal/testutil"
)

func TestDnfManager_Name(t *testing.T) {
	d := pkgmgr.NewDnf(&testutil.FakeRunner{})
	if got := d.Name(); got != "dnf" {
		t.Errorf("Name() = %q, want \"dnf\"", got)
	}
}

func TestDnfManager_IsAvailable(t *testing.T) {
	tests := []struct {
		name    string
		results []*port.RunResult
		wantErr bool
	}{
		{
			name:    "dnf available",
			results: []*port.RunResult{{ExitCode: 0}},
			wantErr: false,
		},
		{
			name:    "dnf not available",
			results: []*port.RunResult{{ExitCode: 1}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &testutil.FakeRunner{Results: tt.results}
			d := pkgmgr.NewDnf(fake)
			got := d.IsAvailable(context.Background())
			if (got != nil) != tt.wantErr {
				t.Errorf("IsAvailable() error = %v, wantErr %v", got, tt.wantErr)
			}
			if len(fake.Calls) != 1 {
				t.Fatalf("expected 1 call, got %d", len(fake.Calls))
			}
			call := fake.Calls[0]
			if call.Cmd != "dnf" || call.Args[0] != "--version" {
				t.Errorf("unexpected call: %s %v", call.Cmd, call.Args)
			}
		})
	}
}

func TestDnfManager_IsInstalled(t *testing.T) {
	tests := []struct {
		name    string
		ref     domain.PackageRef
		results []*port.RunResult
		want    bool
	}{
		{
			name:    "installed",
			ref:     domain.PackageRef{Dnf: "firefox"},
			results: []*port.RunResult{{ExitCode: 0, Stdout: "firefox.x86_64  118.0-1.fc39"}},
			want:    true,
		},
		{
			name:    "not installed",
			ref:     domain.PackageRef{Dnf: "firefox"},
			results: []*port.RunResult{{ExitCode: 0, Stdout: ""}},
			want:    false,
		},
		{
			name: "empty ref returns false",
			ref:  domain.PackageRef{},
			want: false,
		},
		{
			name:    "dry-run treated as not installed",
			ref:     domain.PackageRef{Dnf: "firefox"},
			results: []*port.RunResult{{ExitCode: 0, DryRun: true}},
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &testutil.FakeRunner{Results: tt.results}
			d := pkgmgr.NewDnf(fake)
			got, err := d.IsInstalled(context.Background(), tt.ref)
			if err != nil {
				t.Fatalf("IsInstalled() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("IsInstalled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDnfManager_Install(t *testing.T) {
	tests := []struct {
		name     string
		ref      domain.PackageRef
		wantArgs []string
		wantErr  bool
	}{
		{
			name:     "install package",
			ref:      domain.PackageRef{Dnf: "firefox"},
			wantArgs: []string{"install", "-y", "firefox"},
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
			d := pkgmgr.NewDnf(fake)
			_, err := d.Install(context.Background(), tt.ref)
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
			if call.Cmd != "dnf" {
				t.Errorf("cmd = %q, want \"dnf\"", call.Cmd)
			}
			for i, a := range tt.wantArgs {
				if call.Args[i] != a {
					t.Errorf("args[%d] = %q, want %q", i, call.Args[i], a)
				}
			}
		})
	}
}
