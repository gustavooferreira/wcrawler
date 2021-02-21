package core

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type HTMLParser struct {
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
