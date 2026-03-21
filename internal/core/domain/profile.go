package domain

import "sort"

// Profile is a named preset that maps to a set of category IDs.
// When selected in the TUI, its categories are pre-checked.
// An empty Categories slice (with AllCategories=false) means no categories
// are pre-selected; AllCategories=true pre-selects every category.
type Profile struct {
	ID            string
	Name          string // includes emoji, e.g. "🏗️ Developer"
	Description   string
	Categories    []string // category IDs pre-selected for this profile
	AllCategories bool     // true for the "full" profile — selects every category
}

// UserSelection holds the final choices the user made in the TUI.
type UserSelection struct {
	Profile         *Profile
	Platform        *Platform         // the detected platform at selection time
	ToolsByCategory map[string][]Tool // category ID → selected tools
}

// AllTools returns a flat list of all selected tools in deterministic order.
func (s *UserSelection) AllTools() []Tool {
	catIDs := make([]string, 0, len(s.ToolsByCategory))
	for id := range s.ToolsByCategory {
		catIDs = append(catIDs, id)
	}
	sort.Strings(catIDs)
	var tools []Tool
	for _, id := range catIDs {
		tools = append(tools, s.ToolsByCategory[id]...)
	}
	return tools
}

// TotalCount returns the total number of selected tools.
func (s *UserSelection) TotalCount() int {
	n := 0
	for _, ts := range s.ToolsByCategory {
		n += len(ts)
	}
	return n
}
