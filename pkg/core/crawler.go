package core

import (
	"fmt"
	"net/http"
	"sync"
)

// Crawler brings everything together and is responsible for starting goroutines and manage them.
type Crawler struct {
	connector Connector

	// Read-only vars
	BaseURL      string
	File         string
	Stats        bool
	WorkersCount uint
	Depth        uint

	// Channels
	tasks   chan Task
	results chan Result
}

func NewCrawler(client *http.Client, baseURL string, file string, stats bool, workersCount uint, depth uint) (*Crawler, error) {

	if !IsAbsoluteURL(baseURL) {
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

	// Create channels and WaitGroup
	c.tasks = make(chan Task, int(c.WorkersCount)*2)
	c.results = make(chan Result, int(c.WorkersCount)*2)

	var wg sync.WaitGroup

	// Start merger goroutine (deals with records manager)
	go c.Merger()

	// Start workers (n workers)
	for i := 0; i < int(c.WorkersCount); i++ {
		wg.Add(1)
		go c.WorkerRun(&wg)
	}

	// Start goroutine that handles Ctrl-C (which will close all channels and drain them as well)

	// wait for group
	wg.Wait()

	// If stats == true, then write last message
	fmt.Printf("FINISHED!!\n")
}

// WorkerRun represents the workers doing work in a goroutine
// Receives tasks in a channel and returns results on another
// When tasks channel is closed, the worker returns
func (c *Crawler) WorkerRun(wg *sync.WaitGroup) {
	defer wg.Done()

	end := false

	for {
		select {
		case t, ok := <-c.tasks:
			if !ok {
				end = true
			}

			statusCode, links, err := c.connector.GetLinks(t.URL)

			// handle error
			if err != nil {

			}

			r := Result{ParentURL: t.URL, StatusCode: statusCode, URLs: links, depth: t.depth + 1}

			c.results <- r
		}

		if end {
			break
		}
	}
	// Get the links from the URL

	// if URLs are relative, make sure to join them with ParentURL

	// Returns a list of URLs (and some stats data, like error, status code, etc)
}

func (c *Crawler) Merger() {
	// Add baseURL to tasks channel
	task := Task{URL: c.BaseURL, depth: 0}
	c.tasks <- task

	end := false

	for {
		select {
		case r := <-c.results:
			fmt.Printf("URLs: %+v\n", r.URLs)
			close(c.tasks)
			end = true
		}

		if end {
			break
		}
	}

	// Write to file if enabled
}
