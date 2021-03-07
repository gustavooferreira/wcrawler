package wcrawler_test

import (
	"bytes"
	"testing"

	"github.com/gustavooferreira/wcrawler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddRecords(t *testing.T) {
	rm := wcrawler.NewRecordManager()
	addEntries(rm)

	value, ok := rm.Get("http://example1.com")

	require.Equal(t, true, ok)

	expectedES := wcrawler.NewEdgesSet()
	expectedES.Add(1, 2)
	assert.Equal(t, expectedES, value.Edges)
}

func TestSaveToWriter(t *testing.T) {
	expected := `{"http://example1.com":{"index":0,"initPoint":true,"url":"http://example1.com",` +
		`"host":"example1.com","depth":0,"edges":[1,2],"statusCode":200},` +
		`"http://example1.com/about":{"index":1,"initPoint":false,"url":"http://example1.com/about",` +
		`"host":"example1.com","depth":1,"edges":[],"statusCode":200},` +
		`"http://example1.com/main":{"index":2,"initPoint":false,"url":"http://example1.com/main",` +
		`"host":"example1.com","depth":1,"edges":[3],"statusCode":200},` +
		`"http://example123.com/":{"index":3,"initPoint":false,"url":"http://example123.com/",` +
		`"host":"example123.com","depth":2,"edges":[],"statusCode":200}}
`

	rm := wcrawler.NewRecordManager()
	addEntries(rm)

	var buf bytes.Buffer
	rm.SaveToWriter(&buf, false)
	assert.Equal(t, expected, buf.String())
}

func TestLoadFromWriter(t *testing.T) {
	input := `{"http://example1.com":{"index":0,"initPoint":true,"url":"http://example1.com",` +
		`"host":"example1.com","depth":0,"edges":[1,2],"statusCode":200},` +
		`"http://example1.com/about":{"index":1,"initPoint":false,"url":"http://example1.com/about",` +
		`"host":"example1.com","depth":1,"edges":[],"statusCode":200},` +
		`"http://example1.com/main":{"index":2,"initPoint":false,"url":"http://example1.com/main",` +
		`"host":"example1.com","depth":1,"edges":[3],"statusCode":200},` +
		`"http://example123.com/":{"index":3,"initPoint":false,"url":"http://example123.com/",` +
		`"host":"example123.com","depth":2,"edges":[],"statusCode":200}}
`

	rm := wcrawler.NewRecordManager()

	var buf bytes.Buffer
	buf.WriteString(input)
	rm.LoadFromReader(&buf)

	value, ok := rm.Get("http://example1.com")
	require.Equal(t, true, ok)

	expectedES := wcrawler.NewEdgesSet()
	expectedES.Add(1, 2)
	assert.Equal(t, expectedES, value.Edges)

	assert.Equal(t, true, value.InitPoint)

	value, ok = rm.Get("http://example1.com/about")
	require.Equal(t, true, ok)
	assert.Equal(t, false, value.InitPoint)
}

func addEntries(rm *wcrawler.RecordManager) {
	rmEntry1 := wcrawler.RMEntry{
		ParentURL: "",
		URL: wcrawler.URLEntity{
			NetLoc: "example1.com",
			Raw:    "http://example1.com",
		},
		Depth:      0,
		StatusCode: 200,
		ErrString:  "",
	}
	rm.AddRecord(rmEntry1)

	rmEntry2 := wcrawler.RMEntry{
		ParentURL: "http://example1.com",
		URL: wcrawler.URLEntity{
			NetLoc: "example1.com",
			Raw:    "http://example1.com/about",
		},
		Depth:      1,
		StatusCode: 200,
		ErrString:  "",
	}
	rm.AddRecord(rmEntry2)

	rmEntry3 := wcrawler.RMEntry{
		ParentURL: "http://example1.com",
		URL: wcrawler.URLEntity{
			NetLoc: "example1.com",
			Raw:    "http://example1.com/main",
		},
		Depth:      1,
		StatusCode: 200,
		ErrString:  "",
	}
	rm.AddRecord(rmEntry3)

	rmEntry4 := wcrawler.RMEntry{
		ParentURL: "http://example1.com/main",
		URL: wcrawler.URLEntity{
			NetLoc: "example123.com",
			Raw:    "http://example123.com/",
		},
		Depth:      2,
		StatusCode: 200,
		ErrString:  "",
	}
	rm.AddRecord(rmEntry4)
}
