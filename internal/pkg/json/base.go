package json

import (
	jsoniter "github.com/json-iterator/go"
)

func Encode(v any) string {
	data, err := jsoniter.ConfigFastest.Marshal(v)
	if err != nil {
		return ""
	}
	return string(data)
}

func EncodePretty(v any) string {
	data, err := jsoniter.ConfigFastest.MarshalIndent(v, "", "  ")
	if err != nil {
		return ""
	}
	return string(data)
}

func Decode[T any](data []byte) *T {
	value := new(T)
	if err := jsoniter.ConfigFastest.Unmarshal(data, value); err != nil {
		return nil
	}
	return value
}
