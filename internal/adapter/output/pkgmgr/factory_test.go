package pkgmgr_test

import (
	"testing"

	"github.com/uttejg/newbox/internal/adapter/output/pkgmgr"
	"github.com/uttejg/newbox/internal/core/domain"
	"github.com/uttejg/newbox/internal/testutil"
)

func TestNewForPlatform(t *testing.T) {
	tests := []struct {
		name     string
		platform *domain.Platform
		wantName string
	}{
		{
			name:     "macOS returns composite",
			platform: &domain.Platform{OS: domain.OSMacOS, PackageManager: domain.PkgMgrBrew},
			wantName: "composite",
		},
		{
			name:     "Linux apt returns composite",
			platform: &domain.Platform{OS: domain.OSLinux, PackageManager: domain.PkgMgrApt},
			wantName: "composite",
		},
		{
			name:     "Linux dnf returns composite",
			platform: &domain.Platform{OS: domain.OSLinux, PackageManager: domain.PkgMgrDnf},
			wantName: "composite",
		},
		{
			name:     "Linux pacman returns composite",
			platform: &domain.Platform{OS: domain.OSLinux, PackageManager: domain.PkgMgrPacman},
			wantName: "composite",
		},
		{
			name:     "Linux no native mgr returns composite (flatpak only)",
			platform: &domain.Platform{OS: domain.OSLinux, PackageManager: domain.PkgMgrNone},
			wantName: "composite",
		},
		{
			name:     "Windows returns winget",
			platform: &domain.Platform{OS: domain.OSWindows, PackageManager: domain.PkgMgrWinget},
			wantName: "winget",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mgr := pkgmgr.NewForPlatform(tt.platform, &testutil.FakeRunner{})
			if got := mgr.Name(); got != tt.wantName {
				t.Errorf("Name() = %q, want %q", got, tt.wantName)
			}
		})
	}
}
