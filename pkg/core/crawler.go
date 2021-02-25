package core

import (
	"fmt"
	"net/url"
	"sync"
)

// Crawler brings everything together and is responsible for starting goroutines and manage them.
type Crawler struct {
	connector Connector

	// Read-only vars
	InitialURL      string
	File            string
	Stats           bool
	StayInSubdomain bool
	WorkersCount    uint
	Depth           uint

	// Channels
	tasks   chan Task
	results chan Result
}

// NewCrawler returns a new Crawler.
func NewCrawler(connector Connector, initialURL string, file string, stats bool, stayinsubdomain bool, workersCount uint, depth uint) (*Crawler, error) {

	if !IsAbsoluteURL(initialURL) {
		return nil, fmt.Errorf("URL provided is not valid")
	}

	if workersCount == 0 {
		return nil, fmt.Errorf("the number of workers needs to be greater than 0")
	}

	if depth == 0 {
		return nil, fmt.Errorf("recursion depth needs to be greater than 0")
	}

	return &Crawler{
			connector:       connector,
			InitialURL:      initialURL,
			File:            file,
			Stats:           stats,
			StayInSubdomain: stayinsubdomain,
			WorkersCount:    workersCount,
			Depth:           depth},
		nil
}

// Run starts crawling.
func (c *Crawler) Run() {

	// Create channels and WaitGroup
	c.tasks = make(chan Task, int(c.WorkersCount)*2)
	c.results = make(chan Result, int(c.WorkersCount)*2)

	var wg sync.WaitGroup

	// Start stats goroutine
	// sm := NewStatsManager(c.WorkersCount, c.Depth)
	// wg.Add(1)
	// go c.StatsWriter(sm)

	// Start merger goroutine (deals with records manager)
	wg.Add(1)
	go c.Merger()

	// Start workers (n workers)
	for i := 0; i < int(c.WorkersCount); i++ {
		wg.Add(1)
		go c.WorkerRun(&wg)
	}

	// Start goroutine that handles Ctrl-C (which will close all channels and drain them as well)
	// upon receiving sigInt, inform Merger to stop processing any more links.

	// wait for all goroutines to complete
	wg.Wait()
}

// WorkerRun represents the workers crawling links in a goroutine.
// Receives tasks in a channel and returns results on another.
// When tasks channel is closed, the workers return.
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

			r := Result{ParentURL: t.URL, StatusCode: statusCode, URLs: links, Depth: t.Depth + 1}

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

// Merger gets the results from the workers (links) and keeps all the relevant information
// feeding the new links to workers via another channel.
func (c *Crawler) Merger() {

	// Initialize record manager
	rm := NewRecordManager()

	// Add baseURL to tasks channel
	task := Task{URL: c.InitialURL, Depth: 0}
	c.tasks <- task

	// Keep local counter to know what jobs have been done.
	// when counter reaches zero it means there are no more jobs to be processed
	// and the merger can exit.
	// It starts at 1 because we queued the first job already
	jobsCounter := 1

	// keep local queue, where next links will be queued
	//

	end := false

	for {
		select {
		case r := <-c.results:
			// Got a response means we can decrement the job counter.
			jobsCounter--

			// when processing the new links, make sure every time we queue a new link
			// we increase the jobCounter

			for _, uu := range r.URLs {
				rm.Match(uu.Raw)

				u, err := url.Parse(uu.Raw)
				if err != nil {
					continue
				}
				fmt.Printf("URL: %+v\n", u)
			}

			// check depth, if equal or greater then set, then don't queue more

			// check if URLs in rm if not, add them

			// update the stats on what it has found so far to be printed on the screen.

			fmt.Printf("URLs: %+v\n", r.URLs)

			// check if we are done
			if jobsCounter == 0 {
				close(c.tasks)
				end = true
			}
		}

		// select with default to push as many items from the queue to the channel
		// as possible.

		if end {
			break
		}
	}

	// Write to file if enabled
}

// StatsWriter writes stats to a io.Writer (e.g. os.Stdout)
func (c *Crawler) StatsWriter(sm *StatsManager) {
	sm.RunWriter()
}
