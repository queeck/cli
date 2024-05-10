package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	cases := map[string]struct {
		value    any
		expected string
	}{
		"empty_bytes": {
			value:    []byte{},
			expected: `""`,
		},
		"empty_numbers": {
			value:    []int64{},
			expected: `[]`,
		},
		"empty_zero": {
			value:    0,
			expected: `0`,
		},
		"empty_string": {
			value:    "",
			expected: `""`,
		},
		"empty_nil": {
			value:    nil,
			expected: `null`,
		},
		"simple_object": {
			value:    map[string]any{"json": "awesome"},
			expected: `{"json":"awesome"}`,
		},
		"simple_array": {
			value:    []any{123, "test", -5, 1e10, 1e50},
			expected: `[123,"test",-5,10000000000,1e+50]`,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			actual := Encode(test.value)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestDecodeOnMap(t *testing.T) {
	type SimpleObject struct {
		Value string `json:"json"`
	}

	assert.Equal(t,
		&map[string]any{"json": "awesome"},
		Decode[map[string]any]([]byte(`{"json":"awesome"}`)))

	assert.Equal(t,
		&SimpleObject{Value: "awesome"},
		Decode[SimpleObject]([]byte(`{"json":"awesome"}`)))
}
