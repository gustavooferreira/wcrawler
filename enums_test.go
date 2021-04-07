package wcrawler_test

import (
	"testing"

	"github.com/gustavooferreira/wcrawler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReturningAppStateString(t *testing.T) {
	tests := map[string]struct {
		input          wcrawler.AppState
		expectedOutput string
	}{
		"test 'IDLE' state": {
			input:          wcrawler.AppState_IDLE,
			expectedOutput: "IDLE",
		},
		"test 'Running' state": {
			input:          wcrawler.AppState_Running,
			expectedOutput: "Running",
		},
		"test 'Finished' state": {
			input:          wcrawler.AppState_Finished,
			expectedOutput: "Finished",
		},
		"test 'Unknown' state": {
			input:          wcrawler.AppState_Unknown,
			expectedOutput: "Unknown",
		},
		"test missing state": {
			input:          -100,
			expectedOutput: "Unknown",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			value := test.input.String()
			assert.Equal(t, test.expectedOutput, value)
		})
	}
}

func TestParsingAppStateString(t *testing.T) {
	tests := map[string]struct {
		input          string
		expectedErr    bool
		expectedOutput wcrawler.AppState
	}{
		"test 'IDLE' state": {
			input:          "IDLE",
			expectedOutput: wcrawler.AppState_IDLE,
		},
		"test 'Running' state": {
			input:          "Running",
			expectedOutput: wcrawler.AppState_Running,
		},
		"test 'Finished' state": {
			input:          "Finished",
			expectedOutput: wcrawler.AppState_Finished,
		},
		"test 'Unknown' state": {
			expectedOutput: wcrawler.AppState_Unknown,
			input:          "Unknown",
		},
		"test missing state": {
			input:       "qwueyqwie",
			expectedErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var value wcrawler.AppState
			err := value.Parse(test.input)

			if test.expectedErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.expectedOutput, value)
		})
	}
}
