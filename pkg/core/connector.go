package core

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// Connector is responsible to connect to the links and manage connections to websites.
type Connector struct {
	client *http.Client
}

// We want to make sure that we use the same http.Client to reuse connection to get links from other pages beign served by the same server.
// Set timeouts and what not.
func (c *Connector) GetLinks(baseURL string) []string {
	result := []string{}

	resp, err := c.client.Get(baseURL)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()

	// TODO: Make sure body is utf-8 encoded
	links, err := parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("%v\n", links)

	for _, link := range links {

		u, err := url.Parse(link)
		if err != nil {
			fmt.Printf("ERROR: this one was not cool: %s", link)
			continue
		}

		fmt.Printf("URL: %+v\n", *u)
	}

	// body, err := io.ReadAll(resp.Body)
	// fmt.Println(string(body))

	return result
}
