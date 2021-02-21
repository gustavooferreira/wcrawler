package core_test

import (
	"testing"

	"github.com/gustavooferreira/wcrawler/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestIsURL(t *testing.T) {
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
			value := core.IsURL(test.url)
			assert.Equal(t, test.expected, value)
		})
	}
}
