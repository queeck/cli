package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

const (
	delimiter = `.`
)

type Line struct {
	Key   string
	Value string
}

func FindMaxKeyLength(lines []Line) int {
	maxKeyLength := 0
	for _, line := range lines {
		if len(line.Key) > maxKeyLength {
			maxKeyLength = len(line.Key)
		}
	}
	return maxKeyLength
}

func Keys(configData []byte) []string {
	lines := ViewLines(configData)
	results := make([]string, 0, len(lines))
	for _, line := range lines {
		results = append(results, line.Key)
	}
	return results
}

func View(configData []byte) string {
	lines := ViewLines(configData)
	maxKeyLength := FindMaxKeyLength(lines)
	results := make([]string, 0, len(lines))
	for _, line := range lines {
		results = append(results, ViewLine(line, maxKeyLength))
	}
	return strings.Join(results, "\n")
}

func ViewLines(configData []byte) []Line {
	return Walk(gjson.ParseBytes(configData), "")
}

func ViewLine(line Line, maxKeyLength int) string {
	pad := strconv.Itoa(maxKeyLength)
	return fmt.Sprintf("%-"+pad+"s = %v", line.Key, line.Value)
}

func Walk(node gjson.Result, key string) []Line {
	lines := make([]Line, 0, 1)
	if !node.Exists() {
		return append(lines, Line{Key: key, Value: ""})
	}
	if node.IsArray() {
		for n, childNode := range node.Array() {
			children := Walk(childNode, Join(key, strconv.Itoa(n)))
			lines = append(lines, children...)
		}
		return lines
	}
	if node.IsObject() {
		for childKey, childNode := range node.Map() {
			children := Walk(childNode, Join(key, childKey))
			lines = append(lines, children...)
		}
		return lines
	}
	return append(lines, Line{Key: key, Value: node.String()})
}

func Join(key string, chunks ...string) string {
	if key == "" {
		return strings.Join(chunks, delimiter)
	}
	return strings.Join(append([]string{key}, chunks...), delimiter)
}
