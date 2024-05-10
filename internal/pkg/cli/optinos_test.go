package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithFlags(t *testing.T) {
	cases := map[string]struct {
		flags    []string
		expected *Settings
	}{
		"empty": {
			flags:    nil,
			expected: &Settings{},
		},
		"verbose": {
			flags: []string{"verbose"},
			expected: &Settings{
				flags: map[string]bool{"verbose": true},
			},
		},
		"v": {
			flags: []string{"v"},
			expected: &Settings{
				flags:         map[string]bool{"v": true},
				shortFlagsMap: map[string]bool{"v": true},
			},
		},
		"shorts": {
			flags: []string{"v", "verbose", "h", "help", "s", "simple"},
			expected: &Settings{
				flags: map[string]bool{
					"v":       true,
					"verbose": true,
					"h":       true,
					"help":    true,
					"s":       true,
					"simple":  true,
				},
				shortFlagsMap: map[string]bool{
					"h": true,
					"s": true,
					"v": true,
				},
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			actual := &Settings{}
			WithFlags(test.flags...)(actual)
			require.Equal(t, test.expected, actual)
		})
	}
}
