package domain_test

import (
	"testing"

	"github.com/uttejg/newbox/internal/core/domain"
)

func TestOS_String(t *testing.T) {
	tests := []struct {
		name string
		os   domain.OS
		want string
	}{
		{"macOS", domain.OSMacOS, "macOS"},
		{"Linux", domain.OSLinux, "Linux"},
		{"Windows", domain.OSWindows, "Windows"},
		{"Unknown", domain.OSUnknown, "Unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.os.String(); got != tt.want {
				t.Errorf("OS.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestArch_String(t *testing.T) {
	tests := []struct {
		name string
		arch domain.Arch
		want string
	}{
		{"amd64", domain.ArchAMD64, "amd64"},
		{"arm64", domain.ArchARM64, "arm64"},
		{"unknown", domain.ArchUnknown, "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arch.String(); got != tt.want {
				t.Errorf("Arch.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDistro_String(t *testing.T) {
	tests := []struct {
		name   string
		distro domain.Distro
		want   string
	}{
		{"Debian", domain.DistroDebian, "Debian"},
		{"Ubuntu", domain.DistroUbuntu, "Ubuntu"},
		{"Fedora", domain.DistroFedora, "Fedora"},
		{"Arch", domain.DistroArch, "Arch"},
		{"None", domain.DistroNone, ""},
		{"Unknown", domain.DistroUnknown, "Unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.distro.String(); got != tt.want {
				t.Errorf("Distro.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPackageManagerType_String(t *testing.T) {
	tests := []struct {
		name string
		pm   domain.PackageManagerType
		want string
	}{
		{"brew", domain.PkgMgrBrew, "brew"},
		{"winget", domain.PkgMgrWinget, "winget"},
		{"apt", domain.PkgMgrApt, "apt"},
		{"dnf", domain.PkgMgrDnf, "dnf"},
		{"pacman", domain.PkgMgrPacman, "pacman"},
		{"none", domain.PkgMgrNone, "none"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pm.String(); got != tt.want {
				t.Errorf("PackageManagerType.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPlatform_Summary(t *testing.T) {
	tests := []struct {
		name     string
		platform domain.Platform
		want     string
	}{
		{
			name: "macOS ARM64",
			platform: domain.Platform{
				OS:   domain.OSMacOS,
				Arch: domain.ArchARM64,
			},
			want: "macOS arm64",
		},
		{
			name: "Linux Ubuntu amd64",
			platform: domain.Platform{
				OS:     domain.OSLinux,
				Arch:   domain.ArchAMD64,
				Distro: domain.DistroUbuntu,
			},
			want: "Linux (Ubuntu) amd64",
		},
		{
			name: "Windows amd64",
			platform: domain.Platform{
				OS:   domain.OSWindows,
				Arch: domain.ArchAMD64,
			},
			want: "Windows amd64",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.platform.Summary(); got != tt.want {
				t.Errorf("Platform.Summary() = %q, want %q", got, tt.want)
			}
		})
	}
}
