package pkgmgr_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/uttejg/newbox/internal/adapter/output/pkgmgr"
	"github.com/uttejg/newbox/internal/core/domain"
	"github.com/uttejg/newbox/internal/core/port"
	"github.com/uttejg/newbox/internal/testutil"
)

func TestWingetManager_Name(t *testing.T) {
	w := pkgmgr.NewWinget(&testutil.FakeRunner{})
	if w.Name() != "winget" {
		t.Errorf("expected winget, got %s", w.Name())
	}
}

func TestWingetManager_IsAvailable(t *testing.T) {
	tests := []struct {
		name    string
		results []*port.RunResult
		err     error
		want    bool
	}{
		{"available", []*port.RunResult{{ExitCode: 0}}, nil, true},
		{"not available", []*port.RunResult{{ExitCode: 1}}, nil, false},
		{"error", nil, fmt.Errorf("not found"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &testutil.FakeRunner{Results: tt.results, Err: tt.err}
			w := pkgmgr.NewWinget(fake)
			got := w.IsAvailable(context.Background())
			if got != tt.want {
				t.Errorf("IsAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWingetManager_IsInstalled(t *testing.T) {
	tests := []struct {
		name     string
		ref      domain.PackageRef
		results  []*port.RunResult
		want     bool
		wantArgs []string
	}{
		{
			name:     "installed",
			ref:      domain.PackageRef{Winget: "Mozilla.Firefox"},
			results:  []*port.RunResult{{ExitCode: 0, Stdout: "Mozilla.Firefox  Firefox  120.0"}},
			want:     true,
			wantArgs: []string{"list", "--id", "Mozilla.Firefox", "--exact"},
		},
		{
			name:     "not installed - non-zero exit",
			ref:      domain.PackageRef{Winget: "Mozilla.Firefox"},
			results:  []*port.RunResult{{ExitCode: 1, Stdout: "No installed package found"}},
			want:     false,
			wantArgs: []string{"list", "--id", "Mozilla.Firefox", "--exact"},
		},
		{
			name:     "not installed - id absent from stdout",
			ref:      domain.PackageRef{Winget: "Mozilla.Firefox"},
			results:  []*port.RunResult{{ExitCode: 0, Stdout: "some other output"}},
			want:     false,
			wantArgs: []string{"list", "--id", "Mozilla.Firefox", "--exact"},
		},
		{
			name: "no winget ID",
			ref:  domain.PackageRef{},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &testutil.FakeRunner{Results: tt.results}
			w := pkgmgr.NewWinget(fake)
			got, err := w.IsInstalled(context.Background(), tt.ref)
			if err != nil {
				t.Fatalf("IsInstalled() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("IsInstalled() = %v, want %v", got, tt.want)
			}
			if tt.wantArgs != nil {
				if len(fake.Calls) != 1 {
					t.Fatalf("expected 1 call, got %d", len(fake.Calls))
				}
				call := fake.Calls[0]
				if call.Cmd != "winget" {
					t.Errorf("cmd = %q, want \"winget\"", call.Cmd)
				}
				for i, a := range tt.wantArgs {
					if call.Args[i] != a {
						t.Errorf("args[%d] = %q, want %q", i, call.Args[i], a)
					}
				}
			}
		})
	}
}

func TestWingetManager_Install(t *testing.T) {
	tests := []struct {
		name     string
		ref      domain.PackageRef
		wantArgs []string
		wantErr  bool
	}{
		{
			name: "install package",
			ref:  domain.PackageRef{Winget: "Mozilla.Firefox"},
			wantArgs: []string{
				"install", "--id", "Mozilla.Firefox",
				"--exact", "--silent",
				"--accept-package-agreements",
			},
		},
		{
			name:    "no winget ID returns error",
			ref:     domain.PackageRef{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &testutil.FakeRunner{}
			w := pkgmgr.NewWinget(fake)
			res, err := w.Install(context.Background(), tt.ref)
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
			if call.Cmd != "winget" {
				t.Errorf("cmd = %q, want \"winget\"", call.Cmd)
			}
			for i, a := range tt.wantArgs {
				if call.Args[i] != a {
					t.Errorf("args[%d] = %q, want %q", i, call.Args[i], a)
				}
			}
		})
	}
}
