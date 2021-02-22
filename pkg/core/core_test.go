package core_test

import (
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
		"url 1": {url: "", expectedURL: "", expectedErr: true},
		"url 2": {url: "qwe:\n//1.1.1.1/path_to_file?qwe=213", expectedURL: "", expectedErr: true},
		"url 3": {url: "/path/to/file", expectedURL: "", expectedErr: true},
		"url 4": {url: "google.com/path/to/file", expectedURL: "", expectedErr: true},
		"url 5": {url: "http://google.com/path/to/file", expectedURL: "http://google.com/path/to/file", expectedErr: false},
		"url 6": {url: "https://google.com", expectedURL: "https://google.com", expectedErr: false},
		"url 7": {url: "https://example.com:9999", expectedURL: "https://example.com:9999", expectedErr: false},
		"url 8": {url: "https://example.com/", expectedURL: "https://example.com/", expectedErr: false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			value, err := core.ExtractParentURL(test.url)
			if test.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expectedURL, value)
			}
		})
	}
}

func TestExtractURL(t *testing.T) {
	tests := map[string]struct {
		parentURL   string
		url         string
		expectedURL core.URLEntity
		expectedErr bool
	}{
		"url 1": {parentURL: "http://example.com", url: "qwe:\nqwe", expectedURL: core.URLEntity{}, expectedErr: true},
		"url 2": {
			parentURL:   "http://example.com",
			url:         "https://test123.com/path/to/file",
			expectedURL: core.URLEntity{Host: "test123.com", String: "https://test123.com/path/to/file"},
			expectedErr: false,
		},
		"url 3": {
			parentURL:   "http://example.com",
			url:         "/path/to/file",
			expectedURL: core.URLEntity{Host: "example.com", String: "http://example.com/path/to/file"},
			expectedErr: false,
		},
		"url 4": {
			parentURL:   "http://example.com/base",
			url:         "path/to/file",
			expectedURL: core.URLEntity{Host: "example.com", String: "http://example.com/base/path/to/file"},
			expectedErr: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			value, err := core.ExtractURL(test.parentURL, test.url)
			if test.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expectedURL, value)
			}
		})
	}
}
