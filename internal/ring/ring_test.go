package ring_test

import (
	"testing"

	"github.com/gustavooferreira/wcrawler/internal/ring"
	"github.com/stretchr/testify/assert"
)

func TestRingBuffer(t *testing.T) {
	tests := map[string]struct {
		size           int
		entries        []string
		expectedOutput []string
	}{
		"test 1": {
			size:           3,
			entries:        []string{"a", "b", "c"},
			expectedOutput: []string{"a", "b", "c"},
		},
		"test 2": {
			size:           2,
			entries:        []string{"a", "b", "c"},
			expectedOutput: []string{"b", "c"},
		},
		"test 3": {
			size:           5,
			entries:        []string{},
			expectedOutput: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			rb := ring.New(test.size)

			for _, elem := range test.entries {
				rb.Add(elem)
			}

			value := rb.ReadAll()

			assert.Equal(t, test.expectedOutput, value)
		})
	}
}
