package runner_test

import (
	"context"
	"strings"
	"testing"

	"github.com/uttejg/newbox/internal/adapter/output/runner"
)

func TestDryRunRunner(t *testing.T) {
	tests := []struct {
		name    string
		cmd     string
		args    []string
		wantLog string
		wantLen int
	}{
		{
			name:    "brew install formula",
			cmd:     "brew",
			args:    []string{"install", "git"},
			wantLog: "brew install git",
			wantLen: 1,
		},
		{
			name:    "brew install cask",
			cmd:     "brew",
			args:    []string{"install", "--cask", "signal"},
			wantLog: "brew install --cask signal",
			wantLen: 1,
		},
		{
			name:    "no args",
			cmd:     "brew",
			args:    []string{"--version"},
			wantLog: "brew --version",
			wantLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &runner.DryRunRunner{}
			res, err := r.Run(context.Background(), tt.cmd, tt.args)
			if err != nil {
				t.Fatalf("Run() error = %v", err)
			}
			if !res.DryRun {
				t.Error("expected DryRun=true")
			}
			if res.ExitCode != 0 {
				t.Errorf("expected ExitCode=0, got %d", res.ExitCode)
			}
			if !strings.Contains(res.Stdout, tt.wantLog) {
				t.Errorf("Stdout %q does not contain %q", res.Stdout, tt.wantLog)
			}
			if res.Command != tt.wantLog {
				t.Errorf("Command = %q, want %q", res.Command, tt.wantLog)
			}
			if len(r.Log) != tt.wantLen {
				t.Errorf("Log length = %d, want %d", len(r.Log), tt.wantLen)
			}
		})
	}
}

func TestDryRunRunner_AccumulatesLog(t *testing.T) {
	r := &runner.DryRunRunner{}
	ctx := context.Background()

	cmds := []struct {
		cmd  string
		args []string
	}{
		{"brew", []string{"install", "git"}},
		{"brew", []string{"install", "--cask", "signal"}},
		{"brew", []string{"--version"}},
	}

	for _, c := range cmds {
		if _, err := r.Run(ctx, c.cmd, c.args); err != nil {
			t.Fatalf("Run() error = %v", err)
		}
	}

	if len(r.Log) != len(cmds) {
		t.Errorf("Log length = %d, want %d", len(r.Log), len(cmds))
	}
}
