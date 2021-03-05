package core_test

import (
	"testing"

	"github.com/gustavooferreira/wcrawler/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEdgesSetCreation(t *testing.T) {

	set := core.NewEdgesSet()

	require.NotNil(t, set, "EdgesSet must not be nil")

	length := len(set)

	assert.Equal(t, 0, length)
}

func TestEdgesSetUninitializedOperations(t *testing.T) {
	assert := assert.New(t)

	var set core.EdgesSet

	// testing len function on type
	length := len(set)
	assert.Equal(0, length)

	// testing adding elements
	testFuncAdd := func() {
		set.Add(1)
	}
	assert.Panics(testFuncAdd, "Add function should panic on an uninitialized EdgesSet")

	// test removing an element
	testFuncRemove := func() {
		set.Remove(1)
	}
	assert.NotPanics(testFuncRemove, "Remove function should not panic even with an uninitialized EdgesSet")

	// test counting elements
	count := set.Count()
	assert.Equal(0, count, "count of uninitialized EdgesSet should be zero")

	// test dumping all elements
	dump := set.Dump()
	assert.Equal([]int{}, dump, "dump should return an empty array, on an uninitialized EdgesSet")

	// testing marshaling
	bytes, err := set.MarshalJSON()
	require.NoError(t, err, "expected empty array '[]'")
	assert.Equal("[]", string(bytes))

	// testing unmarshaling
	jsonString := `[1,2,3,4]`
	err = set.UnmarshalJSON([]byte(jsonString))
	require.NoError(t, err, "unmarsheling into an uninitialized EdgesSet should initialize the type and populate its values")

	assert.Equal([]int{1, 2, 3, 4}, set.Dump())
}

func TestEdgesSetOperations(t *testing.T) {
	assert := assert.New(t)

	// set := core.EdgesSet{}
	set := core.NewEdgesSet()

	set.Add(1, 2, 3, 4, 5, 6)

	// testing len function on type
	length := len(set)
	assert.Equal(6, length)

	// testing Count method
	count := set.Count()
	assert.Equal(6, count)

	// testing removing an element
	set.Remove(3)
	assert.Equal([]int{1, 2, 4, 5, 6}, set.Dump())

	// testing marshaling
	bytes, err := set.MarshalJSON()
	require.NoError(t, err, "expected no error marsheling an initialized type")
	assert.Equal("[1,2,4,5,6]", string(bytes))

	// testing unmarshaling
	jsonString := `[10,20,30,40]`
	err = set.UnmarshalJSON([]byte(jsonString))
	require.NoError(t, err)
	assert.Equal([]int{10, 20, 30, 40}, set.Dump(), "unmarsheling should override any values already in the type")

	// testing unmarshaling fail
	jsonString = `[10,20,30,40,"50"]`
	err = set.UnmarshalJSON([]byte(jsonString))
	require.Error(t, err, "mistypes should return an error")
}
