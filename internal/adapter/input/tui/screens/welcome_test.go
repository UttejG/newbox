package screens_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/uttejg/newbox/internal/adapter/input/tui/screens"
	"github.com/uttejg/newbox/internal/core/domain"
)

func TestWelcomeModel_EnterEmitsWelcomeDone(t *testing.T) {
	m := screens.NewWelcome(&domain.Platform{OS: domain.OSMacOS}, false)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected a command on Enter, got nil")
	}
	msg := cmd()
	if _, ok := msg.(screens.WelcomeDone); !ok {
		t.Errorf("expected WelcomeDone, got %T", msg)
	}
}

func TestWelcomeModel_QuitEmitsQuit(t *testing.T) {
	m := screens.NewWelcome(&domain.Platform{OS: domain.OSMacOS}, false)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Fatal("expected a command on q, got nil")
	}
}

func TestWelcomeModel_DryRunBadgePresent(t *testing.T) {
	m := screens.NewWelcome(&domain.Platform{OS: domain.OSMacOS}, true)
	if !strings.Contains(m.View(), "DRY RUN") {
		t.Error("expected DRY RUN badge in dry-run mode")
	}
}

func TestWelcomeModel_NoDryRunBadge(t *testing.T) {
	m := screens.NewWelcome(&domain.Platform{OS: domain.OSMacOS}, false)
	if strings.Contains(m.View(), "DRY RUN") {
		t.Error("expected no DRY RUN badge in normal mode")
	}
}

func TestWelcomeModel_UnrelatedKeyNoCmd(t *testing.T) {
	m := screens.NewWelcome(&domain.Platform{OS: domain.OSMacOS}, false)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	if cmd != nil {
		t.Error("expected no command for unrelated key")
	}
}
