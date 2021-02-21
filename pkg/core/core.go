package core

import (
	"fmt"
	"net/url"
)

// ----------
// Interfaces
// ----------

// Connector describes the connector interface
type Connector interface {
	GetLinks(baseURL string) (statusCode int, links []string, err error)
}

// ---------------
// Utils functions
// ---------------

// IsAbsoluteURL validates URL
func IsAbsoluteURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}

// ExtractParentURL takes any URL and returns a URL string with scheme,authority,path ready
// to be used as a parent URL.
func ExtractParentURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("URL not in a valid format: %s", err)
	}

	if (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return "", fmt.Errorf("URL provided is not absolute")
	}

	rawURL = fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, u.Path)
	return rawURL, nil
}

// ExtractURL behaves the same way as parent URL, except that in also includes query params.
// If URL provided is relative, it will join the URLs.
func ExtractURL(parentURL string, rawURL string) (URLEntity, error) {
	// parentURL must be well formatted, including scheme, authority
	// and potentially a path, but nothing else.

	// Make sure we only take into consideration the scheme, authority, path and query parts of the URL

	// Validate parentURL

	// Validate rawURL (if URL is relative, join with parent URL)

	// u, err := url.Parse(rawURL)
	// if err != nil {
	// 	return false, err
	// }

	// if (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
	// 	return false, fmt.Errorf("URL provided is not absolute")
	// }

	// rawURL = u.Scheme
	return URLEntity{}, nil
}
