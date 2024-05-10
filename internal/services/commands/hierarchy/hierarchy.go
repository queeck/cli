package hierarchy

import (
	"github.com/queeck/cli/internal/services/commands"
)

type Hierarchy interface {
	CommandBuilder() CommandBuilder
	Children() []Hierarchy
}

type CommandBuilder func(commands.Bus) commands.Command

type hiera struct {
	builder  CommandBuilder
	children []Hierarchy
}

func (h *hiera) CommandBuilder() CommandBuilder {
	return h.builder
}

func (h *hiera) Children() []Hierarchy {
	return h.children
}

func Node(builder CommandBuilder, children ...Hierarchy) Hierarchy {
	return &hiera{builder: builder, children: children}
}
