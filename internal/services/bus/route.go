package bus

import (
	"strings"

	"github.com/queeck/cli/internal/services/commands"
)

const (
	RouteSeparator = `.`
)

func (s *Service) route(command commands.Command) string {
	for route, currentCommand := range s.routes {
		if currentCommand == command {
			return route
		}
	}
	return ""
}

func makeRoute(routes ...string) string {
	return strings.Join(routes, RouteSeparator)
}

func parentRoute(route string) string {
	routes := strings.Split(route, RouteSeparator)
	if len(routes) == 0 || len(routes) == 1 {
		return ""
	}
	parent := strings.Join(routes[:len(routes)-1], RouteSeparator)
	return parent
}

func isParentRoute(parent, child string) bool {
	return strings.HasPrefix(child, parent+RouteSeparator)
}
