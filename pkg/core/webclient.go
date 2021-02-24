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
func (c *WebClient) GetLinks(rawURL string) (statusCode int, links []URLEntity, err error) {
	result := []URLEntity{}

	resp, err := c.client.Get(rawURL)
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

		// TODO: take care of relative links here

		l := URLEntity{Host: u.Host, String: link}
		result = append(result, l)
	}

	// body, err := io.ReadAll(resp.Body)
	// fmt.Println(string(body))

	return statusCode, result, nil
}

// Parse <base> tag if it exists
// Parse all <a> tags
// Cater for the fact that a <a> link might be a mailto or a phone or something else.
// Check if "Opaque" field in URL struct is set
// Validate whether they are absolute or relative tags. Also check if the relative tags start with a /
func parse(r io.Reader) ([]string, error) {

	// TODO: same webpages will have the <base> tag inside <head> which should be used
	// when joining relative URLs.

	insideHead := false
	base := ""

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

			if t.Data == "head" {
				insideHead = true
			}

			if t.Data == "base" && insideHead == true {
				ok, url := getHref(t)
				if ok {
					base = url
				}
			}

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

			// Check type of link
			// if relative, join it to base
			_ = base

			url = strings.TrimSpace(url)

			result = append(result, url)

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
