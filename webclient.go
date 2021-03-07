package wcrawler

import (
	"io"
	"net/http"
	"net/http/httptrace"
	"time"

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
func (c *WebClient) GetLinks(rawURL string) (statusCode int, links []URLEntity, latency time.Duration, err error) {
	// make sure to use the same http.Client to reuse connections to get links
	// from other pages being served by the same server.
	// Check for robot.txt, maybe?

	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return 0, links, latency, err
	}

	var start time.Time

	trace := &httptrace.ClientTrace{
		GotFirstResponseByte: func() {
			latency = time.Since(start)
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start = time.Now()

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, links, latency, err
	}
	defer resp.Body.Close()

	statusCode = resp.StatusCode

	if statusCode < 200 || statusCode >= 300 {
		return statusCode, links, latency, nil
	}

	links, err = c.parse(rawURL, resp.Body)

	return statusCode, links, latency, err
}

// parse parses the webpage looking for links.
func (c *WebClient) parse(rawURL string, r io.Reader) (links []URLEntity, err error) {
	// Parse <base> tag inside <head> tag if it exists
	// Parse all <a> tags
	// Cater for the fact that a <a> link might be a mailto or a phone or something else.
	// Check if "Opaque" field in URL struct is set
	// Validate whether they are absolute or relative tags. Also check if the relative tags start with a /

	// TODO: Make sure body is utf-8 encoded

	insideHead := false
	baseURL := rawURL

	links = []URLEntity{}

	z := html.NewTokenizer(r)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// EOF
			return links, nil
		case tt == html.StartTagToken || tt == html.SelfClosingTagToken:
			t := z.Token()

			if t.Data == "head" {
				insideHead = true
			}

			// This only works assuming href in <base> is absolute.
			// TODO: confirm this in the HTML spec.
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

			// rawURL = strings.TrimSpace(rawURL)

			// Deals with absolute and relative URLs.
			urlEntity, err := JoinURLs(baseURL, rawURL)
			if err != nil {
				continue
			}

			links = append(links, urlEntity)

		case tt == html.EndTagToken:
			t := z.Token()
			if t.Data == "head" {
				insideHead = false
			}

		}
	}
}

// getHref returns the href attribute from a Token.
func getHref(t html.Token) (ok bool, href string) {
	// Iterate over all of the Token's attributes until it finds an "href"
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	return
}
