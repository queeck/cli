package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	cases := map[string]struct {
		args     []string
		options  []Option
		expected *Args
	}{
		"empty": {
			args: nil,
			expected: &Args{
				commands:            []string{},
				flagsMap:            map[string]bool{},
				options:             map[string]string{},
				positionalArguments: []string{},
			},
		},
		"config": {
			args: []string{"config"},
			expected: &Args{
				commands:            []string{"config"},
				flagsMap:            map[string]bool{},
				options:             map[string]string{},
				positionalArguments: []string{},
			},
		},
		"config_simple": {
			args: []string{"config", "view", "--simple"},
			expected: &Args{
				commands:            []string{"config", "view"},
				flagsMap:            map[string]bool{"simple": true},
				options:             map[string]string{},
				positionalArguments: []string{},
			},
		},
		"queue_publish_complex": {
			args: []string{
				"queue",
				"publish",
				"--dry",
				"-c",
				"test",
				"-d=50",
				"--scheme",
				"test.entity.update",
				"--data",
				"{\"uid\":\"face\"}",
				"--pretty",
			},
			expected: &Args{
				commands: []string{"queue", "publish"},
				flagsMap: map[string]bool{
					"dry":    true,
					"pretty": true,
				},
				options: map[string]string{
					"c":      "test",
					"d":      "50",
					"scheme": "test.entity.update",
					"data":   "{\"uid\":\"face\"}",
				},
				positionalArguments: []string{},
			},
		},
		"queue_shorthand": {
			args: []string{
				"queue",
				"publish",
				"-fpD",
			},
			options: []Option{WithFlags("c", "d", "p", "f", "D")},
			expected: &Args{
				commands: []string{"queue", "publish"},
				flagsMap: map[string]bool{
					"f": true,
					"p": true,
					"D": true,
				},
				options:             map[string]string{},
				positionalArguments: []string{},
			},
		},
		"curl_with_data": {
			args: []string{
				"curl",
				"-X",
				"POST",
				"https://www.yourwebsite.com/login/",
				"-d",
				"username=username&password=password",
			},
			expected: &Args{
				commands: []string{"curl"},
				flagsMap: map[string]bool{},
				options: map[string]string{
					"X": "POST",
					"d": "username=username&password=password",
				},
				positionalArguments: []string{"https://www.yourwebsite.com/login/"},
			},
		},
		"curl_with_cookie": {
			args: []string{
				"curl",
				"-b",
				"session=abc123; user=JohnDoe",
				"https://example.com",
			},
			expected: &Args{
				commands: []string{"curl"},
				flagsMap: map[string]bool{},
				options: map[string]string{
					"b": "session=abc123; user=JohnDoe",
				},
				positionalArguments: []string{"https://example.com"},
			},
		},
		"git_without_positional": {
			args: []string{"git", "pull", "origin", "master"},
			expected: &Args{
				commands:            []string{"git", "pull", "origin", "master"},
				flagsMap:            map[string]bool{},
				options:             map[string]string{},
				positionalArguments: []string{},
			},
		},
		"git_with_positional": {
			args:    []string{"git", "pull", "origin", "master"},
			options: []Option{WithCommandsCount(2)},
			expected: &Args{
				commands:            []string{"git", "pull"},
				flagsMap:            map[string]bool{},
				options:             map[string]string{},
				positionalArguments: []string{"origin", "master"},
			},
		},
		"git_without_flag": {
			args: []string{"git", "pull", "-v", "origin", "master"},
			expected: &Args{
				commands: []string{"git", "pull"},
				flagsMap: map[string]bool{},
				options: map[string]string{
					"v": "origin",
				},
				positionalArguments: []string{"master"},
			},
		},
		"git_with_flag": {
			args:    []string{"git", "pull", "-v", "origin", "master"},
			options: []Option{WithFlags("v", "verbose")},
			expected: &Args{
				commands: []string{"git", "pull"},
				flagsMap: map[string]bool{
					"v": true,
				},
				options:             map[string]string{},
				positionalArguments: []string{"origin", "master"},
			},
		},
		"git_remote_add": {
			args:    []string{"git", "remote", "add", "origin", "https://devhub.example/queeck/cli.git"},
			options: []Option{WithCommandsCount(3)},
			expected: &Args{
				commands:            []string{"git", "remote", "add"},
				flagsMap:            map[string]bool{},
				options:             map[string]string{},
				positionalArguments: []string{"origin", "https://devhub.example/queeck/cli.git"},
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			actual := Parse(test.args, test.options...)
			require.Equal(t, test.expected, actual)
		})
	}
}
