package core_test

import (
	"bytes"
	"testing"

	"github.com/gustavooferreira/wcrawler/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddRecords(t *testing.T) {
	rm := core.NewRecordManager()
	addEntries(rm)

	value, err := rm.Get("http://example1.com")
	require.NoError(t, err)

	assert.Equal(t, []uint{1, 2}, value.Edges)
}

func TestSaveToWriter(t *testing.T) {
	expected := `{"http://example1.com":{"index":0,"parent_url":"",` +
		`"url":"http://example1.com","host":"example1.com","depth":0,"edges":[1,2]},` +
		`"http://example1.com/about":{"index":1,"parent_url":"http://example1.com",` +
		`"url":"http://example1.com/about","host":"example1.com","depth":1,"edges":[]},` +
		`"http://example1.com/main":{"index":2,"parent_url":"http://example1.com",` +
		`"url":"http://example1.com/main","host":"example1.com","depth":1,"edges":[3]},` +
		`"http://example123.com/":{"index":3,"parent_url":"http://example1.com/main",` +
		`"url":"http://example123.com/","host":"example123.com","depth":2,"edges":[]}}
`

	rm := core.NewRecordManager()
	addEntries(rm)

	var buf bytes.Buffer
	rm.SaveToWriter(&buf)
	assert.Equal(t, expected, buf.String())
}

func TestLoadFromWriter(t *testing.T) {
	input := `{"http://example1.com":{"index":0,"parent_url":"",` +
		`"url":"http://example1.com","host":"example1.com","depth":0,"edges":[1,2]},` +
		`"http://example1.com/about":{"index":1,"parent_url":"http://example1.com",` +
		`"url":"http://example1.com/about","host":"example1.com","depth":1,"edges":[]},` +
		`"http://example1.com/main":{"index":2,"parent_url":"http://example1.com",` +
		`"url":"http://example1.com/main","host":"example1.com","depth":1,"edges":[3]},` +
		`"http://example123.com/":{"index":3,"parent_url":"http://example1.com/main",` +
		`"url":"http://example123.com/","host":"example123.com","depth":2,"edges":[]}}
`

	rm := core.NewRecordManager()

	var buf bytes.Buffer
	buf.WriteString(input)
	rm.LoadFromReader(&buf)

	r, err := rm.Get("http://example1.com")
	require.NoError(t, err)

	assert.Equal(t, []uint{1, 2}, r.Edges)
}

func addEntries(rm *core.RecordManager) {
	rmEntry1 := core.RMEntry{
		ParentURL: "",
		URL: core.URLEntity{
			Host:   "example1.com",
			String: "http://example1.com",
		},
		Depth: 0,
	}
	rm.AddRecord(rmEntry1)

	rmEntry2 := core.RMEntry{
		ParentURL: "http://example1.com",
		URL: core.URLEntity{
			Host:   "example1.com",
			String: "http://example1.com/about",
		},
		Depth: 1,
	}
	rm.AddRecord(rmEntry2)

	rmEntry3 := core.RMEntry{
		ParentURL: "http://example1.com",
		URL: core.URLEntity{
			Host:   "example1.com",
			String: "http://example1.com/main",
		},
		Depth: 1,
	}
	rm.AddRecord(rmEntry3)

	rmEntry4 := core.RMEntry{
		ParentURL: "http://example1.com/main",
		URL: core.URLEntity{
			Host:   "example123.com",
			String: "http://example123.com/",
		},
		Depth: 2,
	}
	rm.AddRecord(rmEntry4)
}
