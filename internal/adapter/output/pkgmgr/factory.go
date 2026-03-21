package pkgmgr

import (
	"os/exec"

	"github.com/uttejg/newbox/internal/core/domain"
	"github.com/uttejg/newbox/internal/core/port"
)

// NewForPlatform returns the appropriate PackageManager for the given platform.
// On macOS it returns a composite of brew+mas.
// On Linux it returns the detected native manager (apt/dnf/pacman) and, if flatpak
// is installed on the system, adds it as a fallback for cross-distro packages.
func NewForPlatform(platform *domain.Platform, runner port.CommandRunner) port.PackageManager {
	switch platform.OS {
	case domain.OSMacOS:
		return NewComposite(NewBrew(runner), NewMAS(runner))
	case domain.OSWindows:
		return NewWinget(runner)
	case domain.OSLinux:
		var managers []port.PackageManager
		switch platform.PackageManager {
		case domain.PkgMgrApt:
			managers = append(managers, NewApt(runner))
		case domain.PkgMgrDnf:
			managers = append(managers, NewDnf(runner))
		case domain.PkgMgrPacman:
			managers = append(managers, NewPacman(runner))
		}
		// Only include flatpak if it is actually installed on the system.
		// Flatpak is a fallback; its absence should not block preflight.
		if _, err := exec.LookPath("flatpak"); err == nil {
			managers = append(managers, NewFlatpak(runner))
		}
		return NewComposite(managers...)

	default:
		return NewBrew(runner)
	}
}
