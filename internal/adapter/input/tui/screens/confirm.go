package screens

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/uttejg/newbox/internal/adapter/input/tui/styles"
	"github.com/uttejg/newbox/internal/core/domain"
)

type ConfirmProceed struct{ Selection *domain.UserSelection }
type ConfirmBack struct{}

type ConfirmModel struct {
	selection  *domain.UserSelection
	categories []domain.Category // ordered category names for display
}

func NewConfirm(selection *domain.UserSelection, categories []domain.Category) ConfirmModel {
	return ConfirmModel{selection: selection, categories: categories}
}

func (m ConfirmModel) Init() tea.Cmd { return nil }

func (m ConfirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.selection.TotalCount() > 0 {
				return m, func() tea.Msg { return ConfirmProceed{Selection: m.selection} }
			}
		case "esc":
			return m, func() tea.Msg { return ConfirmBack{} }
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ConfirmModel) View() string {
	title := styles.TitleStyle.Render("Review Selections")

	total := m.selection.TotalCount()
	sub := styles.SubtitleStyle.Render(fmt.Sprintf("%d tools selected", total))

	var body string
	for _, cat := range m.categories {
		tools, ok := m.selection.ToolsByCategory[cat.ID]
		if !ok || len(tools) == 0 {
			continue
		}

		catName := styles.CatHeaderStyle.Render(
			fmt.Sprintf("%s (%d)", cat.Name, len(tools)),
		)
		body += "  " + catName + "\n"
		for _, t := range tools {
			marker := ""
			if t.DotfilesDefault {
				marker = styles.DotfilesMarker.String()
			}
			body += "    • " + t.Name + marker + "\n"
		}
		body += "\n"
	}

	if total == 0 {
		body = styles.WarningTextStyle.Render("  Nothing selected. Go back and choose some tools.") + "\n\n"
	}

	var proceedHint string
	if total > 0 {
		proceedHint = styles.ProceedStyle.Render("  Enter: Install")
	} else {
		proceedHint = styles.MutedTextStyle.Render("  (select tools first)")
	}
	help := styles.HelpStyle.Render("  Esc: go back  •  q: quit")

	return "\n" + title + "\n" + sub + "\n\n" + body + proceedHint + "\n" + help + "\n"
}
