package domain

// PackageManager identifies a specific package manager.
type PackageManager string

const (
	PackageManagerFormula PackageManager = "formula" // Homebrew formula
	PackageManagerCask    PackageManager = "cask"    // Homebrew cask
	PackageManagerMAS     PackageManager = "mas"     // Mac App Store
	PackageManagerWinget  PackageManager = "winget"  // Windows Package Manager
	PackageManagerApt     PackageManager = "apt"     // Debian/Ubuntu
	PackageManagerDnf     PackageManager = "dnf"     // Fedora/RHEL
	PackageManagerPacman  PackageManager = "pacman"  // Arch Linux
	PackageManagerFlatpak PackageManager = "flatpak" // Flatpak (Linux fallback)
)

// PackageRef maps package managers to their identifier for a given tool.
// Using a map means new package managers can be added without touching the domain struct.
type PackageRef map[PackageManager]string

// IsEmpty returns true if no package reference is set.
func (p PackageRef) IsEmpty() bool {
	return len(p) == 0
}

// Tool represents a single installable application or CLI tool.
type Tool struct {
	Name            string
	Description     string
	Website         string
	DotfilesDefault bool // pre-selected because it's in the existing dotfiles
	MacOS           PackageRef
	Windows         PackageRef
	Linux           PackageRef
}

// IsAvailableOn returns true if this tool has a package mapping for the given OS.
func (t *Tool) IsAvailableOn(os OS) bool {
	return !t.PackageRefFor(os).IsEmpty()
}

// PackageRefFor returns the PackageRef for the given OS, or nil if not available.
func (t *Tool) PackageRefFor(os OS) PackageRef {
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
