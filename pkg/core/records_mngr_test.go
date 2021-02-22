package core_test

import (
	"testing"

	"github.com/gustavooferreira/wcrawler/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddRecordsToRecordManager(t *testing.T) {

	rm := core.NewRecordManager()

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

	value, err := rm.Get("http://example1.com")
	if err != nil {
		require.Fail(t, err.Error())
	}

	assert.Equal(t, []uint{1, 2}, value.Edges)
}
