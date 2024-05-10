package config

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/queeck/cli/internal/services/commands/models/config/view"

	"github.com/queeck/cli/internal/models"
	"github.com/queeck/cli/internal/pkg/keymaps"
	"github.com/queeck/cli/internal/pkg/runtime"
	"github.com/queeck/cli/internal/pkg/styles"
	"github.com/queeck/cli/internal/services/commands"
)

const (
	Code = "config"

	defaultHeight = 12 // default lines count for screen
)

var commandsList = []models.Command{
	{Code: view.Code, Description: "Command for view config keys and values"},
	{Code: "get", Description: "Command for get config value"},
	{Code: "set", Description: "Command for set config value"},
}

var _ commands.Command = &Config{} // check for interface compatibility

type Config struct {
	bus      commands.Bus
	keys     keymaps.Default
	help     help.Model
	selected int
	quitting bool
}

func New(bus commands.Bus) commands.Command {
	return &Config{
		keys: keymaps.DefaultKeyMap(),
		help: help.New(),
		bus:  bus,
	}
}

func (m *Config) Code() string {
	return Code
}

func (m *Config) Commands() []models.Command {
	return commandsList
}

func (m *Config) Init() tea.Cmd {
	return nil
}

func (m *Config) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
			return m.bus.CommandRoot(), nil
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

func (m *Config) View() string {
	if m.quitting {
		return m.bus.Templates().RenderCommonQuit()
	}

	selected := m.Commands()[m.selected]
	list := m.bus.SelectedCommands(m, selected.Code)
	equivalent := styles.ColorForegroundHighlight(runtime.Executable() + " config " + selected.Code)
	text := m.bus.Templates().RenderCommonSelectCommands(list, selected.Description, equivalent)

	helpView := m.help.View(m.keys)
	height := max(defaultHeight-strings.Count(text, "\n")-strings.Count(helpView, "\n"), 0)

	screen := text + strings.Repeat("\n", height) + helpView
	return screen
}

func (m *Config) next() {
	m.selected = (m.selected + 1) % len(commandsList)
}

func (m *Config) prev() {
	if m.selected == 0 {
		m.selected = len(commandsList)
	}
	m.selected = (m.selected - 1) % len(commandsList)
}

func (m *Config) choose() (tea.Model, tea.Cmd) {
	if command := m.bus.Child(m, commandsList[m.selected].Code); command != nil {
		return command, nil
	}
	return m, nil
}
