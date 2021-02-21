package core

type Record struct {
	// Index allows easy referencing of records (this is used for when writing the struct to file)
	Index     uint
	ParentURL string
	URL       string
	Hostname  string
	depth     uint
	edges     []*Record
}

// Task is what gets sent to the channel for workers to pull data from the web
type Task struct {
	URL   string
	depth uint
}

// URL might be an absolute URL or a relative URL.
type Result struct {
	ParentURL  string
	StatusCode uint
	URLs       []string
	depth      uint
}
