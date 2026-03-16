package pkgmgr_test

import (
	"context"
	"errors"
	"testing"

	"github.com/uttejg/newbox/internal/adapter/output/pkgmgr"
	"github.com/uttejg/newbox/internal/core/domain"
	"github.com/uttejg/newbox/internal/core/port"
	"github.com/uttejg/newbox/internal/testutil"
)

func TestMASManager_Name(t *testing.T) {
	m := pkgmgr.NewMAS(&testutil.FakeRunner{})
	if got := m.Name(); got != "mas" {
		t.Errorf("Name() = %q, want \"mas\"", got)
	}
}

func TestMASManager_IsAvailable(t *testing.T) {
	tests := []struct {
		name    string
		results []*port.RunResult
		wantErr error
		want    bool
	}{
		{
			name:    "mas available",
			results: []*port.RunResult{{ExitCode: 0}},
			want:    true,
		},
		{
			name:    "mas not available (exit code 1)",
			results: []*port.RunResult{{ExitCode: 1}},
			want:    false,
		},
		{
			name:    "mas not available (run error)",
			wantErr: errors.New("command not found"),
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &testutil.FakeRunner{Results: tt.results, Err: tt.wantErr}
			m := pkgmgr.NewMAS(fake)
			got := m.IsAvailable(context.Background())
			if got != tt.want {
				t.Errorf("IsAvailable() = %v, want %v", got, tt.want)
			}
			if len(fake.Calls) != 1 {
				t.Fatalf("expected 1 call, got %d", len(fake.Calls))
			}
			if fake.Calls[0].Cmd != "mas" || fake.Calls[0].Args[0] != "version" {
				t.Errorf("unexpected call: %s %v", fake.Calls[0].Cmd, fake.Calls[0].Args)
			}
		})
	}
}

func TestMASManager_IsInstalled(t *testing.T) {
	tests := []struct {
		name    string
		ref     domain.PackageRef
		results []*port.RunResult
		want    bool
		wantErr bool
	}{
		{
			name:    "app found in list",
			ref:     domain.PackageRef{MAS: "409203825"},
			results: []*port.RunResult{{ExitCode: 0, Stdout: "409203825 Xcode (14.0)\n"}},
			want:    true,
		},
		{
			name:    "app not in list",
			ref:     domain.PackageRef{MAS: "409203825"},
			results: []*port.RunResult{{ExitCode: 0, Stdout: "123456789 OtherApp (1.0)\n"}},
			want:    false,
		},
		{
			name: "empty MAS ID returns false without call",
			ref:  domain.PackageRef{},
			want: false,
		},
		{
			name:    "dry-run result treated as not installed",
			ref:     domain.PackageRef{MAS: "409203825"},
			results: []*port.RunResult{{ExitCode: 0, DryRun: true, Stdout: "[dry-run]"}},
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &testutil.FakeRunner{Results: tt.results}
			m := pkgmgr.NewMAS(fake)
			got, err := m.IsInstalled(context.Background(), tt.ref)
			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("IsInstalled() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("IsInstalled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMASManager_Install(t *testing.T) {
	tests := []struct {
		name     string
		ref      domain.PackageRef
		wantArgs []string
		wantErr  bool
	}{
		{
			name:     "installs by MAS ID",
			ref:      domain.PackageRef{MAS: "409203825"},
			wantArgs: []string{"install", "409203825"},
		},
		{
			name:    "empty MAS ID returns error",
			ref:     domain.PackageRef{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &testutil.FakeRunner{}
			m := pkgmgr.NewMAS(fake)
			res, err := m.Install(context.Background(), tt.ref)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("Install() error = %v", err)
			}
			if res == nil {
				t.Fatal("expected non-nil result")
			}
			if len(fake.Calls) != 1 {
				t.Fatalf("expected 1 call, got %d", len(fake.Calls))
			}
			call := fake.Calls[0]
			if call.Cmd != "mas" {
				t.Errorf("cmd = %q, want \"mas\"", call.Cmd)
			}
			for i, a := range tt.wantArgs {
				if call.Args[i] != a {
					t.Errorf("args[%d] = %q, want %q", i, call.Args[i], a)
				}
			}
		})
	}
}
