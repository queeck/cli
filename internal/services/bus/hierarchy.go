package bus

import (
	"github.com/queeck/cli/internal/services/commands"
	"github.com/queeck/cli/internal/services/commands/hierarchy"
)

func (s *Service) walk(hierarchy hierarchy.Hierarchy, parents ...string) map[string]commands.Command {
	if hierarchy == nil {
		return nil
	}
	builder := hierarchy.CommandBuilder()
	command := builder(s)
	code := command.Code()
	route := makeRoute(append(parents, code)...)
	routes := map[string]commands.Command{
		route: command,
	}
	if len(hierarchy.Children()) != 0 {
		for _, child := range hierarchy.Children() {
			for route, builder := range s.walk(child, append(parents, code)...) {
				routes[route] = builder
			}
		}
	}
	return routes
}
