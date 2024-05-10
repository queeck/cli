package models

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/queeck/cli/internal/services/commands"
)

type InitMsg struct{}

func (InitMsg) Choose(bus commands.Bus) (tea.Model, tea.Cmd) {
	command := bus.CommandByCLIArguments(bus.Arguments())
	if command != nil {
		return command, nil
	}

	return bus.CommandRoot(), nil
}

func InitMessage() tea.Msg {
	return InitMsg{}
}
