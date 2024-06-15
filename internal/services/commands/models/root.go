package models

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/queeck/cli/internal/pkg/keymaps"
	"github.com/queeck/cli/internal/pkg/runtime"
	"github.com/queeck/cli/internal/pkg/styles"
	"github.com/queeck/cli/internal/services/commands"
	"github.com/queeck/cli/internal/services/templates"
)

const (
	defaultHeight = 12 // default lines count for screen

	Code = "root"
)

var commandsList = []commands.Variant{
	{Code: "env", Description: "Commands for environment — vars, add, remove, change"},
	{Code: "config", Description: "Commands for configuration — view, get, set"},
	{Code: "push", Description: "Commands for build and psh current service container to registry"},
	{Code: "deploy", Description: "Commands for deploy from pushed registry image"},
}

var _ commands.Command = &Model{} // check for interface compatibility

type Model struct {
	keymap             keymaps.DefaultKeymap
	help               help.Model
	templates          templates.RendererCommon
	commandByArguments func() commands.Command
	commandConfigView  func() commands.Command
	child              func(command commands.Command, code string) commands.Command
	selectedCommands   func(command commands.Command, code string) string
	selected           int
	quitting           bool
}

func New(bus commands.Bus) commands.Command {
	commandByArguments := func() commands.Command {
		// used as callback because on this moment command bus doesn't have registered routes
		return bus.CommandByCLIArguments(bus.Arguments())
	}
	return &Model{
		keymap:             keymaps.Default(),
		help:               help.New(),
		templates:          bus.Templates(),
		commandByArguments: commandByArguments,
		commandConfigView:  bus.CommandConfigView,
		child:              bus.Child,
		selectedCommands:   bus.SelectedCommands,
	}
}

type InitMsg struct{}

func InitMessage() tea.Msg {
	return InitMsg{}
}

func (m *Model) Code() string {
	return Code
}

func (m *Model) Commands() []commands.Variant {
	return commandsList
}

func (m *Model) Init() tea.Cmd {
	return InitMessage
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case InitMsg:
		if command := m.commandByArguments(); command != nil && command != m {
			cmd := command.Init()
			if cmd == nil {
				return command, nil
			}
			return command.Update(cmd())
		}

		return m, nil
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can gracefully truncate
		// its view as needed.
		m.help.Width = msg.Width
		m.commandConfigView().Update(msg) // Prepare sizing for viewport

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Up):
			m.prev()
		case key.Matches(msg, m.keymap.Down):
			m.next()
		case key.Matches(msg, m.keymap.Left):
			//
		case key.Matches(msg, m.keymap.Right):
			return m.choose()
		case key.Matches(msg, m.keymap.Select):
			return m.choose()
		case key.Matches(msg, m.keymap.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keymap.Quit):
			m.quitting = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *Model) View() string {
	if m.quitting {
		return m.templates.RenderCommonQuit()
	}

	selected := m.Commands()[m.selected]
	list := m.selectedCommands(m, selected.Code)
	equivalent := styles.ColorForegroundHighlight(runtime.Executable() + " " + selected.Code)
	text := m.templates.RenderCommonSelectCommands(list, selected.Description, equivalent)

	helpView := m.help.View(m.keymap)
	height := max(defaultHeight-strings.Count(text, "\n")-strings.Count(helpView, "\n"), 0)

	screen := text + strings.Repeat("\n", height) + helpView
	return screen
}

func (m *Model) next() {
	m.selected = (m.selected + 1) % len(commandsList)
}

func (m *Model) prev() {
	if m.selected == 0 {
		m.selected = len(commandsList)
	}
	m.selected = (m.selected - 1) % len(commandsList)
}

func (m *Model) choose() (tea.Model, tea.Cmd) {
	if command := m.child(m, commandsList[m.selected].Code); command != nil {
		return command, nil
	}
	return m, nil
}
