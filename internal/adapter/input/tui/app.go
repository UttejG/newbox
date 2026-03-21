package tui

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/uttejg/newbox/internal/adapter/input/tui/screens"
	"github.com/uttejg/newbox/internal/core/domain"
	"github.com/uttejg/newbox/internal/core/port"
)

type screen int

const (
	screenWelcome screen = iota
	screenProfile
	screenCategories
	screenTools
	screenConfirm
	screenPlanning // Plan() running in background
	screenInstall
	screenDone
)

// planResultMsg carries the outcome of an async Plan() call.
type planResultMsg struct {
	plan *domain.InstallPlan
	err  error
}

// AppModel is the root Bubbletea model. It owns all screen state and handles transitions.
type AppModel struct {
	current screen

	platform       *domain.Platform
	catalogService port.CatalogService
	installSvc     port.InstallService
	dryRun         bool

	// ctx/cancel allow in-progress installs to be cancelled (e.g. on Ctrl+C).
	ctx    context.Context
	cancel context.CancelFunc

	// Screen models (created lazily)
	welcome      screens.WelcomeModel
	profile      screens.ProfileModel
	categories   screens.CategoriesModel
	tools        screens.ToolsModel
	confirm      screens.ConfirmModel
	installModel screens.InstallModel

	// State accumulated across screens
	selectedProfile    *domain.Profile
	selectedCategories []domain.Category
	allCategories      []domain.Category // for confirm display
	selection          *domain.UserSelection

	err error
}

// NewApp creates the root application model.
func NewApp(platform *domain.Platform, catalogSvc port.CatalogService, installSvc port.InstallService, dryRun bool) *AppModel {
	ctx, cancel := context.WithCancel(context.Background())
	return &AppModel{
		current:        screenWelcome,
		platform:       platform,
		catalogService: catalogSvc,
		installSvc:     installSvc,
		dryRun:         dryRun,
		ctx:            ctx,
		cancel:         cancel,
		welcome:        screens.NewWelcome(platform, dryRun),
	}
}

func (m *AppModel) Init() tea.Cmd {
	return m.welcome.Init()
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle global quit — cancel any in-progress install before exiting.
	if km, ok := msg.(tea.KeyMsg); ok && (km.String() == "ctrl+c") {
		m.cancel()
		return m, tea.Quit
	}

	switch m.current {
	case screenWelcome:
		return m.updateWelcome(msg)
	case screenProfile:
		return m.updateProfile(msg)
	case screenCategories:
		return m.updateCategories(msg)
	case screenTools:
		return m.updateTools(msg)
	case screenConfirm:
		return m.updateConfirm(msg)
	case screenPlanning:
		if result, ok := msg.(planResultMsg); ok {
			if result.err != nil {
				m.err = result.err
				return m, nil
			}
			ch := make(chan domain.ProgressEvent, 100)
			m.installModel = screens.NewInstall(result.plan, m.dryRun, ch)
			m.current = screenInstall
			ctx := m.ctx
			go func() {
				defer close(ch)
				_ = m.installSvc.Execute(ctx, result.plan, ch)
			}()
			return m, m.installModel.Init()
		}
		return m, nil
	case screenInstall:
		return m.updateInstall(msg)
	}
	return m, nil
}

func (m *AppModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nError: %v\n\nPress q to quit.\n", m.err)
	}

	switch m.current {
	case screenWelcome:
		return m.welcome.View()
	case screenProfile:
		return m.profile.View()
	case screenCategories:
		return m.categories.View()
	case screenTools:
		return m.tools.View()
	case screenConfirm:
		return m.confirm.View()
	case screenPlanning:
		return "\n  Planning installation…\n"
	case screenInstall:
		return m.installModel.View()
	case screenDone:
		return m.doneView()
	}
	return ""
}

func (m *AppModel) updateWelcome(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	updated, cmd := m.welcome.Update(msg)
	m.welcome = updated.(screens.WelcomeModel)

	if _, ok := msg.(screens.WelcomeDone); ok {
		return m.transitionToProfile()
	}
	return m, cmd
}

