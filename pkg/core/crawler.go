package core

import (
	"fmt"
	"net/http"
	"net/url"
)

// Crawler brings everything together and is responsible for starting goroutines and manage them.
type Crawler struct {
	connector    Connector
	BaseURL      string
	File         string
	Stats        bool
	WorkersCount uint
	Depth        uint
}

func NewCrawler(client *http.Client, baseURL string, file string, stats bool, workersCount uint, depth uint) (*Crawler, error) {

	if !IsURL(baseURL) {
		return nil, fmt.Errorf("URL provided is not valid")
	}

	if workersCount == 0 {
		return nil, fmt.Errorf("the number of workers needs to be greater or equal to 1")
	}

	if depth == 0 {
		return nil, fmt.Errorf("recursion depth needs to be greater or equal to 1")
	}

	connector := WebClient{client: client}
	return &Crawler{connector: &connector, BaseURL: baseURL, File: file, Stats: stats, WorkersCount: workersCount, Depth: depth}, nil
}

func (c *Crawler) Run() {

	// Create tasks channel
	tasks := make(chan Task, int(c.WorkersCount)*2)
	results := make(chan Result, int(c.WorkersCount)*2)

	// Start merger goroutine (deals with records manager)

	// Start workers (n workers as given by the cli option)

	for i := 0; i < int(c.WorkersCount); i++ {
		go c.WorkerRun(tasks, results)
	}

	c.connector.GetLinks(c.BaseURL)

	// Start goroutine that handles Ctrl-C (which will close all channels and drain them as well)

	// wait for group
}

// WorkerRun represents the workers doing work in a goroutine
// Receives tasks in a channel and returns results on another
// When tasks channel is closed, the worker returns
func (c *Crawler) WorkerRun(tasks <-chan Task, results chan<- Result) {
	// Get the links from the URL

	// if URLs are relative, make sure to join them with ParentURL

	// Returns a list of URLs (and some stats data, like error, status code, etc)
}

// IsURL validates URL
func IsURL(urlS string) bool {
	u, err := url.Parse(urlS)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}
