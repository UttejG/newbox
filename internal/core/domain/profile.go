package domain

import "sort"

// Profile is a named preset that maps to a set of category IDs.
// When selected in the TUI, its categories are pre-checked.
type Profile struct {
	ID          string
	Name        string // includes emoji, e.g. "🏗️ Developer"
	Description string
	Categories  []string // category IDs to pre-check; nil/empty means none pre-selected
	AllCategories bool   // when true, all categories are pre-checked (overrides Categories)
}

// UserSelection holds the final choices the user made in the TUI.
type UserSelection struct {
	Profile     *Profile
	ToolsByCategory map[string][]Tool // category ID → selected tools
}

// AllTools returns a flat list of all selected tools in deterministic category order.
func (s *UserSelection) AllTools() []Tool {
	keys := make([]string, 0, len(s.ToolsByCategory))
	for k := range s.ToolsByCategory {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var tools []Tool
	for _, k := range keys {
		tools = append(tools, s.ToolsByCategory[k]...)
	}
	return tools
}

// TotalCount returns the total number of selected tools.
func (s *UserSelection) TotalCount() int {
	return len(s.AllTools())
}
