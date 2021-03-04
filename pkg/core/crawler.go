package core

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/oleiade/lane"
)

// Crawler brings everything together and is responsible for starting goroutines and manage them.
type Crawler struct {
	connector Connector

	// Read-only vars
	InitialURL      string
	IOWriter        io.Writer
	Stats           bool
	WorkersCount    int
	Depth           int
	StayInSubdomain bool
	SubDomain       string
	Retry           int

	statsManager StatsManager

	// Channels
	tasks   chan Task
	results chan Result
}

// NewCrawler returns a new Crawler.
func NewCrawler(connector Connector, initialURL string, retry int, ioWriter io.Writer, stats bool, stayinsubdomain bool, workersCount int, depth int) (*Crawler, error) {

	urlEntity, err := ExtractURL(initialURL)
	if err != nil {
		return nil, fmt.Errorf("URL has to be an absolute URL (including scheme)")
	}

	if workersCount == 0 {
		return nil, fmt.Errorf("the number of workers needs to be greater than 0")
	}

	if depth < 0 {
		return nil, fmt.Errorf("recursion depth needs to be greater or equal to 0")
	}

	return &Crawler{
			connector:       connector,
			InitialURL:      urlEntity.Raw,
			IOWriter:        ioWriter,
			Stats:           stats,
			WorkersCount:    workersCount,
			Depth:           depth,
			StayInSubdomain: stayinsubdomain,
			SubDomain:       urlEntity.NetLoc,
			Retry:           retry},
		nil
}

// Run starts crawling.
func (c *Crawler) Run() {

	// Create channels and WaitGroup
	c.tasks = make(chan Task, int(c.WorkersCount)*2)
	c.results = make(chan Result, int(c.WorkersCount)*2)

	var wg sync.WaitGroup

	// Start stats goroutine
	if c.Stats {
		sm := NewStatsManager(c.WorkersCount, c.Depth)
		wg.Add(1)
		c.statsManager = sm
		go c.StatsWriter(&wg)

		c.statsManager.UpdateStats(SetAppState(AppState_Running))
	}

	// Start merger goroutine (deals with records manager)
	wg.Add(1)
	go c.Merger(&wg)

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

	for t := range c.tasks {
		if c.Stats {
			c.statsManager.UpdateStats(IncDecWorkersRunning(1))
		}

		var statusCode int
		var links []URLEntity
		var latency time.Duration
		var err error

		// retry
		for i := 0; i <= c.Retry; i++ {
			statusCode, links, latency, err = c.connector.GetLinks(t.URL)
			if err == nil {
				// Only retry if timeout!
				break
			}
		}

		r := Result{
			ParentURL:  t.URL,
			StatusCode: statusCode,
			Links:      links,
			Depth:      t.Depth,
			Err:        err,
		}

		c.results <- r

		if c.Stats {
			c.statsManager.UpdateStats(IncDecWorkersRunning(-1),
				IncDecTotalRequestsCount(1),
				AddLatencySample(latency))
		}
	}
}

// Merger gets the results from the workers (links) and keeps all the relevant information
// feeding the new links to workers via another channel.
func (c *Crawler) Merger(wg *sync.WaitGroup) {
	defer wg.Done()

	// Keep local counter to know what jobs have been done.
	// when counter reaches zero it means there are no more jobs to be processed
	// and the merger can exit.
	jobsCounter := 0
	var err error

	// Create queue for queuing jobs
	queue := lane.NewQueue()

	// Initialize record manager
	rm := NewRecordManager()

	// Add baseURL to tasks channel
	task := Task{URL: c.InitialURL, Depth: 0}
	c.tasks <- task

	// Add baseURL as an entry to Record Manager
	re := RMEntry{ParentURL: "", URL: URLEntity{NetLoc: c.SubDomain, Raw: c.InitialURL}, Depth: 0}
	rm.AddRecord(re)

	jobsCounter++

	if c.Stats {
		c.statsManager.UpdateStats(SetLinksInQueue(jobsCounter))
	}

	// ---------

	for {
		r := <-c.results
		// Got a response means we can decrement the job counter
		jobsCounter--

		// Update parent URL entry in Record Manager
		err = rm.Update(r.ParentURL, r.StatusCode, r.Err)
		if err != nil {
			// log
			// continue
		}

		// when processing the new links, make sure every time we queue a new link
		// we increase the jobCounter

		// Check which new jobs to queue
		// Check depth, if equal or greater then set, then don't queue more
		// Also check that we didn't get an error or an unexpected status code
		// If Depth is equal to zero then don't stop ever.
		if (r.Depth < c.Depth || c.Depth == 0) && r.Err == nil && r.StatusCode >= 200 && r.StatusCode < 300 {
			for _, uu := range r.Links {
				if c.StayInSubdomain && c.SubDomain != uu.NetLoc {
					continue
				}

				// if already in the cache, we don't want to query it again
				if !rm.Exists(uu.Raw) {
					queue.Enqueue(Task{URL: uu.Raw, Depth: r.Depth + 1})
					jobsCounter++

					rme := RMEntry{ParentURL: r.ParentURL, URL: uu, Depth: r.Depth + 1}
					rm.AddRecord(rme)
				}
			}
		}

		if c.Stats {
			errCount := 0
			if r.Err != nil || r.StatusCode < 200 || r.StatusCode >= 300 {
				errCount = 1
			}

			c.statsManager.UpdateStats(
				SetLinksInQueue(jobsCounter),
				SetLinksCount(rm.Count()),
				SetDepth(r.Depth),
				IncDecErrorsCount(errCount))
		}

		// fill tasks channel until either channel blocks or queue is empty
		for {
			// Check if channel is full
			// This is fine because this goroutine is the only one writing to the channel,
			// so it won't block when we actually try to write to the channel.
			// If it says the channel is full and the very next millisecond it's not,
			// there is no problem as we will come back to this to refill it.
			if len(c.tasks) == cap(c.tasks) {
				break
			}

			// Check if we can dequeue an item from the queue, if yes, try to push it to the channel
			if queue.Empty() {
				// No more items to dequeue
				break
			} else {
				c.tasks <- queue.Dequeue().(Task)
			}
		}

		// check if we are done (i.e., no more jobs)
		if jobsCounter == 0 {
			close(c.tasks)
			break
		}
	}

	// Write to file
	err = rm.SaveToWriter(c.IOWriter, true)
	if err != nil {
		// log
	}

	if c.Stats {
		c.statsManager.UpdateStats(SetAppState(AppState_Finished))
	}
}

// StatsWriter writes stats to a io.Writer (e.g. os.Stdout)
func (c *Crawler) StatsWriter(wg *sync.WaitGroup) {
	defer wg.Done()
	c.statsManager.RunOutputFlusher()
}
