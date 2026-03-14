package detector

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/uttejg/newbox/internal/core/domain"
)

// SystemDetector detects the platform using runtime info and system commands.
type SystemDetector struct{}

// Detect returns the detected platform information.
func (d *SystemDetector) Detect() (*domain.Platform, error) {
	p := &domain.Platform{
		OS:   detectOS(),
		Arch: detectArch(),
	}

	if p.OS == domain.OSLinux {
		p.Distro = detectDistro()
	}

	p.PackageManager = detectPackageManager(p.OS, p.Distro)

	return p, nil
}

func detectOS() domain.OS {
	switch runtime.GOOS {
	case "darwin":
		return domain.OSMacOS
	case "linux":
		return domain.OSLinux
	case "windows":
		return domain.OSWindows
	default:
		return domain.OSUnknown
	}
}

func detectArch() domain.Arch {
	switch runtime.GOARCH {
	case "amd64":
		return domain.ArchAMD64
	case "arm64":
		return domain.ArchARM64
	default:
		return domain.ArchUnknown
	}
}

func detectDistro() domain.Distro {
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return domain.DistroUnknown
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "ID=") {
			id := strings.TrimPrefix(line, "ID=")
			id = strings.Trim(id, "\"")
			switch strings.ToLower(id) {
			case "debian":
				return domain.DistroDebian
			case "ubuntu":
				return domain.DistroUbuntu
			case "fedora":
				return domain.DistroFedora
			case "arch":
				return domain.DistroArch
			default:
				return domain.DistroUnknown
			}
		}
	}

	return domain.DistroUnknown
}

func detectPackageManager(os domain.OS, distro domain.Distro) domain.PackageManagerType {
	switch os {
	case domain.OSMacOS:
		if commandExists("brew") {
			return domain.PkgMgrBrew
		}
		return domain.PkgMgrNone
	case domain.OSWindows:
		if commandExists("winget") {
			return domain.PkgMgrWinget
		}
		return domain.PkgMgrNone
	case domain.OSLinux:
		return detectLinuxPackageManager(distro)
	default:
		return domain.PkgMgrNone
	}
}

func detectLinuxPackageManager(distro domain.Distro) domain.PackageManagerType {
	// Try distro-specific first, then fall back to detection
	switch distro {
	case domain.DistroDebian, domain.DistroUbuntu:
		if commandExists("apt") {
			return domain.PkgMgrApt
		}
	case domain.DistroFedora:
		if commandExists("dnf") {
			return domain.PkgMgrDnf
		}
	case domain.DistroArch:
		if commandExists("pacman") {
			return domain.PkgMgrPacman
		}
	}

	// Fall back: try each in order of popularity
	if commandExists("apt") {
		return domain.PkgMgrApt
	}
	if commandExists("dnf") {
		return domain.PkgMgrDnf
	}
	if commandExists("pacman") {
		return domain.PkgMgrPacman
	}

	return domain.PkgMgrNone
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// FormatDetectionInfo returns a formatted string with platform details.
func FormatDetectionInfo(p *domain.Platform) string {
	info := fmt.Sprintf("Detected: %s", p.Summary())
	if p.PackageManager != domain.PkgMgrNone {
		info += fmt.Sprintf(" — Package manager: %s", p.PackageManager)
	} else {
		info += " — No supported package manager found"
	}
	return info
}
