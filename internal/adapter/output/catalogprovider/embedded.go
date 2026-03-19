// Package catalogprovider loads the embedded YAML tool catalog.
package catalogprovider

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/uttejg/newbox/catalog"
	"github.com/uttejg/newbox/internal/core/domain"
)

// yamlTool mirrors the YAML structure for a single tool entry.
type yamlTool struct {
	Name            string      `yaml:"name"`
	Description     string      `yaml:"description"`
	Website         string      `yaml:"website"`
	DotfilesDefault bool        `yaml:"dotfiles_default"`
	MacOS           *yamlPkgRef `yaml:"macos"`
	Windows         *yamlPkgRef `yaml:"windows"`
	Linux           *yamlPkgRef `yaml:"linux"`
}

type yamlPkgRef struct {
	Formula string `yaml:"formula"`
	Cask    string `yaml:"cask"`
	MAS     string `yaml:"mas"`
	Winget  string `yaml:"winget"`
	Apt     string `yaml:"apt"`
	Dnf     string `yaml:"dnf"`
	Pacman  string `yaml:"pacman"`
	Flatpak string `yaml:"flatpak"`
}

type yamlCategory struct {
	ID          string     `yaml:"id"`
	Name        string     `yaml:"name"`
	Description string     `yaml:"description"`
	Tools       []yamlTool `yaml:"tools"`
}

type yamlCatalog struct {
	Categories []yamlCategory `yaml:"categories"`
}

type yamlProfile struct {
	ID            string   `yaml:"id"`
	Name          string   `yaml:"name"`
	Description   string   `yaml:"description"`
	AllCategories bool     `yaml:"all_categories"`
	Categories    []string `yaml:"categories"`
}

type yamlProfiles struct {
	Profiles []yamlProfile `yaml:"profiles"`
}

// EmbeddedProvider implements port.CatalogProvider using go:embed YAML files.
type EmbeddedProvider struct{}

// LoadCategories parses tools.yaml and returns all categories.
func (p *EmbeddedProvider) LoadCategories() ([]domain.Category, error) {
	var raw yamlCatalog
	if err := yaml.Unmarshal(catalog.ToolsYAML, &raw); err != nil {
		return nil, fmt.Errorf("parsing embedded tools.yaml: %w", err)
	}

	categories := make([]domain.Category, 0, len(raw.Categories))
	for _, rc := range raw.Categories {
		cat := domain.Category{
			ID:          rc.ID,
			Name:        rc.Name,
			Description: rc.Description,
			Tools:       make([]domain.Tool, 0, len(rc.Tools)),
		}
		for _, rt := range rc.Tools {
			tool := domain.Tool{
				Name:            rt.Name,
				Description:     rt.Description,
				Website:         rt.Website,
				DotfilesDefault: rt.DotfilesDefault,
				MacOS:           toPackageRef(rt.MacOS),
				Windows:         toPackageRef(rt.Windows),
				Linux:           toPackageRef(rt.Linux),
			}
			cat.Tools = append(cat.Tools, tool)
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

// LoadProfiles parses profiles.yaml and returns all profiles.
func (p *EmbeddedProvider) LoadProfiles() ([]domain.Profile, error) {
	var raw yamlProfiles
	if err := yaml.Unmarshal(catalog.ProfilesYAML, &raw); err != nil {
		return nil, fmt.Errorf("parsing embedded profiles.yaml: %w", err)
	}

	profiles := make([]domain.Profile, 0, len(raw.Profiles))
	for _, rp := range raw.Profiles {
		profiles = append(profiles, domain.Profile{
			ID:            rp.ID,
			Name:          rp.Name,
			Description:   rp.Description,
			AllCategories: rp.AllCategories,
			Categories:    rp.Categories,
		})
	}
	return profiles, nil
}

func toPackageRef(r *yamlPkgRef) *domain.PackageRef {
	if r == nil {
		return nil
	}
	ref := &domain.PackageRef{
		Formula: r.Formula,
		Cask:    r.Cask,
		MAS:     r.MAS,
		Winget:  r.Winget,
		Apt:     r.Apt,
		Dnf:     r.Dnf,
		Pacman:  r.Pacman,
		Flatpak: r.Flatpak,
	}
	if ref.IsEmpty() {
		return nil
	}
	return ref
}
