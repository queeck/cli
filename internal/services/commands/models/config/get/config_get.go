package get

import (
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/queeck/cli/internal/pkg/cli"
	"github.com/queeck/cli/internal/pkg/keymaps"
	"github.com/queeck/cli/internal/services/commands"
	"github.com/queeck/cli/internal/services/config"
	"github.com/queeck/cli/internal/services/templates"
)

const (
	Code = "get"
)

var _ commands.Command = &Model{} // check for interface compatibility

type Model struct {
	keymap    keymaps.InputKeymap
	help      help.Model
	inputKey  textinput.Model
	templates templates.RendererConfigGet
	arguments cli.Arguments
	parent    func(commands.Command) commands.Command
	config    config.Config
	key       string
	value     string
	has       bool
	quitting  bool
	selected  bool
}

func newInputKey(suggestions []string) textinput.Model {
	input := textinput.New()
	input.Placeholder = "<key.with.dots>"
	input.Prompt = "config/"
	input.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	input.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	input.Focus()
	input.SetSuggestions(suggestions)
	input.ShowSuggestions = true
	input.CharLimit = 120
	input.Width = 120
	return input
}

func New(bus commands.Bus) commands.Command {
	return &Model{
		keymap:    keymaps.Input(),
		help:      help.New(),
		inputKey:  newInputKey(bus.Config().Keys()),
		templates: bus.Templates(),
		arguments: bus.Arguments(),
		parent:    bus.Parent,
		config:    bus.Config(),
	}
}

type SelectedMsg struct{}

func SelectedMessage() tea.Cmd {
	return func() tea.Msg {
		return SelectedMsg{}
	}
}

func (m *Model) Code() string {
	return Code
}

func (m *Model) Commands() []commands.Variant {
	return nil
}

func (m *Model) Init() tea.Cmd {
	args := m.arguments.Commands()
	index := slices.Index(args, Code)
	if index+1 < len(args) {
		m.key = args[index+1]
		return SelectedMessage()
	}
	return textinput.Blink
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg, tea.MouseMsg:
		return m, nil
	case SelectedMsg:
		m.selected = true
		return m.get()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Left):
			return m.parent(m), nil
		case key.Matches(msg, m.keymap.Select):
			m.key = m.inputKey.Value()
			return m.get()
		case key.Matches(msg, m.keymap.Quit):
			return m.quit()
		}
	}

	m.inputKey, cmd = m.inputKey.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	if m.quitting {
		return withNewLine(m.result())
	}

	return m.templates.RenderConfigGetScreen(m.inputKey.View(), m.help.View(m.keymap))
}

func (m *Model) result() string {
	if m.selected {
		return m.templates.RenderConfigGetValue(m.value)
	}
	if !m.has {
		if m.key == "" {
			return m.templates.RenderConfigGetKeyIsEmpty()
		}
		return m.templates.RenderConfigGetKeyNotFound(m.key)
	}
	return m.templates.RenderConfigGetValueWithKey(m.key, m.value)
}

func (m *Model) get() (tea.Model, tea.Cmd) {
	m.value, m.has = m.config.GetString(m.key)
	return m.quit()
}

func (m *Model) quit() (tea.Model, tea.Cmd) {
	m.quitting = true
	return m, tea.Quit
}

func withNewLine(text string) string {
	if !strings.HasSuffix(text, "\n") {
		text += "\n"
	}
	return text
}
