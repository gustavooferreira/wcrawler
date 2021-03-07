package graph

import (
	"io"
	"strconv"

	"github.com/gustavooferreira/wcrawler"
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

	rm := wcrawler.NewRecordManager()

	err := rm.LoadFromReader(v.reader)
	if err != nil {
		return err
	}

	records := rm.Dump()

	// Create struct as expected by 3d js library
	nodes := []Node{}
	links := []Link{}

	idMapping := map[int]string{}

	// Nodes
	for _, r := range records {
		node := Node{ID: strconv.Itoa(int(r.Index)), URL: r.URL, Domain: r.Host, LinksCount: r.Edges.Count()}
		nodes = append(nodes, node)

		idMapping[r.Index] = strconv.Itoa(int(r.Index))
	}

	// Add links
	for _, r := range records {
		for _, edge := range r.Edges.Dump() {
			link := Link{
				// ID:     fmt.Sprintf("%d-%d", r.Index, edge),
				Source: idMapping[r.Index],
				Target: idMapping[edge],
			}
			links = append(links, link)
		}
	}

	elements := Elements{Nodes: nodes, Links: links}

	err = GenerateHTML(elements, v.writer)
	if err != nil {
		return err
	}

	return nil
}
