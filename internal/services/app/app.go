package app

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/queeck/cli/internal/pkg/cli"
	serviceBus "github.com/queeck/cli/internal/services/bus"
	"github.com/queeck/cli/internal/services/commands"
	hiera "github.com/queeck/cli/internal/services/commands/hierarchy"
	serviceConfig "github.com/queeck/cli/internal/services/config"
	serviceDirectory "github.com/queeck/cli/internal/services/directory"
	serviceState "github.com/queeck/cli/internal/services/state"
	serviceTemplates "github.com/queeck/cli/internal/services/templates"
)

var (
	defaultArgs   = os.Args[1:]
	defaultInput  = os.Stdin
	defaultOutput = os.Stdout
)

type App struct {
	input   input
	output  output
	bus     bus
	options *appOptions
}

func Default() (*App, error) {
	directory, err := serviceDirectory.New()
	if err != nil {
		return nil, fmt.Errorf("failed to make directory service: %w", err)
	}

	if !directory.IsLocalCreated() {
		if err = directory.CreateLocal(); err != nil {
			return nil, fmt.Errorf("failed to create local directory: %w", err)
		}
	}

	config, err := serviceConfig.Reload(directory)
	if err != nil {
		return nil, fmt.Errorf("failed to reload config service: %w", err)
	}

	hierarchy := hiera.Default()

	arguments := cli.Parse(defaultArgs)

	state := serviceState.New()

	templates := serviceTemplates.New()

	commandBus := serviceBus.New(hierarchy, arguments, state, config, templates)

	app := New(commandBus, defaultInput, defaultOutput)

	return app, nil
}

func New(bus bus, input input, output output, options ...Option) *App {
	opts := &appOptions{}
	for _, option := range options {
		option(opts)
	}
	return &App{
		bus:     bus,
		input:   input,
		output:  output,
		options: opts,
	}
}

func (a *App) Run(ctx context.Context) error {
	options := []tea.ProgramOption{
		tea.WithMouseCellMotion(),
		tea.WithInput(a.input),
		tea.WithOutput(a.output),
		tea.WithContext(ctx),
	}
	if a.options.fps != 0 {
		options = append(options, tea.WithFPS(a.options.fps))
	}
	program := tea.NewProgram(a.bus.CommandRoot(), options...)
	_, err := program.Run()
	return err
}

func (a *App) CommandBus() commands.Bus {
	return a.bus
}
