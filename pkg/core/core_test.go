package core_test

import (
	"fmt"
	"testing"

	"github.com/gustavooferreira/wcrawler/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsAbsoluteURL(t *testing.T) {
	tests := map[string]struct {
		url      string
		expected bool
	}{
		"url 1": {url: "/path/to/file", expected: false},
		"url 2": {url: "google.com/path/to/file", expected: false},
		"url 3": {url: "http://google.com/path/to/file", expected: true},
		"url 4": {url: "https://google.com/path/to/file", expected: true},
		"url 5": {url: "https://example.com:9999", expected: true},
		"url 6": {url: "https://example.com/", expected: true},
		"url 7": {url: "https://example.com:9999/path/to/file", expected: true},
		"url 8": {url: "https://example.com:9999/path/to/file#fragment", expected: true},
		"url 9": {url: "https://example.com:9999/path/to/file?q=1&w=2#fragment", expected: true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			value := core.IsAbsoluteURL(test.url)
			assert.Equal(t, test.expected, value)
		})
	}
}

func TestExtractParentURL(t *testing.T) {
	tests := map[string]struct {
		url         string
		expectedURL string
		expectedErr bool
	}{
		"url 0": {url: "", expectedURL: "", expectedErr: true},
		"url 1": {url: "qwe:\n//1.1.1.1/path_to_file?qwe=213", expectedURL: "", expectedErr: true},
		"url 2": {url: "/path/to/file", expectedURL: "", expectedErr: true},
		"url 3": {url: "google.com/path/to/file", expectedURL: "", expectedErr: true},
		"url 4": {url: "http://google.com/path/to/file", expectedURL: "http://google.com/path/to/file", expectedErr: false},
		"url 5": {url: "https://google.com", expectedURL: "https://google.com", expectedErr: false},
		"url 6": {url: "https://example.com:9999", expectedURL: "https://example.com:9999", expectedErr: false},
		"url 7": {url: "https://example.com/", expectedURL: "https://example.com/", expectedErr: false},
		// "url 7": {url: "https://example.com:9999/path/to/file", expected: true},
		// "url 8": {url: "https://example.com:9999/path/to/file#fragment", expected: true},
		// "url 9": {url: "https://example.com:9999/path/to/file?q=1&w=2#fragment", expected: true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var errBool bool
			var errMsg string
			value, err := core.ExtractParentURL(test.url)
			if err != nil {
				errBool = true
				errMsg = fmt.Sprintf(" - err: %s", err.Error())
			}
			require.Equal(t, test.expectedErr, errBool, "error field"+errMsg)

			assert.Equal(t, test.expectedURL, value)
		})
	}
}
