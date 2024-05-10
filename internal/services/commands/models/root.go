package models

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/queeck/cli/internal/models"
	"github.com/queeck/cli/internal/pkg/keymaps"
	"github.com/queeck/cli/internal/pkg/runtime"
	"github.com/queeck/cli/internal/pkg/styles"
	"github.com/queeck/cli/internal/services/commands"
)

const (
	defaultHeight = 12 // default lines count for screen

	Code = "root"
)

var commandsList = []models.Command{
	{Code: "env", Description: "Commands for environment — vars, add, remove, change"},
	{Code: "config", Description: "Commands for configuration — view, get, set"},
	{Code: "push", Description: "Commands for build and psh current service container to registry"},
	{Code: "deploy", Description: "Commands for deploy from pushed registry image"},
}

var _ commands.Command = &Root{} // check for interface compatibility

type Root struct {
	bus      commands.Bus
	keys     keymaps.Default
	help     help.Model
	selected int
	quitting bool
}

func New(commands commands.Bus) commands.Command {
	return &Root{
		keys: keymaps.DefaultKeyMap(),
		help: help.New(),
		bus:  commands,
	}
}

func (m *Root) Code() string {
	return Code
}

func (m *Root) Commands() []models.Command {
	return commandsList
}

func (m *Root) Init() tea.Cmd {
	return InitMessage
}

func (m *Root) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case InitMsg:
		return msg.Choose(m.bus)
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can gracefully truncate
		// its view as needed.
		m.help.Width = msg.Width
		m.bus.CommandConfigView().Update(msg) // Prepare sizing for viewport

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
			m.prev()
		case key.Matches(msg, m.keys.Down):
			m.next()
		case key.Matches(msg, m.keys.Left):
			//
		case key.Matches(msg, m.keys.Right):
			return m.choose()
		case key.Matches(msg, m.keys.Select):
			return m.choose()
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *Root) View() string {
	if m.quitting {
		return m.bus.Templates().RenderCommonQuit()
	}

	selected := m.Commands()[m.selected]
	list := m.bus.SelectedCommands(m, selected.Code)
	equivalent := styles.ColorForegroundHighlight(runtime.Executable() + " " + selected.Code)
	text := m.bus.Templates().RenderCommonSelectCommands(list, selected.Description, equivalent)

	helpView := m.help.View(m.keys)
	height := max(defaultHeight-strings.Count(text, "\n")-strings.Count(helpView, "\n"), 0)

	screen := text + strings.Repeat("\n", height) + helpView
	return screen
}

func (m *Root) next() {
	m.selected = (m.selected + 1) % len(commandsList)
}

func (m *Root) prev() {
	if m.selected == 0 {
		m.selected = len(commandsList)
	}
	m.selected = (m.selected - 1) % len(commandsList)
}

func (m *Root) choose() (tea.Model, tea.Cmd) {
	if command := m.bus.Child(m, commandsList[m.selected].Code); command != nil {
		return command, nil
	}
	return m, nil
}
