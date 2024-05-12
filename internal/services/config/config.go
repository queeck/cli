package config

import (
	"os"
	"path"
	"slices"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	serviceDirectory "github.com/queeck/cli/internal/services/directory"
)

const (
	defaultFilename = "config.json"
	defaultData     = `{}`
)

type Config interface {
	Get(path string) (value any, has bool)
	GetString(path string) (value string, has bool)
	Set(path string, value any) error
	View() string
	Keys() []string
	Path() string
}

type config struct {
	directory directory
	data      []byte
	path      string
}

func Reload(directory directory) (Config, error) {
	filePath := path.Join(directory.Local(), defaultFilename)
	data, err := os.ReadFile(filePath)
	if err != nil && os.IsNotExist(err) {
		err = os.WriteFile(filePath, []byte(defaultData), serviceDirectory.Permissions)
	}
	if err != nil {
		return nil, err
	}
	cfg := &config{
		directory: directory,
		data:      data,
		path:      filePath,
	}
	return cfg, nil
}

func (c *config) Save() error {
	filePath := path.Join(c.directory.Local(), defaultFilename)
	return os.WriteFile(filePath, c.data, serviceDirectory.Permissions)
}

func (c *config) get(path string) gjson.Result {
	return gjson.ParseBytes(c.data).Get(path)
}

func (c *config) Get(path string) (any, bool) {
	result := c.get(path)
	if !result.Exists() {
		return nil, false
	}
	return result.Value(), true
}

func (c *config) GetString(path string) (string, bool) {
	result := c.get(path)
	if !result.Exists() {
		return "", false
	}
	return result.String(), true
}

func (c *config) Set(path string, value any) error {
	data, err := sjson.SetBytes(c.data, path, value)
	if err != nil {
		return err
	}
	c.data = data
	return nil
}

func (c *config) Path() string {
	return c.path
}

func (c *config) Keys() []string {
	list := keys(c.data)
	slices.Sort(list)
	return list
}

func (c *config) View() string {
	return view(c.data)
}
