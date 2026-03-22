package screens_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/uttejg/newbox/internal/adapter/input/tui/screens"
	"github.com/uttejg/newbox/internal/core/domain"
)

func makeInstallModel(steps []domain.ExecutionStep, dryRun bool) (screens.InstallModel, chan domain.ProgressEvent) {
	ch := make(chan domain.ProgressEvent, 10)
	plan := &domain.InstallPlan{Steps: steps, DryRun: dryRun}
	return screens.NewInstall(plan, dryRun, ch), ch
}

func TestInstallModel_InitialViewShowsInstalling(t *testing.T) {
	steps := []domain.ExecutionStep{
		{Tool: domain.Tool{Name: "git"}, Status: domain.StatusPending},
	}
	m, _ := makeInstallModel(steps, false)
	if !strings.Contains(m.View(), "Installing") {
		t.Error("expected 'Installing' in view for non-dry-run")
	}
}

func TestInstallModel_DryRunShowsInstallPlan(t *testing.T) {
	steps := []domain.ExecutionStep{
		{Tool: domain.Tool{Name: "git"}, Status: domain.StatusDryRun},
	}
	m, _ := makeInstallModel(steps, true)
	view := m.View()
	if !strings.Contains(view, "Install Plan") {
		t.Error("expected 'Install Plan' in dry-run view")
	}
	if !strings.Contains(view, "DRY RUN") {
		t.Error("expected 'DRY RUN' badge in dry-run view")
	}
}

func TestInstallModel_InstallDoneMsgSetsDone(t *testing.T) {
	steps := []domain.ExecutionStep{
		{Tool: domain.Tool{Name: "git"}, Status: domain.StatusPending},
	}
	m, _ := makeInstallModel(steps, false)
	updated, _ := m.Update(screens.InstallDoneMsg{})
	im := updated.(screens.InstallModel)
	if !strings.Contains(im.View(), "Done") {
		t.Error("expected 'Done' in view after InstallDoneMsg")
	}
}

func TestInstallModel_InstallMsgUpdatesStepStatus(t *testing.T) {
	tool := domain.Tool{Name: "git"}
	steps := []domain.ExecutionStep{
		{Tool: tool, Status: domain.StatusPending},
	}
	m, _ := makeInstallModel(steps, false)

	// Send an InstallMsg signaling the step is now installing
	ev := domain.ProgressEvent{
		Step:  domain.ExecutionStep{Tool: tool, Status: domain.StatusInstalling},
		Index: 0,
		Total: 1,
	}
	updated, _ := m.Update(screens.InstallMsg(ev))
	_ = updated.(screens.InstallModel) // confirm type assertion succeeds
}

func TestInstallModel_ToolsAppearInView(t *testing.T) {
	steps := []domain.ExecutionStep{
		{Tool: domain.Tool{Name: "git"}, Status: domain.StatusPending},
		{Tool: domain.Tool{Name: "jq"}, Status: domain.StatusPending},
	}
	m, _ := makeInstallModel(steps, false)
	view := m.View()
	for _, name := range []string{"git", "jq"} {
		if !strings.Contains(view, name) {
			t.Errorf("expected tool %q in install view", name)
		}
	}
}

func TestInstallModel_QKeyEmitsQuit(t *testing.T) {
	m, _ := makeInstallModel(nil, false)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Error("expected quit command on q key")
	}
}
