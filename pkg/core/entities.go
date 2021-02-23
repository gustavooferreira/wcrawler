package core

// Record represents an entry in the RecordManager (internal state).
type Record struct {
	// Index allows easy referencing of records (used in the edges)
	Index uint `json:"index"`
	// ParentURL holds the URL for the first ParentURL found (there may be more)
	ParentURL string `json:"parent_url"`
	URL       string `json:"url"`
	Host      string `json:"host"`
	Depth     uint   `json:"depth"`
	Edges     []uint `json:"edges"`
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
	URLs       []URLEntity
	Depth      uint
	Err        error
}
