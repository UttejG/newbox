package domain

import (
	"fmt"
	"strings"
)

// OS represents the operating system.
type OS int

const (
	OSUnknown OS = iota
	OSMacOS
	OSLinux
	OSWindows
)

func (o OS) String() string {
	switch o {
	case OSMacOS:
		return "macOS"
	case OSLinux:
		return "Linux"
	case OSWindows:
		return "Windows"
	default:
		return "Unknown"
	}
}

// Arch represents the CPU architecture.
type Arch int

const (
	ArchUnknown Arch = iota
	ArchAMD64
	ArchARM64
)

func (a Arch) String() string {
	switch a {
	case ArchAMD64:
		return "amd64"
	case ArchARM64:
		return "arm64"
	default:
		return "unknown"
	}
}

// Distro represents a Linux distribution.
type Distro int

const (
	DistroNone Distro = iota
	DistroDebian
	DistroUbuntu
	DistroFedora
	DistroArch
	DistroUnknown
)

func (d Distro) String() string {
	switch d {
	case DistroDebian:
		return "Debian"
	case DistroUbuntu:
		return "Ubuntu"
	case DistroFedora:
		return "Fedora"
	case DistroArch:
		return "Arch"
	case DistroNone:
		return ""
	default:
		return "Unknown"
	}
}

// PackageManagerType represents a system package manager.
type PackageManagerType int

const (
	PkgMgrNone PackageManagerType = iota
	PkgMgrBrew
	PkgMgrWinget
	PkgMgrApt
	PkgMgrDnf
	PkgMgrPacman
)

func (p PackageManagerType) String() string {
	switch p {
	case PkgMgrBrew:
		return "brew"
	case PkgMgrWinget:
		return "winget"
	case PkgMgrApt:
		return "apt"
	case PkgMgrDnf:
		return "dnf"
	case PkgMgrPacman:
		return "pacman"
	default:
		return "none"
	}
}

// Platform holds information about the current system.
type Platform struct {
	OS             OS
	Arch           Arch
	Distro         Distro
	PackageManager PackageManagerType
}

// Summary returns a human-readable description of the platform.
func (p *Platform) Summary() string {
	var parts []string
	parts = append(parts, p.OS.String())

	if p.Distro != DistroNone {
		parts = append(parts, fmt.Sprintf("(%s)", p.Distro))
	}

	parts = append(parts, p.Arch.String())

	return strings.Join(parts, " ")
}

// FormatInfo returns a single-line description of the detected platform and
// package manager, suitable for display in a UI or log line.
func (p *Platform) FormatInfo() string {
	info := fmt.Sprintf("Detected: %s", p.Summary())
	if p.PackageManager != PkgMgrNone {
		return info + fmt.Sprintf(" — Package manager: %s", p.PackageManager)
	}
	return info + " — No supported package manager found"
}
