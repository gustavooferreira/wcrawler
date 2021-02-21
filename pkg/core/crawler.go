package core

import (
	"net/http"
)

// Crawler brings everything together and is responsible for starting goroutines and manage them.
type Crawler struct {
	client       *http.Client
	file         string
	stats        bool
	workersCount uint
	depth        uint
}

func NewCrawler(client *http.Client, file string, stats bool, workersCount uint, depth uint) *Crawler {
	return &Crawler{client: client, file: file, stats: stats, workersCount: workersCount, depth: depth}
}

func (c *Crawler) Run(baseURL string) {

	connector := Connector{client: c.client}
	connector.GetLinks(baseURL)

	// Start workers (n workers as given by the cli option)

	// Start merger goroutine

	// Start goroutine that handles Ctrl-C (which will close all channels and drain them as well)

	// wait for group
}
