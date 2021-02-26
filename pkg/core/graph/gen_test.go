package graph_test

import (
	"testing"

	"github.com/gustavooferreira/wcrawler/pkg/core/graph"
	"github.com/stretchr/testify/assert"
)

func TestHTMLTemplateString(t *testing.T) {
	// s := graph.GetTemplate()
	assert.Equal(t, "", "")
}

func TestGeneration(t *testing.T) {

	gg := graph.NewGraphGenerator()
	gg.LoadInfo()
	// gg.GenerateHTML()

	assert.Equal(t, "", "")
}
