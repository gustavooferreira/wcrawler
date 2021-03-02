package core_test

import (
	"testing"

	"github.com/gustavooferreira/wcrawler/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractParentURL(t *testing.T) {
	tests := map[string]struct {
		url            string
		expectedURLEnt core.URLEntity
		expectedErr    bool
	}{
		"url 1": {url: "",
			expectedURLEnt: core.URLEntity{NetLoc: "", Raw: ""},
			expectedErr:    true},
		"url 2": {url: "qwe:\n//1.1.1.1/path_to_file?qwe=213",
			expectedURLEnt: core.URLEntity{NetLoc: "", Raw: ""},
			expectedErr:    true},
		"url 3": {url: "/path/to/file",
			expectedURLEnt: core.URLEntity{NetLoc: "", Raw: ""},
			expectedErr:    true},
		"url 4": {url: "google.com/path/to/file",
			expectedURLEnt: core.URLEntity{NetLoc: "", Raw: ""},
			expectedErr:    true},
		"url 5": {url: "http://google.com/path/to/file",
			expectedURLEnt: core.URLEntity{NetLoc: "google.com", Raw: "http://google.com/path/to/file"},
			expectedErr:    false},
		"url 6": {url: "https://google.com",
			expectedURLEnt: core.URLEntity{NetLoc: "google.com", Raw: "https://google.com"},
			expectedErr:    false},
		"url 7": {url: "https://example.com:9999",
			expectedURLEnt: core.URLEntity{NetLoc: "example.com:9999", Raw: "https://example.com:9999"},
			expectedErr:    false},
		"url 8": {url: "https://example.com/",
			expectedURLEnt: core.URLEntity{NetLoc: "example.com", Raw: "https://example.com/"},
			expectedErr:    false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			value, err := core.ExtractURL(test.url)
			if test.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expectedURLEnt, value)
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
			expectedURL: core.URLEntity{NetLoc: "test123.com", Raw: "https://test123.com/path/to/file"},
			expectedErr: false,
		},
		"url 3": {
			parentURL:   "http://example.com",
			url:         "/path/to/file",
			expectedURL: core.URLEntity{NetLoc: "example.com", Raw: "http://example.com/path/to/file"},
			expectedErr: false,
		},
		"url 4": {
			parentURL:   "http://example.com/base/",
			url:         "path/to/file",
			expectedURL: core.URLEntity{NetLoc: "example.com", Raw: "http://example.com/base/path/to/file"},
			expectedErr: false,
		},
		"url 5": {
			parentURL:   "http://example.com/base/index.html",
			url:         "../path/to/file",
			expectedURL: core.URLEntity{NetLoc: "example.com", Raw: "http://example.com/path/to/file"},
			expectedErr: false,
		},
		"url 6": {
			parentURL:   "http://example.com/base/path/to/index.html",
			url:         "/path/to/file",
			expectedURL: core.URLEntity{NetLoc: "example.com", Raw: "http://example.com/path/to/file"},
			expectedErr: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			value, err := core.JoinURLs(test.parentURL, test.url)
			if test.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expectedURL, value)
			}
		})
	}
}
