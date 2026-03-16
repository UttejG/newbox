package pkgmgr

import (
	"github.com/uttejg/newbox/internal/core/domain"
	"github.com/uttejg/newbox/internal/core/port"
)

// NewForPlatform returns the appropriate PackageManager for the given platform.
// On macOS it returns a composite of brew+mas.
// On Linux it returns the native manager (apt/dnf/pacman) plus flatpak as fallback.
func NewForPlatform(platform *domain.Platform, runner port.CommandRunner) port.PackageManager {
	switch platform.OS {
	case domain.OSMacOS:
		return NewComposite(NewBrew(runner), NewMAS(runner))

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
		managers = append(managers, NewFlatpak(runner))
		return NewComposite(managers...)

	default:
		return NewBrew(runner)
	}
}
