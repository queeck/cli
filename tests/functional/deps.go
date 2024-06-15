package functional

import (
	"context"
	"io"

	"github.com/queeck/cli/internal/services/commands"
)

//go:generate moq -skip-ensure -out deps_moq_test.go . input output

type app interface {
	Run(ctx context.Context) error
	CommandBus() commands.Bus
}

type input interface {
	io.Reader
}

type output interface {
	io.Writer
}
