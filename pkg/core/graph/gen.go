package graph

import (
	_ "embed"
	"encoding/json"
	"os"
	"text/template"

	"github.com/gustavooferreira/wcrawler/pkg/core"
)

//go:embed "graph.tmpl"
var htmlTemplate string

// GraphGenerator generated an HTML file with the data pulled from the file.
// It uses cytoscape.js library to render the graph in HTML.
type GraphGenerator struct {
	records map[string]core.Record
}

func NewGraphGenerator() *GraphGenerator {
	gg := &GraphGenerator{}
	return gg
}

// LoadInfo takes a RecordManager and loads info into an HTML template.
func (gg *GraphGenerator) LoadInfo() {
	gg.records = make(map[string]core.Record)

	gg.records["test1"] = core.Record{
		Index:      1,
		ParentURL:  "https://www.example.com",
		URL:        "https://www.example2.com",
		Host:       "www.example.com",
		Depth:      2,
		Edges:      []uint{1, 2, 3},
		StatusCode: 200,
		ErrString:  "",
	}
}

// GenerateHTML generates a new HTML file with the loaded data.
func (gg *GraphGenerator) GenerateHTML() {
	// Take the info from gg struct and json indent it to a string

	jsonString, err := json.MarshalIndent(gg.records, "", "    ")
	if err != nil {
		panic(err)
	}

	vars := struct {
		Elements string
	}{
		Elements: string(jsonString),
	}

	t, err := template.New("graph").Parse(htmlTemplate)
	if err != nil {
		panic(err)
	}
	err = t.Execute(os.Stdout, vars)
	if err != nil {
		panic(err)
	}
}

func GetTemplate() string {
	return htmlTemplate
}
