package screens_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/uttejg/newbox/internal/adapter/input/tui/screens"
	"github.com/uttejg/newbox/internal/core/domain"
)

func testProfiles() []domain.Profile {
	return []domain.Profile{
		{ID: "developer", Name: "Developer", Description: "Dev tools"},
		{ID: "creative", Name: "Creative", Description: "Creative tools"},
		{ID: "minimal", Name: "Minimal", Description: "Minimal setup"},
	}
}

func TestProfileModel_DownMovesCursor(t *testing.T) {
	m := screens.NewProfile(testProfiles())
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	pm := updated.(screens.ProfileModel)
	if pm.Selected().ID != "creative" {
		t.Errorf("expected creative after Down, got %s", pm.Selected().ID)
	}
}

func TestProfileModel_UpDoesNotGoNegative(t *testing.T) {
	m := screens.NewProfile(testProfiles())
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	pm := updated.(screens.ProfileModel)
	if pm.Selected().ID != "developer" {
		t.Errorf("expected cursor to stay at developer, got %s", pm.Selected().ID)
	}
}

func TestProfileModel_EnterEmitsProfileSelected(t *testing.T) {
	m := screens.NewProfile(testProfiles())
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected command on Enter")
	}
	msg := cmd()
	ps, ok := msg.(screens.ProfileSelected)
	if !ok {
		t.Fatalf("expected ProfileSelected, got %T", msg)
	}
	if ps.Profile.ID != "developer" {
		t.Errorf("expected developer profile, got %s", ps.Profile.ID)
	}
}

func TestProfileModel_DownThenEnterSelectsSecond(t *testing.T) {
	m := screens.NewProfile(testProfiles())
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	_, cmd := updated.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected command on Enter")
	}
	msg := cmd()
	ps, ok := msg.(screens.ProfileSelected)
	if !ok {
		t.Fatalf("expected ProfileSelected, got %T", msg)
	}
	if ps.Profile.ID != "creative" {
		t.Errorf("expected creative profile, got %s", ps.Profile.ID)
	}
}

func TestProfileModel_ViewContainsProfileNames(t *testing.T) {
	m := screens.NewProfile(testProfiles())
	view := m.View()
	for _, name := range []string{"Developer", "Creative", "Minimal"} {
		if !strings.Contains(view, name) {
			t.Errorf("expected profile name %q in view", name)
		}
	}
}

func TestProfileModel_EscEmitsProfileBack(t *testing.T) {
	m := screens.NewProfile(testProfiles())
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("expected command on Esc")
	}
	msg := cmd()
	if _, ok := msg.(screens.ProfileBack); !ok {
		t.Errorf("expected ProfileBack, got %T", msg)
	}
}
