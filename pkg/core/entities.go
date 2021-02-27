package core

// Record represents an entry in the RecordManager (internal state).
type Record struct {
	// Index allows easy referencing of records (used in the edges)
	Index uint `json:"index"`
	// ParentURL holds the URL for the first ParentURL found (there may be more)
	ParentURL  string `json:"parent_url"`
	URL        string `json:"url"`
	Host       string `json:"host"`
	Depth      int    `json:"depth"`
	Edges      []uint `json:"edges"`
	StatusCode int    `json:"statusCode"`
	ErrString  string `json:"errString,omitempty"`
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
	// Host represents the Host portion of the URL
	Host string
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
	URLs       []URLEntity
	// Depth of the child URLs
	Depth int
	Err   error
}
