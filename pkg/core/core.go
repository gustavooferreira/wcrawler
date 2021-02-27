package core

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

// ----------
// Interfaces
// ----------

// Connector describes the connector interface
type Connector interface {
	GetLinks(rawURL string) (statusCode int, links []URLEntity, err error)
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
func ExtractParentURL(rawURL string) (baseURL string, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("URL not in a valid format: %s", err)
	}

	if (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return "", fmt.Errorf("URL provided is not absolute")
	}

	baseURL = fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, u.Path)
	return baseURL, nil
}

// ExtractURL behaves the same way as parent URL, except that it also includes query params.
// If URL provided is relative, it will join the URLs.
// It will return an error if URL is of an unwanted type, like 'mailto'.
func ExtractURL(baseURL string, rawURL string) (URLEntity, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return URLEntity{}, fmt.Errorf("URL not in a valid format: %s", err)
	}

	// Check if rawURL is of one of those types we don't want to handle, i.e., mailto, telephone, etc.
	if u.Opaque != "" {
		return URLEntity{}, fmt.Errorf("URL not in a processable format")
	}

	// Check whether rawURL is absolute (if it is, get ride of the fragment part)
	if (u.Scheme == "http" || u.Scheme == "https") && u.Host != "" {
		rawURL = fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, u.Path)

		if u.RawQuery != "" {
			rawURL += "?" + u.RawQuery
		}

		return URLEntity{Host: u.Host, Raw: rawURL}, nil
	}

	// If we reached here, it's because URL is in relative format
	baseU, err := url.Parse(baseURL)
	if err != nil {
		return URLEntity{}, fmt.Errorf("base URL not in a valid format: %s", err)
	}

	// Relative URL relative to host
	if strings.HasPrefix(u.Path, "/") {
		rawURL = fmt.Sprintf("%s://%s%s", baseU.Scheme, baseU.Host, u.Path)

		if u.RawQuery != "" {
			rawURL += "?" + u.RawQuery
		}

		return URLEntity{Host: baseU.Host, Raw: rawURL}, nil
	}

	// Relative URL relative to current document
	u.Path = path.Join(baseU.Path, u.Path)
	rawURL = fmt.Sprintf("%s://%s%s", baseU.Scheme, baseU.Host, u.Path)

	if u.RawQuery != "" {
		rawURL += "?" + u.RawQuery
	}

	return URLEntity{Host: baseU.Host, Raw: rawURL}, nil
}
