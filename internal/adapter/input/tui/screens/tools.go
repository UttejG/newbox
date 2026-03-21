package screens

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/uttejg/newbox/internal/adapter/input/tui/keys"
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
		switch {
		case key.Matches(msg, keys.Checklist.Up):
			if m.cursors[m.catIndex] > 0 {
				m.cursors[m.catIndex]--
			}
		case key.Matches(msg, keys.Checklist.Down):
			if m.cursors[m.catIndex] < len(currentCat.Tools)-1 {
				m.cursors[m.catIndex]++
			}
		case key.Matches(msg, keys.Checklist.Toggle):
			idx := m.cursors[m.catIndex]
			if len(m.checked[m.catIndex]) > 0 {
				m.checked[m.catIndex][idx] = !m.checked[m.catIndex][idx]
			}
		case key.Matches(msg, keys.Checklist.Next):
			if m.catIndex < len(m.categories)-1 {
				m.catIndex++
			} else {
				return m, func() tea.Msg { return ToolsDone{ByCategory: m.selectedByCategory()} }
			}
		case key.Matches(msg, keys.Checklist.Prev):
			if m.catIndex > 0 {
				m.catIndex--
			} else {
				return m, func() tea.Msg { return ToolsBack{} }
			}
		case key.Matches(msg, keys.Checklist.Back):
			return m, func() tea.Msg { return ToolsBack{} }
		case key.Matches(msg, keys.Checklist.Quit):
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
	progress := styles.ItemCountStyle.Render(
		fmt.Sprintf("Category %d of %d", m.catIndex+1, len(m.categories)),
	)
	title := styles.TitleStyle.Render(cat.Name)
	selectedCount := m.checkedCountForCat(m.catIndex)
	sub := styles.SubtitleStyle.Render(fmt.Sprintf("%d selected  •  ★ = from your dotfiles", selectedCount))

	var items string
	for i, tool := range cat.Tools {
		cursor := "  "
		nameStyle := styles.ItemNameStyle

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
		items += "      " + styles.ItemDescStyle.Render(tool.Description) + "\n"
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
