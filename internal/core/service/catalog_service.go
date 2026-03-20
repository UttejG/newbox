package service

import (
"fmt"

"github.com/uttejg/newbox/internal/core/domain"
"github.com/uttejg/newbox/internal/core/port"
)

// CatalogService implements port.CatalogService, providing filtered catalog queries.
type CatalogService struct {
provider port.CatalogProvider
}

// NewCatalogService creates a CatalogService backed by the given provider.
func NewCatalogService(provider port.CatalogProvider) *CatalogService {
return &CatalogService{provider: provider}
}

// GetCategories returns categories that have at least one tool available on os.
func (s *CatalogService) GetCategories(os domain.OS) ([]domain.Category, error) {
all, err := s.provider.LoadCategories()
if err != nil {
return nil, fmt.Errorf("loading categories: %w", err)
}

var result []domain.Category
for _, cat := range all {
filtered := cat.FilteredTools(os)
if len(filtered) == 0 {
continue
}
// Return a copy with only OS-compatible tools
result = append(result, domain.Category{
ID:          cat.ID,
Name:        cat.Name,
Description: cat.Description,
Tools:       filtered,
})
}
return result, nil
}

// GetProfile returns the profile with the given ID, or an error if not found.
func (s *CatalogService) GetProfile(id string) (*domain.Profile, error) {
profiles, err := s.provider.LoadProfiles()
if err != nil {
return nil, fmt.Errorf("loading profiles: %w", err)
}

for i := range profiles {
if profiles[i].ID == id {
copy := profiles[i]
return &copy, nil
}
}
return nil, fmt.Errorf("profile %q not found", id)
}

// GetAllProfiles returns all available profiles as a defensive copy.
// The returned slice is independent of any provider-internal cache.
func (s *CatalogService) GetAllProfiles() ([]domain.Profile, error) {
profiles, err := s.provider.LoadProfiles()
if err != nil {
return nil, fmt.Errorf("loading profiles: %w", err)
}
result := make([]domain.Profile, len(profiles))
copy(result, profiles)
return result, nil
}
