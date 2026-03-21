package domain

// Profile is a named preset that maps to a set of category IDs.
// When selected in the TUI, its categories are pre-checked.
type Profile struct {
	ID            string
	Name          string // includes emoji, e.g. "🏗️ Developer"
	Description   string
	Categories    []string // category IDs; empty means "all"
	AllCategories bool     // true for the "full" profile
}

// UserSelection holds the final choices the user made in the TUI.
type UserSelection struct {
	Profile         *Profile
	Platform        *Platform         // the detected platform at selection time
	ToolsByCategory map[string][]Tool // category ID → selected tools
}

// AllTools returns a flat list of all selected tools.
func (s *UserSelection) AllTools() []Tool {
	var tools []Tool
	for _, ts := range s.ToolsByCategory {
		tools = append(tools, ts...)
	}
	return tools
}

// TotalCount returns the total number of selected tools.
func (s *UserSelection) TotalCount() int {
	return len(s.AllTools())
}
