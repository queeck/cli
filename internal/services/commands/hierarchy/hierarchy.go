package hierarchy

import (
	"github.com/queeck/cli/internal/services/commands"
	modelRoot "github.com/queeck/cli/internal/services/commands/models"
	modelConfig "github.com/queeck/cli/internal/services/commands/models/config"
	modelConfigGet "github.com/queeck/cli/internal/services/commands/models/config/get"
	modelConfigSet "github.com/queeck/cli/internal/services/commands/models/config/set"
	modelConfigView "github.com/queeck/cli/internal/services/commands/models/config/view"
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

func Default() Hierarchy {
	return Node(modelRoot.New,
		Node(modelConfig.New,
			Node(modelConfigGet.New),
			Node(modelConfigSet.New),
			Node(modelConfigView.New),
		),
	)
}
