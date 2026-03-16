package screens

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/uttejg/newbox/internal/adapter/input/tui/styles"
	"github.com/uttejg/newbox/internal/core/domain"
)

// InstallMsg wraps a ProgressEvent for BubbleTea messaging.
type InstallMsg domain.ProgressEvent

// InstallDoneMsg signals the progress channel was closed.
type InstallDoneMsg struct{}

// InstallModel is the TUI screen shown during (or previewing) installation.
type InstallModel struct {
	steps    []domain.ExecutionStep
	dryRun   bool
	done     bool
	total    int // total pending/dry-run steps
	progress <-chan domain.ProgressEvent
	width    int
}

// NewInstall creates an InstallModel from the computed plan and progress channel.
func NewInstall(plan *domain.InstallPlan, dryRun bool, ch <-chan domain.ProgressEvent) InstallModel {
	steps := make([]domain.ExecutionStep, len(plan.Steps))
	copy(steps, plan.Steps)
	return InstallModel{
		steps:    steps,
		dryRun:   dryRun,
		total:    len(plan.PendingSteps()),
		progress: ch,
	}
}

func (m InstallModel) Init() tea.Cmd {
	return waitForProgress(m.progress)
}

func (m InstallModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case InstallMsg:
		ev := domain.ProgressEvent(msg)
		for i, s := range m.steps {
			if s.Tool.Name == ev.Step.Tool.Name {
				m.steps[i] = ev.Step
				break
			}
		}
		return m, waitForProgress(m.progress)

	case InstallDoneMsg:
		m.done = true
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "enter":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m InstallModel) View() string {
	var title string
	if m.dryRun {
		badge := lipgloss.NewStyle().Foreground(styles.Warning).Bold(true).Render("[DRY RUN]")
		title = styles.TitleStyle.Render("Install Plan") + "  " + badge
	} else {
		title = styles.TitleStyle.Render("Installing")
	}

	var rows strings.Builder
	doneCount := 0
	for _, s := range m.steps {
		rows.WriteString(m.renderStep(s))
		rows.WriteByte('\n')
		if s.Status == domain.StatusDone || s.Status == domain.StatusFailed {
			doneCount++
		}
	}

	var progressBar string
	if m.total > 0 {
		progressBar = "\n  " + renderProgressBar(doneCount, m.total, 30) +
			fmt.Sprintf("  %d/%d\n", doneCount, m.total)
	}

	var footer string
	if m.done {
		footer = "\n  " + lipgloss.NewStyle().Foreground(styles.Success).Bold(true).
			Render("✓ Done! Press Enter or q to quit") + "\n"
	}

	help := styles.HelpStyle.Render("  q: quit")
	return "\n" + title + "\n\n" + rows.String() + progressBar + footer + help + "\n"
}

func (m InstallModel) renderStep(s domain.ExecutionStep) string {
	icon := statusIcon(s.Status)
	name := lipgloss.NewStyle().Foreground(styles.Text).Render(s.Tool.Name)
	cmd := lipgloss.NewStyle().Foreground(styles.Muted).Render("  " + s.Command)

	status := ""
	switch s.Status {
	case domain.StatusSkipped:
		status = lipgloss.NewStyle().Foreground(styles.Muted).Render("already installed")
	case domain.StatusFailed:
		msg := ""
		if s.Error != nil {
			msg = s.Error.Error()
		}
		status = lipgloss.NewStyle().Foreground(styles.Danger).Render("FAILED: " + msg)
	}

	line := "  " + icon + " " + name
	if s.Command != "" {
		line += "\n      " + cmd
	}
	if status != "" {
		line += "  " + status
	}
	return line
}

func statusIcon(s domain.InstallStatus) string {
	switch s {
	case domain.StatusPending:
		return lipgloss.NewStyle().Foreground(styles.Muted).Render("⬚")
	case domain.StatusDryRun:
		return lipgloss.NewStyle().Foreground(styles.Warning).Render("⬚")
	case domain.StatusInstalling:
		return lipgloss.NewStyle().Foreground(styles.Primary).Render("⏳")
	case domain.StatusDone:
		return lipgloss.NewStyle().Foreground(styles.Success).Render("✅")
	case domain.StatusSkipped:
		return lipgloss.NewStyle().Foreground(styles.Muted).Render("⏭️")
	case domain.StatusFailed:
		return lipgloss.NewStyle().Foreground(styles.Danger).Render("✗")
	default:
		return " "
	}
}

func renderProgressBar(done, total, width int) string {
	if total == 0 {
		return ""
	}
	filled := (done * width) / total
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return lipgloss.NewStyle().Foreground(styles.Primary).Render(bar)
}

func waitForProgress(ch <-chan domain.ProgressEvent) tea.Cmd {
	return func() tea.Msg {
		ev, ok := <-ch
		if !ok {
			return InstallDoneMsg{}
		}
		return InstallMsg(ev)
	}
}
