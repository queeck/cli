package config

import (
	"os"
	"path"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	serviceDirectory "github.com/queeck/cli/internal/services/directory"
)

const (
	defaultFilename = "config.json"
	defaultData     = `{}`
)

type Config interface {
	Get(path string) any
	GetString(path string) string
	Set(path string, value any) error
	View() string
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

func (c *config) Get(path string) any {
	return gjson.ParseBytes(c.data).Get(path).Value()
}

func (c *config) GetString(path string) string {
	return gjson.ParseBytes(c.data).Get(path).String()
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

func (c *config) View() string {
	return View(c.data)
}
