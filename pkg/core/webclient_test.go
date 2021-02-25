package core_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gustavooferreira/wcrawler/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebClient(t *testing.T) {
	tests := map[string]struct {
		path               string
		htmlBody           string
		expectedErr        bool
		expectedStatusCode int
		expectedLinks      []core.URLEntity
	}{
		"parse 1": {
			path:               "/",
			htmlBody:           htmlBody1,
			expectedStatusCode: 200,
			expectedLinks: []core.URLEntity{{
				Base: "http://www.example.com",
				Raw:  "http://www.example.com/path1",
			}},
			expectedErr: false,
		},
	}

	// Setup
	c := &http.Client{}
	wc := core.NewWebClient(c)

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, test.htmlBody)
			}))
			defer ts.Close()

			queryURL := ts.URL + test.path

			statusCode, links, err := wc.GetLinks(queryURL)

			if test.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, test.expectedStatusCode, statusCode)
			assert.Equal(t, test.expectedLinks, links)
		})
	}
}
