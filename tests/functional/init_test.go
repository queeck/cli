package functional

import (
	"time"

	"github.com/queeck/cli/internal/pkg/cli"
	serviceApp "github.com/queeck/cli/internal/services/app"
	serviceBus "github.com/queeck/cli/internal/services/bus"
	hiera "github.com/queeck/cli/internal/services/commands/hierarchy"
	"github.com/queeck/cli/internal/services/config"
	serviceState "github.com/queeck/cli/internal/services/state"
	serviceTemplates "github.com/queeck/cli/internal/services/templates"
)

const (
	defaultFPS     = 120
	defaultTimeout = 10 * time.Millisecond
)

func makeApp(config config.Config, args []string, input input, output output) app {
	hierarchy := hiera.Default()

	arguments := cli.Parse(args)

	state := serviceState.New()

	templates := serviceTemplates.New()

	commandBus := serviceBus.New(hierarchy, arguments, state, config, templates)

	return serviceApp.New(commandBus, input, output, serviceApp.WithFPS(defaultFPS))
}
