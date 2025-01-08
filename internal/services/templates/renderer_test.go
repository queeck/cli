package templates

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRendering_Render(t *testing.T) {
	cases := map[string]struct {
		template  string
		keyValues []string
		expected  string
	}{
		"empty": {
			template:  "",
			keyValues: nil,
			expected:  "",
		},
		"key_values_redundant": {
			template:  "",
			keyValues: []string{"key", "value"},
			expected:  "",
		},
		"missing_key": {
			template:  "key is {{ .key }}",
			keyValues: nil,
			expected:  "key is <no value>",
		},
		"key_value": {
			template:  "key is {{ .key }}",
			keyValues: []string{"key", "value"},
			expected:  "key is value",
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			actual := New().Render(test.template, test.keyValues...)
			require.Equal(t, test.expected, actual)
		})
	}
}

func TestKeyValuesToData(t *testing.T) {
	cases := map[string]struct {
		keyValues []string
		expected  map[string]string
	}{
		"empty": {
			//
		},
		"key_only": {
			keyValues: []string{"key"},
			expected:  map[string]string{"key": ""},
		},
		"key_value": {
			keyValues: []string{"key", "value"},
			expected:  map[string]string{"key": "value"},
		},
		"key_value_key1": {
			keyValues: []string{
				"key", "value",
				"key1",
			},
			expected: map[string]string{
				"key":  "value",
				"key1": "",
			},
		},
		"key_value_key1_value1": {
			keyValues: []string{
				"key", "value",
				"key1", "value1",
			},
			expected: map[string]string{
				"key":  "value",
				"key1": "value1",
			},
		},
		"key_value_key1_value1_key2": {
			keyValues: []string{
				"key", "value",
				"key1", "value1",
				"key2",
			},
			expected: map[string]string{
				"key":  "value",
				"key1": "value1",
				"key2": "",
			},
		},
		"key_value_key1_value1_key2_value2": {
			keyValues: []string{
				"key", "value",
				"key1", "value1",
				"key2", "value2",
			},
			expected: map[string]string{
				"key":  "value",
				"key1": "value1",
				"key2": "value2",
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			actual := keyValuesToData(test.keyValues...)
			require.Equal(t, test.expected, actual)
		})
	}
}
