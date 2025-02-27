package bus

import (
	"github.com/queeck/cli/internal/pkg/cli"
	"github.com/queeck/cli/internal/services/commands"
	"github.com/queeck/cli/internal/services/commands/hierarchy"
	modelRoot "github.com/queeck/cli/internal/services/commands/models"
	"github.com/queeck/cli/internal/services/commands/selected"
	configService "github.com/queeck/cli/internal/services/config"
	serviceState "github.com/queeck/cli/internal/services/state"
	serviceTemplates "github.com/queeck/cli/internal/services/templates"
)

var _ commands.Bus = &Service{}

type Service struct {
	arguments arguments
	state     state
	templates templates
	config    config
	routes    map[string]commands.Command
	window    *window
}

func New(
	hierarchy hierarchy.Hierarchy,
	arguments arguments,
	state state,
	config config,
	templates templates,
) *Service {
	bus := &Service{
		arguments: arguments,
		state:     state,
		templates: templates,
		config:    config,
	}

	bus.routes = bus.walk(hierarchy)

	return bus
}

func (s *Service) State() serviceState.State {
	return s.state
}

func (s *Service) Arguments() cli.Arguments {
	return s.arguments
}

func (s *Service) Templates() serviceTemplates.Renderer {
	return s.templates
}

func (s *Service) Config() configService.Config {
	return s.config
}

func (s *Service) Parent(command commands.Command) commands.Command {
	return s.routes[parentRoute(s.route(command))]
}

func (s *Service) Children(command commands.Command) []commands.Command {
	children := make([]commands.Command, 0)
	route := s.route(command)
	for currentRoute, currentCommand := range s.routes {
		if isParentRoute(route, currentRoute) {
			children = append(children, currentCommand)
		}
	}
	return children
}

func (s *Service) Child(command commands.Command, code string) commands.Command {
	return s.routes[makeRoute(s.route(command), code)]
}

func (s *Service) CommandByCLIArguments(arguments cli.Arguments) commands.Command {
	routes := append([]string{modelRoot.Code}, arguments.Commands()...)
	for {
		if command, has := s.routes[makeRoute(routes...)]; has {
			return command
		}
		routes = routes[0 : len(routes)-1]
		if len(routes) == 0 {
			break
		}
	}
	return nil
}

func (s *Service) SelectedCommands(command commands.Command, code string) string {
	return selected.New(s).Render(command, code)
}

func (s *Service) CommandRoot() commands.Command {
	return s.routes[makeRoute(modelRoot.Code)]
}

func (s *Service) ForEach(fn func(route string, command commands.Command) bool) {
	for route, cmd := range s.routes {
		if !fn(route, cmd) {
			return
		}
	}
}

func (s *Service) Window() commands.Window {
	return s.window
}

func (s *Service) UpdateWindowSize(width, height int) {
	// used for viewport
	s.window = &window{
		height: height,
		width:  width,
	}
}
