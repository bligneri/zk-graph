package assets

import _ "embed"

//go:embed template/force_graph.tmpl
var ForceGraphTemplate string

// GetForceGraphTemplate exposes the embedded template
func GetForceGraphTemplate() string {
	return ForceGraphTemplate
}
