package screens

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/uttejg/newbox/internal/adapter/input/tui/styles"
	"github.com/uttejg/newbox/internal/core/domain"
)

type CategoriesDone struct{ Selected []domain.Category }
type CategoriesBack struct{}

type CategoriesModel struct {
	categories []domain.Category
	checked    []bool
	cursor     int
}

func NewCategories(categories []domain.Category, profile *domain.Profile) CategoriesModel {
	checked := make([]bool, len(categories))
	if profile != nil {
		if profile.AllCategories {
			for i := range checked {
				checked[i] = true
			}
		} else {
			profileCats := make(map[string]bool)
			for _, id := range profile.Categories {
				profileCats[id] = true
			}
			for i, cat := range categories {
				checked[i] = profileCats[cat.ID]
			}
		}
	}
	return CategoriesModel{categories: categories, checked: checked}
}

func (m CategoriesModel) Init() tea.Cmd { return nil }

func (m CategoriesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.categories)-1 {
				m.cursor++
			}
		case " ":
			if len(m.checked) > m.cursor {
				m.checked[m.cursor] = !m.checked[m.cursor]
			}
		case "tab", "enter":
			return m, func() tea.Msg { return CategoriesDone{Selected: m.selectedCategories()} }
		case "esc":
			return m, func() tea.Msg { return CategoriesBack{} }
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m CategoriesModel) View() string {
	title := styles.TitleStyle.Render("Select Categories")
	sub := styles.SubtitleStyle.Render(fmt.Sprintf("Pre-selected by profile  •  %d selected", m.checkedCount()))

	var items string
	for i, cat := range m.categories {
		cursor := "  "
		nameStyle := styles.ItemNameStyle

		if i == m.cursor {
			cursor = styles.SelectedStyle.Render("▸ ")
			nameStyle = styles.SelectedStyle
		}

		var checkbox string
		if m.checked[i] {
			checkbox = styles.CheckedStyle.Render("[✓] ")
		} else {
			checkbox = styles.UncheckedStyle.Render("[ ] ")
		}

		toolCount := styles.ItemCountStyle.Render(fmt.Sprintf(" (%d tools)", len(cat.Tools)))

		items += cursor + checkbox + nameStyle.Render(cat.Name) + toolCount + "\n"
	}

	help := styles.HelpStyle.Render("↑/↓ navigate  •  Space toggle  •  Enter/Tab proceed  •  Esc back  •  q quit")
	return "\n" + title + "\n" + sub + "\n\n" + items + "\n" + help + "\n"
}

func (m CategoriesModel) checkedCount() int {
	count := 0
	for _, c := range m.checked {
		if c {
			count++
		}
	}
	return count
}

func (m CategoriesModel) selectedCategories() []domain.Category {
	var selected []domain.Category
	for i, cat := range m.categories {
		if m.checked[i] {
			selected = append(selected, cat)
		}
	}
	return selected
}
