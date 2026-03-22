package screens_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/uttejg/newbox/internal/adapter/input/tui/screens"
	"github.com/uttejg/newbox/internal/core/domain"
)

func testSelection() *domain.UserSelection {
	return &domain.UserSelection{
		ToolsByCategory: map[string][]domain.Tool{
			"cli": {{Name: "git"}, {Name: "jq"}},
		},
	}
}

func TestConfirmModel_ViewContainsInstall(t *testing.T) {
	cats := []domain.Category{{ID: "cli", Name: "CLI Tools"}}
	m := screens.NewConfirm(testSelection(), cats)
	if !strings.Contains(m.View(), "Install") {
		t.Error("expected 'Install' in confirm view")
	}
}

func TestConfirmModel_ViewContainsToolNames(t *testing.T) {
	cats := []domain.Category{{ID: "cli", Name: "CLI Tools"}}
	m := screens.NewConfirm(testSelection(), cats)
	view := m.View()
	for _, name := range []string{"git", "jq"} {
		if !strings.Contains(view, name) {
			t.Errorf("expected tool name %q in confirm view", name)
		}
	}
}

func TestConfirmModel_EnterEmitsConfirmProceed(t *testing.T) {
	cats := []domain.Category{{ID: "cli", Name: "CLI Tools"}}
	m := screens.NewConfirm(testSelection(), cats)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected command on Enter with non-empty selection")
	}
	msg := cmd()
	cp, ok := msg.(screens.ConfirmProceed)
	if !ok {
		t.Fatalf("expected ConfirmProceed, got %T", msg)
	}
	if cp.Selection.TotalCount() != 2 {
		t.Errorf("expected 2 tools in selection, got %d", cp.Selection.TotalCount())
	}
}

func TestConfirmModel_EmptySelectionEnterNoCmd(t *testing.T) {
	empty := &domain.UserSelection{ToolsByCategory: map[string][]domain.Tool{}}
	m := screens.NewConfirm(empty, nil)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Error("expected no command on Enter with empty selection")
	}
}

func TestConfirmModel_EscEmitsConfirmBack(t *testing.T) {
	m := screens.NewConfirm(testSelection(), nil)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("expected command on Esc")
	}
	msg := cmd()
	if _, ok := msg.(screens.ConfirmBack); !ok {
		t.Errorf("expected ConfirmBack, got %T", msg)
	}
}
