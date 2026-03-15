// Package catalog provides the embedded YAML tool catalog data.
package catalog

import _ "embed"

//go:embed tools.yaml
var ToolsYAML []byte

//go:embed profiles.yaml
var ProfilesYAML []byte