func (m *AppModel) transitionToProfile() (tea.Model, tea.Cmd) {
	profiles, err := m.catalogService.GetAllProfiles()
	if err != nil {
		m.err = err
		return m, nil
	}
	m.profile = screens.NewProfile(profiles)
	m.current = screenProfile
	return m, m.profile.Init()
}

func (m *AppModel) updateProfile(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated, cmd := m.profile.Update(msg)
	m.profile = updated.(screens.ProfileModel)

	if ev, ok := msg.(screens.ProfileSelected); ok {
		profile := ev.Profile
		m.selectedProfile = &profile
		return m.transitionToCategories()
	}
	if _, ok := msg.(screens.ProfileBack); ok {
		m.current = screenWelcome
		return m, m.welcome.Init()
	}
	return m, cmd
}

func (m *AppModel) transitionToCategories() (tea.Model, tea.Cmd) {
	cats, err := m.catalogService.GetCategories(m.platform.OS)
	if err != nil {
		m.err = err
		return m, nil
	}
	m.allCategories = cats
	m.categories = screens.NewCategories(cats, m.selectedProfile)
	m.current = screenCategories
	return m, m.categories.Init()
}

func (m *AppModel) updateCategories(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated, cmd := m.categories.Update(msg)
	m.categories = updated.(screens.CategoriesModel)

	if ev, ok := msg.(screens.CategoriesDone); ok {
		m.selectedCategories = ev.Selected
		return m.transitionToTools()
	}
	if _, ok := msg.(screens.CategoriesBack); ok {
		return m.transitionToProfile()
	}
	return m, cmd
}

func (m *AppModel) transitionToTools() (tea.Model, tea.Cmd) {
	if len(m.selectedCategories) == 0 {
		// Skip tools, go straight to confirm with empty selection
		m.selection = &domain.UserSelection{
			Profile:         m.selectedProfile,
			Platform:        m.platform,
			ToolsByCategory: map[string][]domain.Tool{},
		}
		m.confirm = screens.NewConfirm(m.selection, m.allCategories)
		m.current = screenConfirm
		return m, m.confirm.Init()
	}
	m.tools = screens.NewTools(m.selectedCategories)
	m.current = screenTools
	return m, m.tools.Init()
}

func (m *AppModel) updateTools(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated, cmd := m.tools.Update(msg)
	m.tools = updated.(screens.ToolsModel)

	if ev, ok := msg.(screens.ToolsDone); ok {
		m.selection = &domain.UserSelection{
			Profile:         m.selectedProfile,
			Platform:        m.platform,
			ToolsByCategory: ev.ByCategory,
		}
		m.confirm = screens.NewConfirm(m.selection, m.allCategories)
		m.current = screenConfirm
		return m, m.confirm.Init()
	}
	if _, ok := msg.(screens.ToolsBack); ok {
		return m.transitionToCategories()
	}
	return m, cmd
}

func (m *AppModel) updateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated, cmd := m.confirm.Update(msg)
	m.confirm = updated.(screens.ConfirmModel)

	if _, ok := msg.(screens.ConfirmProceed); ok {
		if m.installSvc != nil {
			return m.transitionToInstall()
		}
		m.current = screenDone
		return m, tea.Quit
	}
	if _, ok := msg.(screens.ConfirmBack); ok {
		return m.transitionToTools()
	}
	return m, cmd
}

func (m *AppModel) transitionToInstall() (tea.Model, tea.Cmd) {
	m.current = screenPlanning
	svc := m.installSvc
	sel := m.selection
	ctx := m.ctx
	return m, func() tea.Msg {
		plan, err := svc.Plan(ctx, sel)
		return planResultMsg{plan: plan, err: err}
	}
}

func (m *AppModel) updateInstall(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated, cmd := m.installModel.Update(msg)
	m.installModel = updated.(screens.InstallModel)
	return m, cmd
}

func (m *AppModel) doneView() string {
	return ""
}

// FinalSelection returns the user's selection after the TUI completes.
func (m *AppModel) FinalSelection() *domain.UserSelection {
	return m.selection
}
