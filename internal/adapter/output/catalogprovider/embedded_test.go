package catalogprovider_test

import (
	"testing"

	"github.com/uttejg/newbox/internal/adapter/output/catalogprovider"
)

func TestEmbeddedProvider_LoadCategories(t *testing.T) {
	p := &catalogprovider.EmbeddedProvider{}
	cats, err := p.LoadCategories()
	if err != nil {
		t.Fatalf("LoadCategories() error: %v", err)
	}
	if len(cats) == 0 {
		t.Fatal("LoadCategories() returned no categories")
	}
	for _, cat := range cats {
		if cat.ID == "" {
			t.Errorf("category missing ID: %+v", cat)
		}
		if cat.Name == "" {
			t.Errorf("category %q missing Name", cat.ID)
		}
		if len(cat.Tools) == 0 {
			t.Errorf("category %q has no tools", cat.ID)
		}
		for _, tool := range cat.Tools {
			if tool.Name == "" {
				t.Errorf("tool in %q missing Name", cat.ID)
			}
			// Every tool must have at least one OS mapping
			if tool.MacOS == nil && tool.Windows == nil && tool.Linux == nil {
				t.Errorf("tool %q in %q has no OS mappings", tool.Name, cat.ID)
			}
		}
	}
}

func TestEmbeddedProvider_LoadProfiles(t *testing.T) {
	p := &catalogprovider.EmbeddedProvider{}
	profiles, err := p.LoadProfiles()
	if err != nil {
		t.Fatalf("LoadProfiles() error: %v", err)
	}
	if len(profiles) == 0 {
		t.Fatal("LoadProfiles() returned no profiles")
	}

	// Must include the standard profiles
	requiredIDs := []string{"developer", "creative", "minimal", "full", "custom"}
	found := make(map[string]bool)
	for _, p := range profiles {
		found[p.ID] = true
	}
	for _, id := range requiredIDs {
		if !found[id] {
			t.Errorf("missing required profile %q", id)
		}
	}
}

func TestCatalog_AllToolsHaveAtLeastOneOS(t *testing.T) {
	p := &catalogprovider.EmbeddedProvider{}
	cats, err := p.LoadCategories()
	if err != nil {
		t.Fatalf("LoadCategories() error: %v", err)
	}
	for _, cat := range cats {
		for _, tool := range cat.Tools {
			hasAny := tool.MacOS != nil || tool.Windows != nil || tool.Linux != nil
			if !hasAny {
				t.Errorf("tool %q in category %q has no OS mappings", tool.Name, cat.ID)
			}
		}
	}
}

func TestCatalog_NoDuplicateToolNames(t *testing.T) {
	p := &catalogprovider.EmbeddedProvider{}
	cats, err := p.LoadCategories()
	if err != nil {
		t.Fatalf("LoadCategories() error: %v", err)
	}
	// Track tool name for global uniqueness across all categories
	seen := make(map[string]string)
	for _, cat := range cats {
		for _, tool := range cat.Tools {
			if prevCat, exists := seen[tool.Name]; exists {
				t.Errorf("duplicate tool %q: first seen in category %q, also in %q", tool.Name, prevCat, cat.ID)
			}
			seen[tool.Name] = cat.ID
		}
	}
}

func TestCatalog_AllCategoryIDsUnique(t *testing.T) {
	p := &catalogprovider.EmbeddedProvider{}
	cats, err := p.LoadCategories()
	if err != nil {
		t.Fatalf("LoadCategories() error: %v", err)
	}
	seen := make(map[string]bool)
	for _, cat := range cats {
		if seen[cat.ID] {
			t.Errorf("duplicate category ID %q", cat.ID)
		}
		seen[cat.ID] = true
	}
}

func TestCatalog_CategoryCount(t *testing.T) {
	p := &catalogprovider.EmbeddedProvider{}
	cats, err := p.LoadCategories()
	if err != nil {
		t.Fatalf("LoadCategories() error: %v", err)
	}
	if len(cats) < 20 {
		t.Errorf("expected at least 20 categories, got %d", len(cats))
	}
}

func TestCatalog_ToolCount(t *testing.T) {
	p := &catalogprovider.EmbeddedProvider{}
	cats, err := p.LoadCategories()
	if err != nil {
		t.Fatalf("LoadCategories() error: %v", err)
	}
	total := 0
	for _, cat := range cats {
		total += len(cat.Tools)
	}
	if total < 80 {
		t.Errorf("expected at least 80 tools, got %d", total)
	}
}
