package graph

//go:embed "graph.tmpl"
var htmlTemplate string

// GraphGenerator generated an HTML file with the data pulled from the file.
// It uses cytoscape.js library to render the graph in HTML.
type GraphGenerator struct {
}
