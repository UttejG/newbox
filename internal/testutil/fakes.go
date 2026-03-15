// Package testutil provides shared fakes and builders for unit tests.
package testutil

import (
	"fmt"

	"github.com/uttejg/newbox/internal/core/domain"
)

// FakeCatalogProvider implements port.CatalogProvider with in-memory data.
type FakeCatalogProvider struct {
	Categories []domain.Category
	Profiles   []domain.Profile
	Err        error // if non-nil, all calls return this error
}

func (f *FakeCatalogProvider) LoadCategories() ([]domain.Category, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	return f.Categories, nil
}

func (f *FakeCatalogProvider) LoadProfiles() ([]domain.Profile, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	return f.Profiles, nil
}

// NewTestTool creates a Tool available on the given OSes for testing.
func NewTestTool(name string, oses ...domain.OS) domain.Tool {
	t := domain.Tool{Name: name, Description: "test tool " + name}
	for _, os := range oses {
		ref := &domain.PackageRef{Formula: "test-" + name}
		switch os {
		case domain.OSMacOS:
			t.MacOS = ref
		case domain.OSLinux:
			t.Linux = ref
		case domain.OSWindows:
			t.Windows = ref
		}
	}
	return t
}

// NewTestCategory creates a Category with the given tools.
func NewTestCategory(id, name string, tools ...domain.Tool) domain.Category {
	return domain.Category{
		ID:          id,
		Name:        name,
		Description: "test category " + id,
		Tools:       tools,
	}
}

// NewTestProfile creates a Profile for testing.
func NewTestProfile(id string, categories ...string) domain.Profile {
	return domain.Profile{
		ID:          id,
		Name:        fmt.Sprintf("Test Profile %s", id),
		Description: "test profile",
		Categories:  categories,
	}
}

// DefaultTestCatalog returns a small catalog useful in most tests.
func DefaultTestCatalog() ([]domain.Category, []domain.Profile) {
	cats := []domain.Category{
		NewTestCategory("messaging", "💬 Messaging",
			NewTestTool("Signal", domain.OSMacOS, domain.OSLinux),
			NewTestTool("Telegram", domain.OSMacOS, domain.OSLinux, domain.OSWindows),
			NewTestTool("WinOnly", domain.OSWindows),
		),
		NewTestCategory("browsers", "🌐 Browsers",
			NewTestTool("Firefox", domain.OSMacOS, domain.OSLinux, domain.OSWindows),
			NewTestTool("Arc", domain.OSMacOS), // macOS only
		),
		NewTestCategory("cli", "🔧 CLI Essentials",
			NewTestTool("git", domain.OSMacOS, domain.OSLinux, domain.OSWindows),
			NewTestTool("jq", domain.OSMacOS, domain.OSLinux, domain.OSWindows),
		),
	}
	profiles := []domain.Profile{
		NewTestProfile("developer", "messaging", "browsers", "cli"),
		NewTestProfile("minimal", "cli"),
		{ID: "full", Name: "🚀 Full", Description: "Everything", AllCategories: true},
		{ID: "custom", Name: "🔧 Custom", Description: "Pick your own"},
	}
	return cats, profiles
}
