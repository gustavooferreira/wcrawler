package wcrawler

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gustavooferreira/wcrawler/internal"
)

// ----------
// Interfaces
// ----------

// Connector describes the connector interface
type Connector interface {
	GetLinks(rawURL string) (statusCode int, links []URLEntity, latency time.Duration, err error)
}

type StatsManager interface {
	UpdateStats(updates ...func(*internal.StatsCLIOutput))
	RunOutputFlusher()
}

// ---------------
// Utils functions
// ---------------

// ExtractURL takes any URL and returns a URL string with scheme,authority,path ready
// to be used as a parent URL.
func ExtractURL(rawURL string) (urlEntity URLEntity, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return urlEntity, fmt.Errorf("URL not in a valid format: %s", err)
	}

	if (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return urlEntity, fmt.Errorf("URL provided is not absolute")
	}

	urlEntity.NetLoc = u.Host
	urlEntity.Raw = fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, u.Path)
	if u.RawQuery != "" {
		urlEntity.Raw += "?" + u.RawQuery
	}

	return urlEntity, nil
}

// JoinURLs behaves the same way as parent URL, except that it also includes query params.
// If URL provided is relative, it will join the URLs.
// It will return an error if URL is of an unwanted type, like 'mailto'.
func JoinURLs(baseURL string, rawURL string) (URLEntity, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return URLEntity{}, fmt.Errorf("URL not in a valid format: %s", err)
	}

	// Check if rawURL is one of those types we don't want to handle, i.e., mailto, telephone, etc.
	if u.Opaque != "" || (u.Scheme != "" && u.Scheme != "http" && u.Scheme != "https") {
		return URLEntity{}, fmt.Errorf("URL not in a supported format")
	}

	baseU, err := url.Parse(baseURL)
	if err != nil {
		return URLEntity{}, fmt.Errorf("base URL not in a valid format: %s", err)
	}

	if baseU.Opaque != "" || (baseU.Scheme != "" && baseU.Scheme != "http" && baseU.Scheme != "https") {
		return URLEntity{}, fmt.Errorf("URL not in a supported format")
	}

	mergedU := baseU.ResolveReference(u)
	rawURL = fmt.Sprintf("%s://%s%s", mergedU.Scheme, mergedU.Host, mergedU.Path)

	if mergedU.RawQuery != "" {
		rawURL += "?" + mergedU.RawQuery
	}

	return URLEntity{NetLoc: mergedU.Host, Raw: rawURL}, nil
}
