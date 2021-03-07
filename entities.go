package wcrawler

import (
	"encoding/json"
	"sort"
)

// Record represents an entry in the RecordManager (internal state).
type Record struct {
	// Index allows easy referencing of records (used in the edges)
	Index int `json:"index"`
	// This indicates whether this is the start of the graph
	// i.e., URL provided.
	InitPoint bool   `json:"initPoint"`
	URL       string `json:"url"`
	Host      string `json:"host"`
	Depth     int    `json:"depth"`
	// Edges      []uint `json:"edges"`
	// This is supposed to be mimicing a hashset
	// We use a struct as a value as it's a bit more space efficient
	Edges      EdgesSet `json:"edges"`
	StatusCode int      `json:"statusCode"`
	ErrString  string   `json:"errString,omitempty"`
}

// RMEntry represents an entry in the RecordManager (external interface).
type RMEntry struct {
	ParentURL  string
	URL        URLEntity
	Depth      int
	StatusCode int
	ErrString  string
}

// URLEntity represents a URL.
type URLEntity struct {
	// NetLoc represents the NetLoc portion of the URL
	NetLoc string
	// Raw represents the entire URL
	Raw string
}

// Task is what gets sent to the channel for workers to pull data from the web.
type Task struct {
	URL   string
	Depth int
}

// Result is what workers return in a channel.
type Result struct {
	ParentURL  string
	StatusCode int
	Links      []URLEntity
	// Depth of the ParentURL
	Depth int
	Err   error
}

type EdgesSet map[int]struct{}

func NewEdgesSet() EdgesSet {
	return make(map[int]struct{})
}

func (es EdgesSet) Add(elems ...int) {
	for _, elem := range elems {
		es[elem] = struct{}{}
	}
}

func (es EdgesSet) Remove(elem int) {
	delete(es, elem)
}

func (es EdgesSet) Count() int {
	return len(es)
}

func (es EdgesSet) Dump() []int {
	arr := []int{}
	for k := range es {
		arr = append(arr, k)
	}

	// sort array before returning
	sort.Ints(arr)

	return arr
}

func (es EdgesSet) MarshalJSON() ([]byte, error) {
	// Create an array with fixed size
	arr := make([]int, len(es))
	index := 0
	for k := range es {
		arr[index] = k
		index++
	}

	// sort array before marshaling
	sort.Ints(arr)

	return json.Marshal(arr)
}

func (es *EdgesSet) UnmarshalJSON(b []byte) error {
	*es = make(map[int]struct{})

	arr := []int{}

	err := json.Unmarshal(b, &arr)
	if err != nil {
		return err
	}

	for _, e := range arr {
		(*es)[e] = struct{}{}
	}
	return nil
}
