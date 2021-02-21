package core

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// WebClient is responsible to connect to the links and manage connections to websites.
// Implements Connector interface
type WebClient struct {
	client *http.Client
}

func NewWebClient(client *http.Client) *WebClient {
	return &WebClient{client: client}
}

// We want to make sure that we use the same http.Client to reuse connection to get links from other pages being served by the same server.
// Set timeouts and what not.
func (c *WebClient) GetLinks(baseURL string) (statusCode int, links []string, err error) {
	result := []string{}

	resp, err := c.client.Get(baseURL)
	if err != nil {
		return 0, result, err
	}
	defer resp.Body.Close()

	statusCode = resp.StatusCode

	if statusCode >= 400 {
		return statusCode, result, nil
	}

	// TODO: Make sure body is utf-8 encoded
	rawLinks, err := parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("%v\n", links)

	for _, link := range rawLinks {

		u, err := url.Parse(link)
		if err != nil {
			fmt.Printf("ERROR: this one was not cool: %s", link)
			continue
		}

		// take care of relative links here
		_ = *u

		result = append(result, link)
	}

	// body, err := io.ReadAll(resp.Body)
	// fmt.Println(string(body))

	return statusCode, result, nil
}

func parse(r io.Reader) ([]string, error) {

	result := []string{}

	z := html.NewTokenizer(r)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// EOF
			return result, nil
		case tt == html.StartTagToken:
			t := z.Token()

			// Check if the token is an <a> tag
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			// Extract the href value, if there is one
			ok, url := getHref(t)
			if !ok {
				continue
			}

			// Make sure the url begines in http**
			// hasProto := strings.Index(url, "http") == 0
			// if hasProto {
			// 	ch <- url
			// }

			url = strings.TrimSpace(url)

			result = append(result, url)
		}
	}
}

// Helper function to pull the href attribute from a Token
func getHref(t html.Token) (ok bool, href string) {
	// Iterate over all of the Token's attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}

	// "bare" return will return the variables (ok, href) as defined in
	// the function definition
	return
}

// FindLinks will find all <a> tags with an href attribute.
// Returns a list of links.
// Can also do a bit of clean up before sending the response, like if local URL append to current URL visiting.
func FindLinks() {

}
