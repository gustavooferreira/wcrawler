package core

// Record represents an entry in the RecordManager (internal state).
type Record struct {
	// Index allows easy referencing of records (this is used for when writing the struct to file)
	Index     uint
	ParentURL string
	URL       string
	Host      string
	Depth     uint
	Edges     []uint
}

// RMEntry represents an entry in the RecordManager (external interface).
type RMEntry struct {
	ParentURL string
	URL       URLEntity
	Depth     uint
}

type URLEntity struct {
	Host   string
	String string
}

// Task is what gets sent to the channel for workers to pull data from the web
type Task struct {
	URL   string
	Depth uint
}

// URL might be an absolute URL or a relative URL.
type Result struct {
	ParentURL  string
	StatusCode int
	URLs       []string
	Depth      uint
	Err        error
}