package domain

// PackageRef holds the package identifier(s) for a specific package manager type.
type PackageRef struct {
Formula string // brew formula (e.g. "git")
Cask    string // brew cask (e.g. "signal")
MAS     string // Mac App Store numeric ID (e.g. "409203825")
Winget  string // winget package ID (e.g. "OpenWhisperSystems.Signal")
Apt     string // apt package name
Dnf     string // dnf package name
Pacman  string // pacman package name
Flatpak string // flatpak app ID (Linux fallback)
}

// IsEmpty returns true if no package reference is set.
func (p PackageRef) IsEmpty() bool {
return p.Formula == "" && p.Cask == "" && p.MAS == "" &&
p.Winget == "" && p.Apt == "" && p.Dnf == "" &&
p.Pacman == "" && p.Flatpak == ""
}

// Tool represents a single installable application or CLI tool.
type Tool struct {
Name            string
Description     string
Website         string
DotfilesDefault bool // pre-selected because it's in the existing dotfiles
MacOS           *PackageRef
Windows         *PackageRef
Linux           *PackageRef
}

// IsAvailableOn returns true if this tool has a package mapping for the given OS.
func (t *Tool) IsAvailableOn(os OS) bool {
switch os {
case OSMacOS:
return t.MacOS != nil && !t.MacOS.IsEmpty()
case OSWindows:
return t.Windows != nil && !t.Windows.IsEmpty()
case OSLinux:
return t.Linux != nil && !t.Linux.IsEmpty()
default:
return false
}
}

// PackageRefFor returns the PackageRef for the given OS, or nil if not available.
func (t *Tool) PackageRefFor(os OS) *PackageRef {
switch os {
case OSMacOS:
return t.MacOS
case OSWindows:
return t.Windows
case OSLinux:
return t.Linux
default:
return nil
}
}
