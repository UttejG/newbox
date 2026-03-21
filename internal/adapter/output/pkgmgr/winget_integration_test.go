//go:build integration && windows

package pkgmgr_test

import (
	"context"
	"testing"

	"github.com/uttejg/newbox/internal/adapter/output/pkgmgr"
	"github.com/uttejg/newbox/internal/adapter/output/runner"
)

func TestWinget_Integration_IsAvailable(t *testing.T) {
	w := pkgmgr.NewWinget(&runner.ExecRunner{})
	if err := w.IsAvailable(context.Background()); err != nil {
		t.Skip("winget not available:", err)
	}
}
