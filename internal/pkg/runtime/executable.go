package runtime

import (
	"os"
	"path/filepath"
)

func Executable() string {
	ex, err := os.Executable()
	if err != nil {
		return "cli"
	}
	_, file := filepath.Split(ex)
	return file
}
