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
	// Check whether rawURL is absolute or not
	u, err := url.Parse(rawURL)
	if err != nil {
		return URLEntity{}, fmt.Errorf("URL not in a valid format: %s", err)
	}

	if (u.Scheme == "http" || u.Scheme == "https") && u.Host != "" {
		rawURL = fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, u.Path)

		if u.RawQuery != "" {
			rawURL += "?" + u.RawQuery
		}

		return URLEntity{Host: u.Host, String: rawURL}, nil
	}

	parentU, err := url.Parse(parentURL)
	if err != nil {
		return URLEntity{}, fmt.Errorf("parent URL not in a valid format: %s", err)
	}

	if strings.HasPrefix(u.Path, "/") {
		rawURL = fmt.Sprintf("%s://%s%s", parentU.Scheme, parentU.Host, u.Path)

		if u.RawQuery != "" {
			rawURL += "?" + u.RawQuery
		}

		return URLEntity{Host: parentU.Host, String: rawURL}, nil
	} else {
		u.Path = path.Join(parentU.Path, u.Path)
		rawURL = fmt.Sprintf("%s://%s%s", parentU.Scheme, parentU.Host, u.Path)

		if u.RawQuery != "" {
			rawURL += "?" + u.RawQuery
		}

		return URLEntity{Host: parentU.Host, String: rawURL}, nil
	}
}
