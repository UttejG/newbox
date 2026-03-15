package screens

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/uttejg/newbox/internal/adapter/input/tui/styles"
	"github.com/uttejg/newbox/internal/core/domain"
)

type ToolsDone struct{ ByCategory map[string][]domain.Tool }
type ToolsBack struct{}

type ToolsModel struct {
	categories []domain.Category
	catIndex   int
	cursors    []int
	checked    [][]bool // checked[catIdx][toolIdx]
}

func NewTools(categories []domain.Category) ToolsModel {
	cursors := make([]int, len(categories))
	checked := make([][]bool, len(categories))
	for i, cat := range categories {
		checked[i] = make([]bool, len(cat.Tools))
		// Pre-check dotfiles defaults
		for j, tool := range cat.Tools {
			checked[i][j] = tool.DotfilesDefault
		}
	}
	return ToolsModel{categories: categories, cursors: cursors, checked: checked}
}

func (m ToolsModel) Init() tea.Cmd { return nil }

func (m ToolsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if len(m.categories) == 0 {
		return m, func() tea.Msg { return ToolsDone{ByCategory: map[string][]domain.Tool{}} }
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		currentCat := m.categories[m.catIndex]
		switch msg.String() {
		case "up", "k":
			if m.cursors[m.catIndex] > 0 {
				m.cursors[m.catIndex]--
			}
		case "down", "j":
			if m.cursors[m.catIndex] < len(currentCat.Tools)-1 {
				m.cursors[m.catIndex]++
			}
		case " ":
			idx := m.cursors[m.catIndex]
			m.checked[m.catIndex][idx] = !m.checked[m.catIndex][idx]
		case "tab", "enter":
			if m.catIndex < len(m.categories)-1 {
				m.catIndex++
			} else {
				return m, func() tea.Msg { return ToolsDone{ByCategory: m.selectedByCategory()} }
			}
		case "shift+tab":
			if m.catIndex > 0 {
				m.catIndex--
			} else {
				return m, func() tea.Msg { return ToolsBack{} }
			}
		case "esc":
			return m, func() tea.Msg { return ToolsBack{} }
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ToolsModel) View() string {
	if len(m.categories) == 0 {
		return styles.SubtitleStyle.Render("No categories selected.")
	}

	cat := m.categories[m.catIndex]
	progress := lipgloss.NewStyle().Foreground(styles.Muted).Render(
		fmt.Sprintf("Category %d of %d", m.catIndex+1, len(m.categories)),
	)
	title := styles.TitleStyle.Render(cat.Name)
	selectedCount := m.checkedCountForCat(m.catIndex)
	sub := styles.SubtitleStyle.Render(fmt.Sprintf("%d selected  •  ★ = from your dotfiles", selectedCount))

	var items string
	for i, tool := range cat.Tools {
		cursor := "  "
		nameStyle := lipgloss.NewStyle().Foreground(styles.Text)
		descStyle := lipgloss.NewStyle().Foreground(styles.Muted)

		if i == m.cursors[m.catIndex] {
			cursor = styles.SelectedStyle.Render("▸ ")
			nameStyle = styles.SelectedStyle
		}

		var checkbox string
		if m.checked[m.catIndex][i] {
			checkbox = styles.CheckedStyle.Render("[✓] ")
		} else {
			checkbox = styles.UncheckedStyle.Render("[ ] ")
		}

		dotfilesMark := ""
		if tool.DotfilesDefault {
			dotfilesMark = styles.DotfilesMarker.String()
		}

		items += cursor + checkbox + nameStyle.Render(tool.Name) + dotfilesMark + "\n"
		items += "      " + descStyle.Render(tool.Description) + "\n"
	}

	help := styles.HelpStyle.Render("↑/↓ navigate  •  Space toggle  •  Tab/Enter next category  •  Shift+Tab prev  •  Esc back")
	return "\n" + progress + "\n" + title + "\n" + sub + "\n\n" + items + "\n" + help + "\n"
}

func (m ToolsModel) checkedCountForCat(catIdx int) int {
	count := 0
	for _, c := range m.checked[catIdx] {
		if c {
			count++
		}
	}
	return count
}

func (m ToolsModel) selectedByCategory() map[string][]domain.Tool {
	result := make(map[string][]domain.Tool)
	for i, cat := range m.categories {
		var tools []domain.Tool
		for j, tool := range cat.Tools {
			if m.checked[i][j] {
				tools = append(tools, tool)
			}
		}
		if len(tools) > 0 {
			result[cat.ID] = tools
		}
	}
	return result
}
