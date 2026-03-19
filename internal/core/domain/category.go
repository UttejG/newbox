package domain

// Category groups related tools under a label shown in the TUI.
type Category struct {
	ID          string
	Name        string // includes emoji, e.g. "💬 Messaging"
	Description string
	Tools       []Tool
}

// FilteredTools returns only tools available on the given OS.
func (c *Category) FilteredTools(os OS) []Tool {
	var result []Tool
	for _, t := range c.Tools {
		if t.IsAvailableOn(os) {
			result = append(result, t)
		}
	}
	return result
}

// IsAvailableOn returns true if the category has at least one tool for the given OS.
func (c *Category) IsAvailableOn(os OS) bool {
	return len(c.FilteredTools(os)) > 0
}
