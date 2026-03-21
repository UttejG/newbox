//go:build integration

package pkgmgr_test

import (
	"context"
	"testing"

	"github.com/uttejg/newbox/internal/adapter/output/pkgmgr"
	"github.com/uttejg/newbox/internal/adapter/output/runner"
	"github.com/uttejg/newbox/internal/core/domain"
)

func TestBrew_Integration_IsAvailable(t *testing.T) {
	brew := pkgmgr.NewBrew(&runner.ExecRunner{})
	if !brew.IsAvailable(context.Background()) {
		t.Skip("brew not available on this machine")
	}
}

func TestBrew_Integration_IsInstalled_Git(t *testing.T) {
	brew := pkgmgr.NewBrew(&runner.ExecRunner{})
	if !brew.IsAvailable(context.Background()) {
		t.Skip("brew not available")
	}
	installed, err := brew.IsInstalled(context.Background(), domain.PackageRef{Formula: "git"})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("git installed via brew: %v", installed)
}
