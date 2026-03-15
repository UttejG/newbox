package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Primary color: teal/cyan
	Primary   = lipgloss.Color("#00D7D7")
	Secondary = lipgloss.Color("#7D56F4")
	Subtle    = lipgloss.Color("#555555")
	Success   = lipgloss.Color("#00CC66")
	Warning   = lipgloss.Color("#FFAA00")
	Danger    = lipgloss.Color("#FF4444")
	Text      = lipgloss.Color("#EEEEEE")
	Muted     = lipgloss.Color("#888888")

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

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(1, 2)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(Muted).
			MarginTop(1)
)
