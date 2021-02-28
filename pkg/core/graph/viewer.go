package graph

import (
	"fmt"
	"io"

	"github.com/gustavooferreira/wcrawler/pkg/core"
)

type Viewer struct {
	reader io.Reader
	writer io.Writer
}

func NewViewer(r io.Reader, w io.Writer) *Viewer {
	v := &Viewer{reader: r, writer: w}
	return v
}

func (v *Viewer) Run() error {

	rm := core.NewRecordManager()

	err := rm.LoadFromReader(v.reader)
	if err != nil {
		return err
	}

	records := rm.Dump()

	// Create struct as expected by cytoscape
	elements := []Element{}

	parents := map[string]uint{}
	idMapping := map[uint]string{}

	for _, r := range records {
		elem := Element{Group: "nodes", Data: Data{ID: r.URL, Parent: r.Host}}
		elements = append(elements, elem)

		if _, ok := parents[r.Host]; !ok {
			parents[r.Host] = 0
		}

		idMapping[r.Index] = r.URL
	}

	// Loop through parents
	for p, _ := range parents {
		elem := Element{Group: "nodes", Data: Data{ID: p}}
		elements = append(elements, elem)
	}

	// Add edges
	for _, r := range records {
		for _, edge := range r.Edges {
			elem := Element{
				Group: "edges",
				Data: Data{
					ID:     fmt.Sprintf("%d-%d", r.Index, edge),
					Source: idMapping[r.Index],
					Target: idMapping[edge],
				},
			}
			elements = append(elements, elem)
		}
	}

	err = GenerateHTML(elements, v.writer)
	if err != nil {
		return err
	}

	return nil
}
