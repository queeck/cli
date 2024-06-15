package utils

import "strings"

func CRLF(text string) string {
	lines := make([]string, 0, 1)
	for _, line := range strings.Split(text, "\n") {
		lines = append(lines, strings.TrimRight(line, "\r"))
	}
	return strings.Join(lines, "\r\n")
}
