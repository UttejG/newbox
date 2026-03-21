package screens

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/uttejg/newbox/internal/adapter/input/tui/styles"
	"github.com/uttejg/newbox/internal/adapter/output/detector"
	"github.com/uttejg/newbox/internal/core/domain"
)

type WelcomeDone struct{}

type WelcomeModel struct {
	platform *domain.Platform
	dryRun   bool
}

func NewWelcome(platform *domain.Platform, dryRun bool) WelcomeModel {
	return WelcomeModel{platform: platform, dryRun: dryRun}
}

func (m WelcomeModel) Init() tea.Cmd { return nil }

func (m WelcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, func() tea.Msg { return WelcomeDone{} }
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m WelcomeModel) View() string {
	banner := styles.TitleStyle.Render(`
 ███╗   ██╗███████╗██╗    ██╗██████╗  ██████╗ ██╗  ██╗
 ████╗  ██║██╔════╝██║    ██║██╔══██╗██╔═══██╗╚██╗██╔╝
 ██╔██╗ ██║█████╗  ██║ █╗ ██║██████╔╝██║   ██║ ╚███╔╝
 ██║╚██╗██║██╔══╝  ██║███╗██║██╔══██╗██║   ██║ ██╔██╗
 ██║ ╚████║███████╗╚███╔███╔╝██████╔╝╚██████╔╝██╔╝ ██╗
 ╚═╝  ╚═══╝╚══════╝ ╚══╝╚══╝ ╚═════╝  ╚═════╝ ╚═╝  ╚═╝`)

	sub := styles.SubtitleStyle.Render("Cross-platform machine setup — choose what to install")

	platformInfo := lipgloss.NewStyle().Foreground(styles.Primary).Render(
		"  " + detector.FormatDetectionInfo(m.platform),
	)

	var dryRunBadge string
	if m.dryRun {
		dryRunBadge = "\n  " + lipgloss.NewStyle().
			Foreground(styles.Warning).
			Bold(true).
			Render("[DRY RUN] — no packages will be installed")
	}

	hint := styles.HelpStyle.Render("  Press Enter to continue  •  q to quit")

	return "\n" + banner + "\n" + sub + "\n\n" + platformInfo + dryRunBadge + "\n\n" + hint + "\n"
}
