package app

import (
	"io"

	"github.com/queeck/cli/internal/services/commands"
)

type bus interface {
	commands.Bus
}

type input interface {
	io.Reader
}

type output interface {
	io.Writer
}
