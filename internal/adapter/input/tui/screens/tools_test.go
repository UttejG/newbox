package screens_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/uttejg/newbox/internal/adapter/input/tui/screens"
	"github.com/uttejg/newbox/internal/core/domain"
)

func testToolCategories() []domain.Category {
	return []domain.Category{
		{
			ID:   "cli",
			Name: "CLI Tools",
			Tools: []domain.Tool{
				{Name: "git", Description: "version control", DotfilesDefault: false},
				{Name: "jq", Description: "JSON processor", DotfilesDefault: false},
			},
		},
		{
			ID:    "browsers",
			Name:  "Browsers",
			Tools: []domain.Tool{{Name: "Firefox", Description: "web browser"}},
		},
	}
}

func TestToolsModel_SpaceTogglesTool(t *testing.T) {
	cats := testToolCategories()
	m := screens.NewTools(cats)
	// Space to toggle first tool (git, initially unchecked)
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	if cmd != nil {
		t.Error("expected no command after Space")
	}
	// Enter to proceed to next category
	updated, _ = updated.Update(tea.KeyMsg{Type: tea.KeyEnter})
	// Enter again on last category to emit ToolsDone
	_, cmd = updated.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected ToolsDone on last Enter")
	}
	msg := cmd()
	done, ok := msg.(screens.ToolsDone)
	if !ok {
		t.Fatalf("expected ToolsDone, got %T", msg)
	}
	tools := done.ByCategory["cli"]
	if len(tools) != 1 || tools[0].Name != "git" {
		t.Errorf("expected git toggled into result, got %v", tools)
	}
}

func TestToolsModel_EnterOnSingleCategoryEmitsDone(t *testing.T) {
	cats := []domain.Category{
		{ID: "cli", Name: "CLI Tools", Tools: []domain.Tool{{Name: "git"}}},
	}
	m := screens.NewTools(cats)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected ToolsDone command on Enter with single category")
	}
	msg := cmd()
	if _, ok := msg.(screens.ToolsDone); !ok {
		t.Fatalf("expected ToolsDone, got %T", msg)
	}
}

func TestToolsModel_EmptyCategoriesEmitsDoneImmediately(t *testing.T) {
	m := screens.NewTools([]domain.Category{})
	// Update with any message — the empty guard emits ToolsDone
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected ToolsDone for empty categories")
	}
	msg := cmd()
	if _, ok := msg.(screens.ToolsDone); !ok {
		t.Fatalf("expected ToolsDone, got %T", msg)
	}
}

func TestToolsModel_ViewContainsToolNames(t *testing.T) {
	cats := testToolCategories()
	m := screens.NewTools(cats)
	view := m.View()
	for _, name := range []string{"git", "jq"} {
		if !strings.Contains(view, name) {
			t.Errorf("expected tool name %q in view", name)
		}
	}
}

func TestToolsModel_DotfilesDefaultPreChecked(t *testing.T) {
	cats := []domain.Category{
		{
			ID:   "cli",
			Name: "CLI Tools",
			Tools: []domain.Tool{
				{Name: "vim", DotfilesDefault: true},
			},
		},
	}
	m := screens.NewTools(cats)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected ToolsDone")
	}
	msg := cmd()
	done := msg.(screens.ToolsDone)
	tools := done.ByCategory["cli"]
	if len(tools) != 1 || tools[0].Name != "vim" {
		t.Errorf("expected vim pre-checked via DotfilesDefault, got %v", tools)
	}
}

func TestToolsModel_EscEmitsToolsBack(t *testing.T) {
	m := screens.NewTools(testToolCategories())
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("expected command on Esc")
	}
	msg := cmd()
	if _, ok := msg.(screens.ToolsBack); !ok {
		t.Errorf("expected ToolsBack, got %T", msg)
	}
}
