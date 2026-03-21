package pkgmgr_test

import (
	"context"
	"testing"

	"github.com/uttejg/newbox/internal/adapter/output/pkgmgr"
	"github.com/uttejg/newbox/internal/core/domain"
	"github.com/uttejg/newbox/internal/core/port"
	"github.com/uttejg/newbox/internal/testutil"
)

func TestFlatpakManager_Name(t *testing.T) {
	f := pkgmgr.NewFlatpak(&testutil.FakeRunner{})
	if got := f.Name(); got != "flatpak" {
		t.Errorf("Name() = %q, want \"flatpak\"", got)
	}
}

func TestFlatpakManager_IsAvailable(t *testing.T) {
	tests := []struct {
		name    string
		results []*port.RunResult
		wantErr bool
	}{
		{
			name:    "flatpak available",
			results: []*port.RunResult{{ExitCode: 0}},
			wantErr: false,
		},
		{
			name:    "flatpak not available",
			results: []*port.RunResult{{ExitCode: 1}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &testutil.FakeRunner{Results: tt.results}
			f := pkgmgr.NewFlatpak(fake)
			got := f.IsAvailable(context.Background())
			if (got != nil) != tt.wantErr {
				t.Errorf("IsAvailable() error = %v, wantErr %v", got, tt.wantErr)
			}
			if len(fake.Calls) != 1 {
				t.Fatalf("expected 1 call, got %d", len(fake.Calls))
			}
			call := fake.Calls[0]
			if call.Cmd != "flatpak" || call.Args[0] != "--version" {
				t.Errorf("unexpected call: %s %v", call.Cmd, call.Args)
			}
		})
	}
}

func TestFlatpakManager_IsInstalled(t *testing.T) {
	tests := []struct {
		name    string
		ref     domain.PackageRef
		results []*port.RunResult
		want    bool
	}{
		{
			name:    "installed",
			ref:     domain.PackageRef{Flatpak: "org.signal.Signal"},
			results: []*port.RunResult{{ExitCode: 0, Stdout: "org.signal.Signal\n"}},
			want:    true,
		},
		{
			name:    "not installed",
			ref:     domain.PackageRef{Flatpak: "org.signal.Signal"},
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
			ref:     domain.PackageRef{Flatpak: "org.signal.Signal"},
			results: []*port.RunResult{{ExitCode: 0, DryRun: true}},
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &testutil.FakeRunner{Results: tt.results}
			f := pkgmgr.NewFlatpak(fake)
			got, err := f.IsInstalled(context.Background(), tt.ref)
			if err != nil {
				t.Fatalf("IsInstalled() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("IsInstalled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlatpakManager_Install(t *testing.T) {
	tests := []struct {
		name     string
		ref      domain.PackageRef
		wantArgs []string
		wantErr  bool
	}{
		{
			name:     "install package",
			ref:      domain.PackageRef{Flatpak: "org.signal.Signal"},
			wantArgs: []string{"install", "-y", "flathub", "org.signal.Signal"},
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
			f := pkgmgr.NewFlatpak(fake)
			_, err := f.Install(context.Background(), tt.ref)
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
			if call.Cmd != "flatpak" {
				t.Errorf("cmd = %q, want \"flatpak\"", call.Cmd)
			}
			for i, a := range tt.wantArgs {
				if call.Args[i] != a {
					t.Errorf("args[%d] = %q, want %q", i, call.Args[i], a)
				}
			}
		})
	}
}
