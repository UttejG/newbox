//go:build integration && linux

package pkgmgr_test

import (
	"context"
	"testing"

	"github.com/uttejg/newbox/internal/adapter/output/pkgmgr"
	"github.com/uttejg/newbox/internal/adapter/output/runner"
)

func TestApt_Integration_IsAvailable(t *testing.T) {
	apt := pkgmgr.NewApt(&runner.ExecRunner{})
	if err := apt.IsAvailable(context.Background()); err != nil {
		t.Skip("apt not available")
	}
}
