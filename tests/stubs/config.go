package stubs

import (
	"slices"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	servicesConfig "github.com/queeck/cli/internal/services/config"
)

func ConfigInMemory(source string) servicesConfig.Config {
	data := []byte(source)

	get := func(path string) gjson.Result {
		return gjson.ParseBytes(data).Get(path)
	}

	return &configMock{
		GetFunc: func(path string) (any, bool) {
			result := get(path)
			if !result.Exists() {
				return nil, false
			}
			return result.Value(), true
		},
		GetStringFunc: func(path string) (string, bool) {
			result := get(path)
			if !result.Exists() {
				return "", false
			}
			return result.String(), true
		},
		KeysFunc: func() []string {
			list := servicesConfig.Keys(data)
			slices.Sort(list)
			return list
		},
		PathFunc: func() string {
			return "tests/stubs/config"
		},
		SaveFunc: func() error {
			return nil
		},
		SetFunc: func(path string, value any) (err error) {
			if data, err = sjson.SetBytes(data, path, value); err != nil {
				return err
			}
			return nil
		},
		TypeFunc: func(path string) servicesConfig.ValueType {
			result := get(path)
			if !result.Exists() {
				return servicesConfig.TypeNull
			}
			if result.IsBool() {
				return servicesConfig.TypeBool
			}
			if result.IsArray() {
				return servicesConfig.TypeArray
			}
			if result.IsObject() {
				return servicesConfig.TypeObject
			}
			if result.Type == gjson.String {
				return servicesConfig.TypeString
			}
			if result.Type == gjson.Number {
				return servicesConfig.TypeNumber
			}
			return servicesConfig.TypeNull
		},
		ViewFunc: func() string {
			return servicesConfig.View(data)
		},
	}
}
