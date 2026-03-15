package port

import "github.com/uttejg/newbox/internal/core/domain"

// CatalogProvider is the output port for loading catalog data.
// Implemented by the embedded YAML loader.
type CatalogProvider interface {
	// LoadCategories returns all categories and their tools.
	LoadCategories() ([]domain.Category, error)
	// LoadProfiles returns all available profiles.
	LoadProfiles() ([]domain.Profile, error)
}

// CatalogService is the input port for querying catalog data.
// Called by the TUI and CLI adapters.
type CatalogService interface {
	// GetCategories returns categories that have at least one tool for the given OS.
	GetCategories(os domain.OS) ([]domain.Category, error)
	// GetProfile returns a profile by ID.
	GetProfile(id string) (*domain.Profile, error)
	// GetAllProfiles returns all available profiles.
	GetAllProfiles() ([]domain.Profile, error)
}
