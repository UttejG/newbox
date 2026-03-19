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
