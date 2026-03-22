package screens_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/uttejg/newbox/internal/adapter/input/tui/screens"
	"github.com/uttejg/newbox/internal/core/domain"
)

func testCategories() []domain.Category {
	return []domain.Category{
		{ID: "cli", Name: "CLI Tools", Tools: []domain.Tool{{Name: "git"}, {Name: "jq"}}},
		{ID: "browsers", Name: "Browsers", Tools: []domain.Tool{{Name: "Firefox"}}},
	}
}

func TestCategoriesModel_SpaceToggles(t *testing.T) {
	m := screens.NewCategories(testCategories(), nil) // nothing pre-selected
	// Press Space to toggle the first category
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	if cmd != nil {
		t.Error("expected no command after Space")
	}
	// Press Enter to emit CategoriesDone
	_, cmd = updated.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected command on Enter")
	}
	msg := cmd()
	done, ok := msg.(screens.CategoriesDone)
	if !ok {
		t.Fatalf("expected CategoriesDone, got %T", msg)
	}
	if len(done.Selected) != 1 || done.Selected[0].ID != "cli" {
		t.Errorf("expected only 'cli' selected after Space, got %v", done.Selected)
	}
}

func TestCategoriesModel_EnterWithNothingSelectedEmitsDone(t *testing.T) {
	m := screens.NewCategories(testCategories(), nil)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected command on Enter")
	}
	msg := cmd()
	done, ok := msg.(screens.CategoriesDone)
	if !ok {
		t.Fatalf("expected CategoriesDone, got %T", msg)
	}
	if len(done.Selected) != 0 {
		t.Errorf("expected empty selection, got %v", done.Selected)
	}
}

func TestCategoriesModel_ProfilePreSelectsCategories(t *testing.T) {
	profile := &domain.Profile{
		ID:         "dev",
		Categories: []string{"cli"},
	}
	m := screens.NewCategories(testCategories(), profile)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected command on Enter")
	}
	msg := cmd()
	done, ok := msg.(screens.CategoriesDone)
	if !ok {
		t.Fatalf("expected CategoriesDone, got %T", msg)
	}
	if len(done.Selected) != 1 || done.Selected[0].ID != "cli" {
		t.Errorf("expected cli pre-selected by profile, got %v", done.Selected)
	}
}

func TestCategoriesModel_ViewContainsCategoryNames(t *testing.T) {
	m := screens.NewCategories(testCategories(), nil)
	view := m.View()
	for _, name := range []string{"CLI Tools", "Browsers"} {
		if !strings.Contains(view, name) {
			t.Errorf("expected category name %q in view", name)
		}
	}
}

func TestCategoriesModel_EscEmitsCategoriesBack(t *testing.T) {
	m := screens.NewCategories(testCategories(), nil)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("expected command on Esc")
	}
	msg := cmd()
	if _, ok := msg.(screens.CategoriesBack); !ok {
		t.Errorf("expected CategoriesBack, got %T", msg)
	}
}
