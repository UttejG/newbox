package screens

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/key"
	"github.com/uttejg/newbox/internal/adapter/input/tui/keys"
	"github.com/uttejg/newbox/internal/adapter/input/tui/styles"
	"github.com/uttejg/newbox/internal/core/domain"
)

type ProfileSelected struct{ Profile domain.Profile }
type ProfileBack struct{}

type ProfileModel struct {
	profiles []domain.Profile
	cursor   int
}

func NewProfile(profiles []domain.Profile) ProfileModel {
	return ProfileModel{profiles: profiles}
}

func (m ProfileModel) Init() tea.Cmd { return nil }

func (m ProfileModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.List.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, keys.List.Down):
			if m.cursor < len(m.profiles)-1 {
				m.cursor++
			}
		case key.Matches(msg, keys.List.Select):
			if len(m.profiles) > 0 {
				selected := m.profiles[m.cursor]
				return m, func() tea.Msg { return ProfileSelected{Profile: selected} }
			}
		case key.Matches(msg, keys.List.Back):
			return m, func() tea.Msg { return ProfileBack{} }
		case key.Matches(msg, keys.List.Quit):
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ProfileModel) View() string {
	title := styles.TitleStyle.Render("Select a Profile")
	sub := styles.SubtitleStyle.Render("Profiles pre-select categories — you can customise afterwards")

	var items string
	for i, p := range m.profiles {
		cursor := "  "
		nameStyle := styles.ItemNameStyle

		if i == m.cursor {
			cursor = styles.SelectedStyle.Render("▸ ")
			nameStyle = styles.SelectedStyle
		}

		items += cursor + nameStyle.Render(p.Name) + "\n"
		items += "    " + styles.ItemDescStyle.Render(p.Description) + "\n\n"
	}

	help := styles.HelpStyle.Render("↑/↓ navigate  •  Enter select  •  Esc back  •  q quit")
	return "\n" + title + "\n" + sub + "\n\n" + items + help + "\n"
}

// Selected returns the currently highlighted profile.
func (m ProfileModel) Selected() domain.Profile {
	if len(m.profiles) == 0 {
		return domain.Profile{}
	}
	return m.profiles[m.cursor]
}
