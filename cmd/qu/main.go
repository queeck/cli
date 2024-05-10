package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/queeck/cli/internal/pkg/cli"
	serviceBus "github.com/queeck/cli/internal/services/bus"
	hiera "github.com/queeck/cli/internal/services/commands/hierarchy"
	modelRoot "github.com/queeck/cli/internal/services/commands/models"
	modelConfig "github.com/queeck/cli/internal/services/commands/models/config"
	modelConfigView "github.com/queeck/cli/internal/services/commands/models/config/view"
	serviceConfig "github.com/queeck/cli/internal/services/config"
	serviceDirectory "github.com/queeck/cli/internal/services/directory"
	serviceState "github.com/queeck/cli/internal/services/state"
	serviceTemplates "github.com/queeck/cli/internal/services/templates"
)

func main() {
	if err := run(); err != nil {
		log.Printf("error occurred:\n%v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	directory, err := serviceDirectory.New()
	if err != nil {
		return err
	}

	if !directory.IsLocalCreated() {
		if err = directory.CreateLocal(); err != nil {
			return err
		}
	}

	config, err := serviceConfig.Reload(directory)
	if err != nil {
		return err
	}

	hierarchy := hiera.Node(
		modelRoot.New,
		hiera.Node(
			modelConfig.New,
			hiera.Node(
				modelConfigView.New,
			),
		),
	)

	arguments := cli.Parse(os.Args[1:])

	state := serviceState.New()

	templates := serviceTemplates.New()

	bus, err := serviceBus.New(hierarchy, arguments, state, config, templates)
	if err != nil {
		return err
	}
	program := tea.NewProgram(
		bus.CommandRoot(),
		// tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	_, err = program.Run()
	return err
}
