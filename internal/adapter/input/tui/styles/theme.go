package styles

import "github.com/charmbracelet/lipgloss"

var (
	// All colors use AdaptiveColor so the TUI is readable on both light and dark terminals.
	Primary   lipgloss.TerminalColor = lipgloss.AdaptiveColor{Light: "#007070", Dark: "#00D7D7"}
	Secondary lipgloss.TerminalColor = lipgloss.AdaptiveColor{Light: "#5A3CB3", Dark: "#7D56F4"}
	Subtle    lipgloss.TerminalColor = lipgloss.AdaptiveColor{Light: "#999999", Dark: "#555555"}
	Success   lipgloss.TerminalColor = lipgloss.AdaptiveColor{Light: "#007A3D", Dark: "#00CC66"}
	Warning   lipgloss.TerminalColor = lipgloss.AdaptiveColor{Light: "#995500", Dark: "#FFAA00"}
	Danger    lipgloss.TerminalColor = lipgloss.AdaptiveColor{Light: "#CC0000", Dark: "#FF4444"}
	Text      lipgloss.TerminalColor = lipgloss.AdaptiveColor{Light: "#111111", Dark: "#EEEEEE"}
	Muted     lipgloss.TerminalColor = lipgloss.AdaptiveColor{Light: "#666666", Dark: "#888888"}

	TitleStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(Muted).
			Italic(true)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true)

	CheckedStyle = lipgloss.NewStyle().
			Foreground(Success)

	UncheckedStyle = lipgloss.NewStyle().
			Foreground(Subtle)

	DotfilesMarker = lipgloss.NewStyle().
			Foreground(Warning).
			SetString(" ★")

	HelpStyle = lipgloss.NewStyle().
			Foreground(Muted).
			MarginTop(1)

	// ItemNameStyle and ItemDescStyle are used in list/checklist item rendering.
	// Defined here to avoid allocating new Style values on every View() call.
	ItemNameStyle = lipgloss.NewStyle().Foreground(Text)
	ItemDescStyle = lipgloss.NewStyle().Foreground(Muted)
	ItemCountStyle = lipgloss.NewStyle().Foreground(Muted)

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(1, 2)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(Muted).
			MarginTop(1)
)
