package core

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// WebClient is responsible to connect to the links and manage connections to websites.
// Implements Connector interface.
type WebClient struct {
	client *http.Client
}

// NewWebClient returns a new WebClient.
func NewWebClient(client *http.Client) *WebClient {
	return &WebClient{client: client}
}

// GetLinks returns all the links found in the webpage.
func (c *WebClient) GetLinks(rawURL string) (statusCode int, links []URLEntity, err error) {
	// make sure to use the same http.Client to reuse connections to get links
	// from other pages being served by the same server.
	// Check for robot.txt, maybe?

	resp, err := c.client.Get(rawURL)
	if err != nil {
		return 0, links, err
	}
	defer resp.Body.Close()

	statusCode = resp.StatusCode

	if statusCode >= 200 && statusCode < 300 {
		return statusCode, links, nil
	}

	links, err = c.parse(resp.Body)

	return statusCode, links, err
}

// parse parses the webpage looking for links.
func (c *WebClient) parse(r io.Reader) (links []URLEntity, err error) {
	// Parse <base> tag if it exists
	// Parse all <a> tags
	// Cater for the fact that a <a> link might be a mailto or a phone or something else.
	// Check if "Opaque" field in URL struct is set
	// Validate whether they are absolute or relative tags. Also check if the relative tags start with a /
	// same webpages will have the <base> tag inside <head> which should be used when joining relative URLs.

	// TODO: Make sure body is utf-8 encoded

	insideHead := false
	baseURL := ""

	links = []URLEntity{}

	z := html.NewTokenizer(r)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// EOF
			return links, nil
		case tt == html.StartTagToken:
			t := z.Token()

			if t.Data == "head" {
				insideHead = true
			}

			if t.Data == "base" && insideHead == true {
				ok, rawURL := getHref(t)
				if ok {
					baseURL = rawURL
				}
			}

			// Check if the token is an <a> tag
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			// Extract the href value, if there is one
			ok, rawURL := getHref(t)
			if !ok {
				continue
			}

			rawURL = strings.TrimSpace(rawURL)
			// Check type of link
			// if relative, join it to base
			_ = baseURL

			u, err := url.Parse(rawURL)
			if err != nil {

			}

			// if link is relative, need to join with base url first before doing this.
			base := fmt.Sprintf("%s://%s", u.Scheme, u.Host)

			links = append(links, URLEntity{Base: base, Raw: rawURL})

		case tt == html.EndTagToken:
			t := z.Token()
			if t.Data == "head" {
				insideHead = false
			}

		}
	}
}

// getHref returns the href attribute from a Token
func getHref(t html.Token) (ok bool, href string) {
	// Iterate over all of the Token's attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	return
}

// for _, link := range rawLinks {

// 	u, err := url.Parse(link)
// 	if err != nil {
// 		fmt.Printf("ERROR: this one was not cool: %s", link)
// 		continue
// 	}

// 	// TODO: take care of relative links here

// 	l := URLEntity{Host: u.Host, Raw: link}
// 	result = append(result, l)
// }
